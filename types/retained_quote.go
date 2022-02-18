package types

type RQState uint32

const (
	RQStateWaitingForDeposit RQState = iota
	RQStateTimeForDepositElapsed
	RQStateCallForUserSucceeded
	RQStateCallForUserFailed
	RQStateRegisterPegInSucceeded
	RQStateRegisterPegInFailed
	RQStateRefundLiquiditySucceeded
	RQStateRefundLiquidityFailed
	RQStateInternalError
)

type RetainedQuote struct {
	QuoteHash   string  `json:"quoteHash" db:"quote_hash"`
	DepositAddr string  `json:"depositAddr" db:"deposit_addr"`
	Signature   string  `json:"signature" db:"signature"`
	ReqLiq      uint64  `json:"reqLiq" db:"req_liq"`
	State       RQState `json:"state" db:"state"`
}
