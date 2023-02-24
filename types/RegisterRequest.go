package types

type ProviderRegisterRequest struct {
    Name             string `json:"name"`
    Fee              uint   `json:"fee"`
    QuoteExpiration  uint   `json:"quoteExpiration"`
    AcceptedQuoteExpiration uint `json:"acceptedQuoteExpiration"`
    MinTransactionValue uint `json:"minTransactionValue"`
    MaxTransactionValue uint `json:"maxTransactionValue"`
    ApiBaseUrl       string `json:"apiBaseUrl"`
    Status           bool   `json:"status"`
}
