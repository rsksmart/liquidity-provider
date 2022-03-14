package providers

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"

	"bytes"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

type LiquidityProvider interface {
	GetQuote(types.Quote, uint64, uint64) *types.Quote
	Address() string
	SignQuote(hash []byte, depositAddr string, reqLiq *big.Int) ([]byte, error)
	SignTx(common.Address, *gethTypes.Transaction) (*gethTypes.Transaction, error)
	SetLiquidity(value *big.Int)
	RefundLiquidity(hash []byte) error
}

type RetainedQuotesRepository interface {
	RetainQuote(*types.RetainedQuote) error
	GetRetainedQuote(hash string) (*types.RetainedQuote, error) // returns nil, if not found
}

type LocalProvider struct {
	account    *accounts.Account
	ks         *keystore.KeyStore
	cfg        ProviderConfig
	liquidity  *big.Int
	repository RetainedQuotesRepository
}

type ProviderConfig struct {
	Keydir         string
	BtcAddr        string
	AccountNum     int
	PwdFile        string
	ChainId        *big.Int
	MaxConf        uint16
	Confirmations  map[int]uint16
	TimeForDeposit uint32
	CallTime       uint32
	CallFee        uint64
	PenaltyFee     uint64
}

type InMemRetainedQuotesRepository struct {
	retainedQuotes map[string]*types.RetainedQuote
}

func NewLocalProvider(config ProviderConfig, repository RetainedQuotesRepository) (*LocalProvider, error) {
	if config.Keydir == "" {
		config.Keydir = "keystore"
	}
	if err := os.MkdirAll(config.Keydir, 0700); err != nil {
		return nil, err
	}
	var f *os.File
	if config.PwdFile != "" {
		var err error
		f, err = os.Open(config.PwdFile)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", config.PwdFile)
		}
		defer f.Close()
	}

	ks := keystore.NewKeyStore(config.Keydir, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := retrieveOrCreateAccount(ks, config.AccountNum, f)

	if err != nil {
		return nil, err
	}
	lp := LocalProvider{
		account:    acc,
		ks:         ks,
		cfg:        config,
		liquidity:  big.NewInt(0),
		repository: repository,
	}
	return &lp, nil
}

func (lp *LocalProvider) GetQuote(q types.Quote, gas uint64, gasPrice uint64) *types.Quote {
	q.LPBTCAddr = lp.cfg.BtcAddr
	q.LPRSKAddr = lp.account.Address.String()
	q.AgreementTimestamp = uint32(time.Now().Unix())
	q.Nonce = int64(rand.Int())
	q.TimeForDeposit = lp.cfg.TimeForDeposit
	q.CallTime = lp.cfg.CallTime
	q.PenaltyFee = lp.cfg.PenaltyFee

	q.Confirmations = lp.cfg.MaxConf
	for _, k := range sortedConfirmations(lp.cfg.Confirmations) {
		v := lp.cfg.Confirmations[k]

		if q.Value < uint64(k) {
			q.Confirmations = v
			break
		}
	}
	callCostInSatoshi := weiToSatoshi(gasPrice * gas)
	q.CallFee = uint64(math.Ceil(callCostInSatoshi)) + lp.cfg.CallFee
	return &q
}

func (lp *LocalProvider) Address() string {
	return lp.account.Address.String()
}

func (lp *LocalProvider) SetLiquidity(value *big.Int) {
	lp.liquidity = value
}

func (lp *LocalProvider) RefundLiquidity(hash []byte) error {
	h := hex.EncodeToString(hash)
	rq, err := lp.repository.GetRetainedQuote(h)
	if err != nil {
		return err
	}
	if rq == nil {
		return fmt.Errorf("retained quote not found: %s", h)
	}
	lp.liquidity.Add(lp.liquidity, big.NewInt(int64(rq.ReqLiq)))
	return nil
}

func (lp *LocalProvider) SignQuote(hash []byte, depositAddr string, reqLiq *big.Int) ([]byte, error) {
	if lp.liquidity.Int64()-reqLiq.Int64() < 0 {
		return nil, fmt.Errorf("not enough liquidity. required: %v", reqLiq)
	}
	quoteHash := hex.EncodeToString(hash)
	rq, err := lp.repository.GetRetainedQuote(quoteHash)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("\x19Ethereum Signed Message:\n32")
	buf.Write(hash)

	signB, err := lp.ks.SignHash(*lp.account, crypto.Keccak256(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	signB[len(signB)-1] += 27 // v must be 27 or 28

	if rq == nil {
		signature := hex.EncodeToString(signB)
		rq := types.RetainedQuote{
			QuoteHash:   quoteHash,
			DepositAddr: depositAddr,
			Signature:   signature,
			ReqLiq:      reqLiq.Uint64(),
			State:       types.RQStateWaitingForDeposit,
		}
		err = lp.repository.RetainQuote(&rq)
		if err != nil {
			return nil, err
		}

		lp.liquidity.Sub(lp.liquidity, reqLiq)
	}

	return signB, nil
}

func (lp *LocalProvider) SignTx(address common.Address, tx *gethTypes.Transaction) (*gethTypes.Transaction, error) {
	if !bytes.Equal(address[:], lp.account.Address[:]) {
		return nil, fmt.Errorf("provider address %v is incorrect", address.Hash())
	}
	return lp.ks.SignTx(*lp.account, tx, lp.cfg.ChainId)
}

func weiToSatoshi(wei uint64) float64 {
	return float64(wei) / math.Pow10(10)
}

func retrieveOrCreateAccount(ks *keystore.KeyStore, accountNum int, in *os.File) (*accounts.Account, error) {
	if cap(ks.Accounts()) == 0 {
		log.Info("no RSK account found")
		acc, err := createAccount(ks, in)
		return acc, err
	} else {
		if cap(ks.Accounts()) <= int(accountNum) {
			return nil, fmt.Errorf("account number %v not found", accountNum)
		}
		acc := ks.Accounts()[accountNum]
		passwd, err := enterPasswd(in)

		if err != nil {
			return nil, err
		}
		err = ks.Unlock(acc, passwd)
		return &acc, err
	}
}

func createAccount(ks *keystore.KeyStore, in *os.File) (*accounts.Account, error) {
	passwd, err := createPasswd(in)

	if err != nil {
		return nil, err
	}
	acc, err := ks.NewAccount(passwd)

	if err != nil {
		return &acc, err
	}
	err = ks.Unlock(acc, passwd)

	if err != nil {
		return &acc, err
	}
	log.Info("new account created: ", acc.Address)
	return &acc, err
}

func enterPasswd(in *os.File) (string, error) {
	fmt.Println("enter password for RSK account")
	fmt.Print("password: ")
	var pwd string
	var err error
	if in == nil {
		pwd, err = readPasswdCons(nil)
	} else {
		pwd, err = readPasswdReader(bufio.NewReader(in))
	}
	fmt.Println()
	return pwd, err
}

func createPasswd(in *os.File) (string, error) {
	fmt.Println("creating password for new RSK account")
	fmt.Println("WARNING: the account will be lost forever if you forget this password!!! Do you understand? (yes/[no])")

	var r *bufio.Reader
	var readPasswd func(*bufio.Reader) (string, error)
	if in == nil {
		r = bufio.NewReader(os.Stdin)
		readPasswd = readPasswdCons
	} else {
		r = bufio.NewReader(in)
		readPasswd = readPasswdReader
	}

	str, _ := r.ReadString('\n')
	if str != "yes\n" {
		return "", errors.New("must say yes")
	}
	fmt.Print("password: ")
	pwd1, err := readPasswd(r)
	fmt.Println()
	if err != nil {
		return "", err
	}

	fmt.Print("repeat password: ")
	pwd2, err := readPasswd(r)
	fmt.Println()
	if err != nil {
		return "", err
	}
	if pwd1 != pwd2 {
		return "", errors.New("passwords do not match")
	}
	return pwd1, nil
}

func readPasswdCons(r *bufio.Reader) (string, error) {
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	return string(bytes), err
}

func readPasswdReader(r *bufio.Reader) (string, error) {
	str, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(str, "\n"), nil
}

func sortedConfirmations(m map[int]uint16) []int {
	keys := make([]int, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

func NewInMemRetainedQuotesRepository() *InMemRetainedQuotesRepository {
	return &InMemRetainedQuotesRepository{
		retainedQuotes: make(map[string]*types.RetainedQuote),
	}
}

func (s InMemRetainedQuotesRepository) RetainQuote(quote *types.RetainedQuote) error {
	s.retainedQuotes[quote.QuoteHash] = quote
	return nil
}

func (s InMemRetainedQuotesRepository) GetRetainedQuote(hash string) (*types.RetainedQuote, error) {
	q, ok := s.retainedQuotes[hash]
	if !ok {
		return nil, nil
	}
	return q, nil
}
