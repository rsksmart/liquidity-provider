package types

import (
	"math/big"
)

type Quote struct {
	FedBTCAddr         string  `json:"fedBTCAddr" db:"fed_addr"`
	LBCAddr            string  `json:"lbcAddr" db:"lbc_addr"`
	LPRSKAddr          string  `json:"lpRSKAddr" db:"lp_rsk_addr"`
	BTCRefundAddr      string  `json:"btcRefundAddr" db:"btc_refund_addr"`
	RSKRefundAddr      string  `json:"rskRefundAddr" db:"rsk_refund_addr"`
	LPBTCAddr          string  `json:"lpBTCAddr" db:"lp_btc_addr"`
	CallFee            big.Int `json:"callFee" db:"call_fee"`
	PenaltyFee         big.Int `json:"penaltyFee" db:"penalty_fee"`
	ContractAddr       string  `json:"contractAddr" db:"contract_addr"`
	Data               string  `json:"data" db:"data"`
	GasLimit           uint    `json:"gasLimit" db:"gas_limit"`
	Nonce              uint    `json:"nonce" db:"nonce"`
	Value              big.Int `json:"value" db:"value"`
	AgreementTimestamp uint    `json:"agreementTimestamp" db:"agreement_timestamp"`
	TimeForDeposit     uint    `json:"timeForDeposit" db:"time_for_deposit"`
	CallTime           uint    `json:"callTime" db:"call_time"`
	Confirmations      uint    `json:"confirmations" db:"confirmations"`
}
