package types

type ProviderRegisterRequest struct {
    Name             string `json:"name" example:"New Provider" description:"Provider Name"`
    Fee              uint   `json:"fee" example:"100" description:"Fee in wei"`
    QuoteExpiration  uint   `json:"quoteExpiration" example:"100" description:"Quote expiration in seconds"`
    AcceptedQuoteExpiration uint `json:"acceptedQuoteExpiration" example:"100" description:"Accepted quote expiration in seconds"`
    MinTransactionValue uint `json:"minTransactionValue" example:"100" description:"Minimum transaction value in wei"`
    MaxTransactionValue uint `json:"maxTransactionValue" example:"100" description:"Maximum transaction value in wei"`
    ApiBaseUrl       string `json:"apiBaseUrl" example:"https://api.example.com" description:"API base URL"`
    Status           bool   `json:"status" example:"true" description:"Provider status"`
}
