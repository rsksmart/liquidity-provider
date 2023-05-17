package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)
type UserEvents struct {
	From      common.Address `json:"from" example:"0x0" description:"From Address"`
	Dest      common.Address `json:"dest" example:"0x0" description:"Destination Address"`
	GasLimit  *big.Int       `json:"gasLimit" example:"10000" description:"Gas Limit"`
	Value     *big.Int       `json:"value" example:"10000" description:"Event Value"`
	Data      []byte         `json:"data" example:"[]" description:"Event Data"`
	Success   bool	         `json:"success" example:"true" description:"Event Status"`
	QuoteHash string         `json:"quoteHash" example:"0x0" description:"QuoteHash"`
}
