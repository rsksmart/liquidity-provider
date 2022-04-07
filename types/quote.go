package types

type Quote struct {
	FedBTCAddr         string `json:"fedBTCAddr" db:"fed_addr"`
	LBCAddr            string `json:"lbcAddr" db:"lbc_addr"`
	LPRSKAddr          string `json:"lpRSKAddr" db:"lp_rsk_addr"`
	BTCRefundAddr      string `json:"btcRefundAddr" db:"btc_refund_addr"`
	RSKRefundAddr      string `json:"rskRefundAddr" db:"rsk_refund_addr"`
	LPBTCAddr          string `json:"lpBTCAddr" db:"lp_btc_addr"`
	CallFee            *Wei   `json:"callFee" db:"call_fee"`
	PenaltyFee         *Wei   `json:"penaltyFee" db:"penalty_fee"`
	ContractAddr       string `json:"contractAddr" db:"contract_addr"`
	Data               string `json:"data" db:"data"`
	GasLimit           uint32 `json:"gasLimit" db:"gas_limit"`
	Nonce              int64  `json:"nonce" db:"nonce"`
	Value              *Wei   `json:"value" db:"value"`
	AgreementTimestamp uint32 `json:"agreementTimestamp" db:"agreement_timestamp"`
	TimeForDeposit     uint32 `json:"timeForDeposit" db:"time_for_deposit"`
	CallTime           uint32 `json:"callTime" db:"call_time"`
	Confirmations      uint16 `json:"confirmations" db:"confirmations"`
	CallOnRegister     bool   `json:"callOnRegister" db:"call_on_register"`
}
