package types

import (
	"math/big"
)

type Quote struct {
	FedBTCAddr         string  `json:"fedBTCAddr" db:"fed_addr" abi:"fedBtcAddress"`
	LBCAddr            string  `json:"lbcAddr" db:"lbc_addr" abi:"lbcAddress"`
	LPRSKAddr          string  `json:"lpRSKAddr" db:"lp_rsk_addr" abi:"liquidityProviderRskAddress"`
	BTCRefundAddr      string  `json:"btcRefundAddr" db:"btc_refund_addr" abi:"btcRefundAddress"`
	RSKRefundAddr      string  `json:"rskRefundAddr" db:"rsk_refund_addr" abi:"rskRefundAddress"`
	LPBTCAddr          string  `json:"lpBTCAddr" db:"lp_btc_addr" abi:"liquidityProviderBtcAddress"`
	CallFee            big.Int `json:"callFee" db:"call_fee" abi:"callFee"`
	ContractAddr       string  `json:"contractAddr" db:"contract_addr" abi:"contractAddress"`
	Data               string  `json:"data" db:"data" abi:"data"`
	GasLimit           uint    `json:"gasLimit" db:"gas_limit" abi:"gasLimit"`
	Nonce              uint    `json:"nonce" db:"nonce" abi:"nonce"`
	Value              big.Int `json:"value" db:"value" abi:"value"`
	AgreementTimestamp uint    `json:"agreementTimestamp" db:"agreement_timestamp" abi:"agreementTimestamp"`
	TimeForDeposit     uint    `json:"timeForDeposit" db:"time_for_deposit" abi:"timeForDeposit"`
	CallTime           uint    `json:"callTime" db:"call_time" abi:"callTime"`
	Confirmations      uint    `json:"confirmations" db:"confirmations" abi:"depositConfirmations"`
}
