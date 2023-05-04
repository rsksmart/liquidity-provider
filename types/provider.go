package types

type GlobalProvider struct {
    Id               uint64 `json:"id" example:"1" description:"Provider Id"`
    Provider         string `json:"provider" example:"0x0" description:"Provider Address"`
    Name             string `json:"name" example:"New Provider" description:"Provider Name"`
    Fee              uint64   `json:"fee" example:"100" description:"Fee in wei"`
    QuoteExpiration  uint64   `json:"quoteExpiration" example:"100" description:"Quote expiration in seconds"`
    AcceptedQuoteExpiration uint64 `json:"acceptedQuoteExpiration" example:"100" description:"Accepted quote expiration in seconds"`
    MinTransactionValue uint64 `json:"minTransactionValue" example:"100" description:"Minimum transaction value in wei"`
    MaxTransactionValue uint64 `json:"maxTransactionValue" example:"100" description:"Maximum transaction value in wei"`
    ApiBaseUrl       string `json:"apiBaseUrl" example:"https://api.example.com" description:"API base URL"`
    Status           bool   `json:"status" example:"true" description:"Provider status"`
    ProviderType     string `json:"providerType" example:"pegin" description:"Provider type"`
}
