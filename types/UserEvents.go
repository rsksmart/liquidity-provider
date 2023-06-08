package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)
type UserEvents struct {
	From      common.Address `json:"from" example:"0x0" description:"From Address"`
	Amount     *big.Int      `json:"amount" example:"10000" description:"Event Value"`
	Timestamp  *big.Int      `json:"timestamp" example:"10000" description:"Event Timestamp"`
	QuoteHash string         `json:"quoteHash" example:"0x0" description:"QuoteHash"`
}
