package types

type RQState uint32

const (
	RQStateWaitingForDeposit RQState = iota
	RQStateTimeForDepositElapsed
	RQStateCallForUserSucceeded
	RQStateCallForUserFailed
	RQStateRegisterPegInSucceeded
	RQStateRegisterPegInFailed
	RQStateWaitingForDepositConfirmations
)

type RetainedQuote struct {
	QuoteHash           string  `json:"quoteHash" db:"quote_hash"`
	DepositAddr         string  `json:"depositAddr" db:"deposit_addr"`
	Signature           string  `json:"signature" db:"signature"`
	FlyoverRedeemScript string  `json:"flyoverRedeemScript" db:"flyover_redeem_script"`
	ReqLiq              *Wei    `json:"reqLiq" db:"req_liq"`
	State               RQState `json:"state" db:"state"`
}
