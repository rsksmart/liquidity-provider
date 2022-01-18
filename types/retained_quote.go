package types

type RetainedQuote struct {
	QuoteHash     string `json:"quoteHash" db:"quote_hash"`
	DepositAddr   string `json:"depositAddr" db:"deposit_addr"`
	Signature     string `json:"signature" db:"signature"`
	CalledForUser bool   `json:"calledForUser" db:"called_for_user"`
	ReqLiq        uint64 `json:"reqLiq" db:"req_liq"`
}
