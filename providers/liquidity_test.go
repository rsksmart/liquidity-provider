package providers

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/rsksmart/liquidity-provider/types"
)

func newLocalProvider(t *testing.T) *LocalProvider {
	f := getFile("yes\n1234\n1234\n", t)
	defer f.Close()

	lp, err := NewLocalProvider(t.TempDir(), 0, f)
	if err != nil {
		t.Fatal("error creating local provider: ", err)
	}
	return lp
}

func TestCreatePassword(t *testing.T) {
	f1 := getFile("yes\n1234\n1234\n", t)
	defer f1.Close()
	pwd, err := createPasswd(f1)
	if err != nil {
		t.Fatal(err)
	}
	if pwd != "1234" {
		t.Fatalf("expected 1234, got %v", pwd)
	}

	f2 := getFile("yes\n1234\n14\n", t)
	defer f2.Close()
	_, err = createPasswd(f2)
	if err == nil {
		t.Fatal("did not fail when passwords do not match")
	}

	f3 := getFile("nah\n1234\n1234\n", t)
	defer f3.Close()
	_, err = createPasswd(f3)
	if err == nil {
		t.Fatal("did not fail when yes is not typed")
	}
}

func TestLocalProvider(t *testing.T) {
	t.Run("new", testNewLocal)
	t.Run("get quote", testGetQuoteLocal)
	t.Run("sign hash", testSignHashLocal)
}

func testNewLocal(t *testing.T) {
	lp := newLocalProvider(t)
	if lp.account == nil {
		t.Fatalf("account is empty")
	}
	if lp.ks == nil {
		t.Fatalf("keystore is empty")
	}
}

func testGetQuoteLocal(t *testing.T) {
	q := types.Quote{
		ContractAddr: "222",
		Value:        *big.NewInt(5),
		Data:         "0x0",
	}
	lp := newLocalProvider(t)
	nq := lp.GetQuote(q, 10, *big.NewInt(3))

	if nq == nil {
		t.Fatal("empty quote")
	}
	if nq.AgreementTimestamp <= 0 {
		t.Fatalf("invalid agreement timestamp: %v", nq.AgreementTimestamp)
	}
	if nq.CallTime <= 0 {
		t.Fatalf("invalid call time: %v", nq.CallTime)
	}
	if nq.CallFee.Cmp(big.NewInt(0)) < 0 {
		t.Fatal("invalid call fee")
	}
	if nq.PenaltyFee.Cmp(big.NewInt(0)) < 0 {
		t.Fatal("invalid penalty fee")
	}
	if nq.Confirmations == 0 {
		t.Fatal("invalid confirmations")
	}
	if nq.Nonce == 0 {
		t.Fatal("nonce is 0")
	}
	if nq.TimeForDeposit == 0 {
		t.Fatal("time for deposit is 0")
	}
}

func testSignHashLocal(t *testing.T) {
	lp := newLocalProvider(t)
	b, err := lp.SignHash([]byte("12345678901234567890123456789012"))

	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("empty signature")
	}
}

func getFile(s string, t *testing.T) *os.File {
	tmpFile, err := ioutil.TempFile(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tmpFile.WriteString(s)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	f, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	return f
}
