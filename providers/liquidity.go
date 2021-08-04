package providers

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
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
	SignHash() []byte
}

type LocalProvider struct {
	account *accounts.Account
	ks      *keystore.KeyStore
}

func NewLocalProvider(keydir string, accountNum int) (*LocalProvider, error) {
	kd := keydir

	if kd == "" {
		kd = "keystore"
	}
	if err := os.MkdirAll(kd, 0700); err != nil {
		return nil, err
	}
	ks := keystore.NewKeyStore(kd, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := retreiveOrCreateAccount(ks, accountNum)

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
	q.Nonce = uint(rand.Uint64())
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

func retreiveOrCreateAccount(ks *keystore.KeyStore, accountNum int) (*accounts.Account, error) {
	if cap(ks.Accounts()) == 0 {
		log.Info("no RSK account found")
		acc, err := createAccount(ks)
		return acc, err
	} else {
		if cap(ks.Accounts()) <= int(accountNum) {
			return nil, fmt.Errorf("account number %v not found", accountNum)
		}
		acc := ks.Accounts()[accountNum]
		passwd, err := readPasswd()

		if err != nil {
			return nil, err
		}
		err = ks.Unlock(acc, passwd)
		return &acc, err
	}
}

func createAccount(ks *keystore.KeyStore) (*accounts.Account, error) {
	passwd, err := createPasswd()

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

func readPasswd() (string, error) {
	fmt.Println("enter password for RSK account")
	fmt.Print("password: ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(bytepw), err
}

func createPasswd() (string, error) {
	fmt.Println("creating password for new RSK account")
	fmt.Println("WARNING: the account will be lost forever if you forget this password!!! Do you understand? (yes/[no])")

	r := bufio.NewReader(os.Stdin)
	str, _ := r.ReadString('\n')

	if str != "yes\n" {
		return "", errors.New("must type yes")
	}
	fmt.Print("password: ")
	bytepw1, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if err != nil {
		return "", err
	}
	fmt.Print("repeat password: ")
	bytepw2, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if err != nil {
		return "", err
	}
	if string(bytepw1) != string(bytepw2) {
		return "", errors.New("passwords do not match")
	}
	return string(bytepw1), nil
}
