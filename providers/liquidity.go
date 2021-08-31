package providers

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

type LiquidityProvider interface {
	GetQuote(types.Quote, uint64, big.Int) *types.Quote
	Address() string
	SignHash(hash []byte) ([]byte, error)
}

type LocalProvider struct {
	account *accounts.Account
	ks      *keystore.KeyStore
}

func NewLocalProvider(keydir string, accountNum int, in *os.File) (*LocalProvider, error) {
	kd := keydir

	if kd == "" {
		kd = "keystore"
	}
	if err := os.MkdirAll(kd, 0700); err != nil {
		return nil, err
	}
	ks := keystore.NewKeyStore(kd, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := retreiveOrCreateAccount(ks, accountNum, in)

	if err != nil {
		return nil, err
	}
	lp := LocalProvider{
		account: acc,
		ks:      ks,
	}
	return &lp, nil
}

func (lp *LocalProvider) GetQuote(q types.Quote, gas uint64, gasPrice big.Int) *types.Quote {
	q.LPRSKAddr = lp.account.Address.String()
	// TODO better way to compute fee, times, etc.
	q.Confirmations = 10
	q.TimeForDeposit = 3600
	q.CallTime = 3600
	q.PenaltyFee = *big.NewInt(10)
	q.Nonce = rand.Int()
	cost := big.NewInt(0).Mul(&gasPrice, big.NewInt(int64(gas)))
	fee := cost.Div(cost, big.NewInt(33))
	finalCost := cost.Add(cost, fee)
	q.CallFee = *finalCost
	q.AgreementTimestamp = uint(time.Now().Unix())
	return &q
}

func (lp *LocalProvider) Address() string {
	return lp.account.Address.String()
}

func (lp *LocalProvider) SignHash(hash []byte) ([]byte, error) {
	return lp.ks.SignHash(*lp.account, hash)
}

func retreiveOrCreateAccount(ks *keystore.KeyStore, accountNum int, in *os.File) (*accounts.Account, error) {
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
