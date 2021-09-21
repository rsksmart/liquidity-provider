package providers

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/rsksmart/liquidity-provider/types"
)

type signature struct {
	h string
	s string
}

type getQuoteData struct {
	inQ       types.Quote
	gas       uint
	gasPrice  uint64
	expectedQ types.Quote
}

var (
	btcAddr = "123"

	expectedSign = [2]signature{
		{
			h: "4545454545454545454545454545454545454545454545454545454545454545",
			s: "5747fc2a9327abf9a3dd8caf454b41dc867f2f33b6fdf1caad1d0b050ce43ceb35cc2214ed808af3598f44d4dfb15e5dfafd88c6367b0c9820b4076a22e3dcdc1b",
		},
		{
			h: "4545454545454545454545454545454545454545454545454523454545454545",
			s: "709e2e47aa3b77fd151d3595753b06f155fab4427adc08e655811feb89edb9fa672d997c3e21449a6fd01b22f4d43483e8bb9dc249488e77cf3d3c265a20652b1c",
		},
	}

	testQuotes = [2]getQuoteData{
		{
			inQ: types.Quote{
				Value:    *big.NewInt(3000000),
				CallFee:  1000,
				GasLimit: 50000,
			},
			gas:      50000,
			gasPrice: 10,
			expectedQ: types.Quote{
				Confirmations: 6,
				CallFee:       501000,
			},
		},
		{
			inQ: types.Quote{
				Value:    *big.NewInt(100000000),
				CallFee:  1000,
				GasLimit: 50000,
			},
			gas:      50000,
			gasPrice: 10,
			expectedQ: types.Quote{
				Confirmations: 60,
				CallFee:       501000,
			},
		},
	}
)

func testSignature(t *testing.T) {
	f := genTmpFile("1234\n1234\n", t)
	defer f.Close()

	cfg := ProviderConfig{
		Keydir:     "./testdata/keystore/keystore",
		AccountNum: 0,
		PwdFile:    f.Name(),
	}

	p, err := NewLocalProvider(cfg)
	if err != nil {
		t.Error(err)
	}

	for _, sign := range expectedSign {
		h, _ := hex.DecodeString(sign.h)

		b, err := p.SignHash(h)
		if err != nil {
			t.Errorf("error signing hash: %v", err)
		}
		if hex.EncodeToString(b) != sign.s {
			t.Errorf("wrong signature. got: %x \n expected: %v", b, sign.s)
		}
	}
}

func testCreatePassword(t *testing.T) {
	f1 := genTmpFile("yes\n1234\n1234\n", t)
	defer f1.Close()
	pwd, err := createPasswd(f1)
	if err != nil {
		t.Fatal(err)
	}
	if pwd != "1234" {
		t.Fatalf("expected 1234, got %v", pwd)
	}

	f2 := genTmpFile("yes\n1234\n14\n", t)
	defer f2.Close()
	_, err = createPasswd(f2)
	if err == nil {
		t.Fatal("did not fail when passwords do not match")
	}

	f3 := genTmpFile("nah\n1234\n1234\n", t)
	defer f3.Close()
	_, err = createPasswd(f3)
	if err == nil {
		t.Fatal("did not fail when yes is not typed")
	}
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
	f, err := os.Open("./testdata/test_config.json")
	if err != nil {
		t.Errorf("error opening test config: %v", err)
	}
	dec := json.NewDecoder(f)
	cfg := ProviderConfig{}
	dec.Decode(&cfg)
	cfg.PwdFile = genTmpFile("1234\n1234\n", t).Name()

	lp, err := NewLocalProvider(cfg)
	if err != nil {
		t.Fatal("error creating local provider: ", err)
	}

	for _, q := range testQuotes {
		nq := lp.GetQuote(q.inQ, uint64(q.gas), q.gasPrice)

		if nq == nil {
			t.Fatal("empty quote")
		}
		if nq.AgreementTimestamp <= 0 {
			t.Fatalf("invalid agreement timestamp: %v", nq.AgreementTimestamp)
		}
		if nq.CallTime != cfg.CallTime {
			t.Fatalf("invalid call time: %v", nq.CallTime)
		}
		if nq.CallFee != q.expectedQ.CallFee {
			t.Fatal("invalid call fee")
		}
		if nq.PenaltyFee != cfg.PenaltyFee {
			t.Fatal("invalid penalty fee")
		}
		if nq.Confirmations != q.expectedQ.Confirmations {
			t.Fatalf("invalid confirmations: %v", nq.Confirmations)
		}
		if nq.Nonce == 0 {
			t.Fatal("nonce is 0")
		}
		if nq.TimeForDeposit != cfg.TimeForDeposit {
			t.Fatal("time for deposit is 0")
		}
		if nq.LPBTCAddr != cfg.BtcAddr {
			t.Fatal("bitcoin address wasn't set")
		}
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

func newLocalProvider(t *testing.T) *LocalProvider {
	f := genTmpFile("yes\n1234\n1234\n", t)
	cfg := ProviderConfig{
		BtcAddr:    btcAddr,
		Keydir:     t.TempDir(),
		AccountNum: 0,
		PwdFile:    f.Name(),
	}
	defer f.Close()

	lp, err := NewLocalProvider(cfg)
	if err != nil {
		t.Fatal("error creating local provider: ", err)
	}
	return lp
}

func genTmpFile(s string, t *testing.T) *os.File {
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

func TestLocalProvider(t *testing.T) {
	t.Run("new", testNewLocal)
	t.Run("get quote", testGetQuoteLocal)
	t.Run("sign hash", testSignHashLocal)
	t.Run("create password", testCreatePassword)
	t.Run("signature", testSignature)

}
