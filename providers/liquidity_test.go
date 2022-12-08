package providers

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/assert"
)

type signature struct {
	h string
	s string
}

type getQuoteData struct {
	inQ       *types.Quote
	gas       uint64
	gasPrice  *types.Wei
	expectedQ *types.Quote
}

var (
	btcAddr = "123"

	expectedSign = [2]signature{
		{
			h: "4545454545454545454545454545454545454545454545454545454545454545",
			s: "329389e8c4cb329cdcb88e44e524abedc7492a8d9b210037879874ce8103d5a2491f947010a18f7fab43d91b35e8ec8cf7760f6ab8806f53ef1c7644e7bf8b741b",
		},
		{
			h: "4545454545454545454545454545454545454545454545454523454545454545",
			s: "02333b55f159f2e21b58f22f3f06cf1ed4aaec74caa43be736ab7d09de119965522557fef9f4ecb1b61950c732fac49a81c18b8d7972c8f7e3d89b82cd1dcb141b",
		},
	}

	testQuotes = [2]getQuoteData{
		{
			inQ: &types.Quote{
				Value:    types.NewWei(3000000),
				CallFee:  types.NewWei(1000),
				GasLimit: 50000,
			},
			gas:      50000,
			gasPrice: types.NewWei(10),
			expectedQ: &types.Quote{
				Confirmations: 6,
				CallFee:       types.NewWei(501000),
			},
		},
		{
			inQ: &types.Quote{
				Value:    types.NewWei(100000000),
				CallFee:  types.NewWei(1000),
				GasLimit: 50000,
			},
			gas:      50000,
			gasPrice: types.NewWei(10),
			expectedQ: &types.Quote{
				Confirmations: 60,
				CallFee:       types.NewWei(501000),
			},
		},
	}
)

type InMemLocalProviderRepository struct {
	retainedQuotes map[string]*types.RetainedQuote
	liquidity      *types.Wei
}

func NewInMemRetainedQuotesRepository() *InMemLocalProviderRepository {
	return &InMemLocalProviderRepository{
		retainedQuotes: make(map[string]*types.RetainedQuote),
		liquidity:      types.NewWei(0),
	}
}

func (r *InMemLocalProviderRepository) RetainQuote(quote *types.RetainedQuote) error {
	r.retainedQuotes[quote.QuoteHash] = quote
	return nil
}

func (r *InMemLocalProviderRepository) HasRetainedQuote(hash string) (bool, error) {
	_, ok := r.retainedQuotes[hash]
	return ok, nil
}

func (r *InMemLocalProviderRepository) GetLiquidity() *types.Wei {
	liq := r.liquidity.Copy()

	for _, rq := range r.retainedQuotes {
		if rq.State == types.RQStateWaitingForDeposit {
			liq.Sub(liq, rq.ReqLiq)
		}
	}

	return liq
}

func (r *InMemLocalProviderRepository) HasLiquidity(_ LiquidityProvider, wei *types.Wei) (bool, error) {
	return r.GetLiquidity().Cmp(wei) >= 0, nil
}

func (r *InMemLocalProviderRepository) SetLiquidity(liq *types.Wei) {
	r.liquidity = liq.Copy()
}

func (r *InMemLocalProviderRepository) SetRetainedQuoteState(hash string, state types.RQState) {
	q, ok := r.retainedQuotes[hash]
	if ok {
		q.State = state
	}
}

func testSignature(t *testing.T) {
	f := genTmpFile("correct horse battery staple\ncorrect horse battery staple\n", t)
	defer f.Close()

	cfg := ProviderConfig{
		Keydir:     "./testdata/keystore/keystore",
		AccountNum: 0,
		PwdFile:    f.Name(),
	}

	repository := NewInMemRetainedQuotesRepository()
	p, err := NewLocalProvider(cfg, repository)
	if err != nil {
		t.Error(err)
	}

	for _, sign := range expectedSign {
		reqLiq := types.NewWei(200)
		repository.SetLiquidity(reqLiq)
		h, _ := hex.DecodeString(sign.h)

		b, err := p.SignQuote(h, "abc", reqLiq)
		if err != nil {
			t.Errorf("error signing hash: %v", err)
		}
		if hex.EncodeToString(b) != sign.s {
			t.Errorf("wrong signature. got: %x \n expected: %v", b, sign.s)
		}
		repository.SetRetainedQuoteState(sign.h, types.RQStateCallForUserSucceeded)
	}
}

func testCreatePassword(t *testing.T) {
	f1 := genTmpFile("yes\ncorrect horse battery staple\ncorrect horse battery staple\n", t)
	defer f1.Close()
	pwd, err := createPasswd(f1)
	if err != nil {
		t.Fatal(err)
	}
	if pwd != "correct horse battery staple" {
		t.Fatalf("expected 1234, got %v", pwd)
	}

	f2 := genTmpFile("yes\ncorrect horse battery staple\ncorrect horse battery step\n", t)
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

	f4 := genTmpFile("yes\n1234\n1234\n", t)
	defer f4.Close()
	_, err = createPasswd(f4)
	if err == nil {
		t.Fatal("did not fail when password is not secure enough")
	}
}

func testNewLocal(t *testing.T) {
	repository := NewInMemRetainedQuotesRepository()
	lp := newLocalProvider(t, repository)
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
	err = dec.Decode(&cfg)
	if err != nil {
		t.Fatal("error decoding config: ", err)
	}
	cfg.PwdFile = genTmpFile("correct horse battery staple\ncorrect horse battery staple\n", t).Name()

	repository := NewInMemRetainedQuotesRepository()
	lp, err := NewLocalProvider(cfg, repository)
	if err != nil {
		t.Fatal("error creating local provider: ", err)
	}

	for _, q := range testQuotes {
		nq, err := lp.GetQuote(q.inQ, q.gas, q.gasPrice)
		if err != nil {
			t.Fatal("error getting quote: ", err)
		}

		if nq == nil {
			t.Fatal("empty quote")
		}
		if nq.AgreementTimestamp <= 0 {
			t.Fatalf("invalid agreement timestamp: %v", nq.AgreementTimestamp)
		}
		if nq.CallTime != cfg.CallTime {
			t.Fatalf("invalid call time: %v", nq.CallTime)
		}
		if nq.CallFee.Cmp(q.expectedQ.CallFee) != 0 {
			t.Fatal("invalid call fee")
		}
		if nq.PenaltyFee.Cmp(cfg.PenaltyFee) != 0 {
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

func testSignQuoteLocal(t *testing.T) {
	repository := NewInMemRetainedQuotesRepository()
	lp := newLocalProvider(t, repository)
	repository.SetLiquidity(types.NewWei(220))
	reqLiq := types.NewWei(200)
	b, err := lp.SignQuote([]byte("12345678901234567890123456789012"), "abc", reqLiq)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("empty signature")
	}

	assert.EqualValues(t, big.NewInt(20), repository.GetLiquidity())
}

func testInsufficientFunds(t *testing.T) {
	repository := NewInMemRetainedQuotesRepository()
	lp := newLocalProvider(t, repository)
	repository.SetLiquidity(types.NewWei(100))
	reqLiq := types.NewWei(101)
	_, err := lp.SignQuote([]byte("12345678901234567890123456789012"), "abc", reqLiq)
	if err != nil {
		assert.Errorf(t, err, "not enough liquidity. required: %v")
	}
}

func testLiquidityFluctuation(t *testing.T) {
	quoteHash := "12345678901234567890123456789012"
	repository := NewInMemRetainedQuotesRepository()
	lp := newLocalProvider(t, repository)
	initialLiq := types.NewWei(100)
	repository.SetLiquidity(initialLiq)
	reqLiq := types.NewWei(90)
	expectedLiq := new(types.Wei).Sub(initialLiq, reqLiq)
	qb, err := hex.DecodeString(quoteHash)
	if err != nil {
		t.Fail()
	}
	_, err = lp.SignQuote(qb, "abc", reqLiq)
	if err != nil {
		t.Fail()
	}

	assert.EqualValues(t, expectedLiq, repository.GetLiquidity())

	repository.SetRetainedQuoteState(quoteHash, types.RQStateCallForUserSucceeded)

	assert.EqualValues(t, initialLiq, repository.GetLiquidity())
}

func testLiquidityAbnormalFluctuation(t *testing.T) {
	quoteHash := "12345678901234567890123456789012"
	repository := NewInMemRetainedQuotesRepository()
	lp := newLocalProvider(t, repository)
	initialLiq := types.NewWei(200)
	repository.SetLiquidity(initialLiq)
	reqLiq := types.NewWei(90)
	expectedLiq := new(types.Wei).Sub(initialLiq, reqLiq)
	qb, err := hex.DecodeString(quoteHash)
	if err != nil {
		t.Fail()
	}
	_, err = lp.SignQuote(qb, "abc", reqLiq)
	if err != nil {
		t.Fail()
	}

	assert.EqualValues(t, expectedLiq, repository.GetLiquidity())

	_, err = lp.SignQuote(qb, "abc", reqLiq)
	if err != nil {
		t.Fail()
	}

	assert.EqualValues(t, expectedLiq, repository.GetLiquidity())

	repository.SetRetainedQuoteState(quoteHash, types.RQStateCallForUserSucceeded)

	assert.EqualValues(t, initialLiq, repository.GetLiquidity())
}

func testSetLiquidity(t *testing.T) {
	repository := NewInMemRetainedQuotesRepository()
	value := types.NewWei(20000)
	repository.SetLiquidity(value)
	assert.EqualValues(t, value, repository.GetLiquidity())
}

func newLocalProvider(t *testing.T, repository LocalProviderRepository) *LocalProvider {
	f := genTmpFile("yes\ncorrect horse battery staple\ncorrect horse battery staple\n", t)
	cfg := ProviderConfig{
		BtcAddr:    btcAddr,
		Keydir:     t.TempDir(),
		AccountNum: 0,
		PwdFile:    f.Name(),
	}
	defer f.Close()

	lp, err := NewLocalProvider(cfg, repository)
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
	t.Run("sign quote", testSignQuoteLocal)
	t.Run("create password", testCreatePassword)
	t.Run("signature", testSignature)
	t.Run("set liquidity", testSetLiquidity)
	t.Run("sign quote with insufficient funds", testInsufficientFunds)
	t.Run("liquidity fluctuation", testLiquidityFluctuation)
	t.Run("liquidity abnormal fluctuation", testLiquidityAbnormalFluctuation)
}
