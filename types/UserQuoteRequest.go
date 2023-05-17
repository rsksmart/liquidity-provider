package types

type UserQuoteRequest struct {
    Address   string  `json:"address" example:"0x0" description:"User Address"`
    FromBlock *uint64 `json:"fromBlock,omitempty" example:"69" description:"optional fromBlock"`
    ToBlock   *uint64 `json:"toBlock,omitempty" example:"500 or latest" description:"optional toBlock"`
}
