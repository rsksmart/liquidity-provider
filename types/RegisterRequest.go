package types

type ProviderRegisterRequest struct {
    Name             string `json:"name" bson:"name" example:"New Provider" description:"Provider Name"`
    Fee              uint   `json:"fee" bson:"fee" example:"100" description:"Fee in wei"`
    QuoteExpiration  uint   `json:"quoteExpiration" bson:"quoteExpiration" example:"100" description:"Quote expiration in seconds"`
    AcceptedQuoteExpiration uint `json:"acceptedQuoteExpiration" bson:"acceptedQuoteExpiration"  example:"100" description:"Accepted quote expiration in seconds"`
    MinTransactionValue uint64 `json:"minTransactionValue" bson:"minTransactionValue" example:"100" description:"Minimum transaction value in wei"`
    MaxTransactionValue uint64 `json:"maxTransactionValue" bson:"maxTransactionValue" example:"100" description:"Maximum transaction value in wei"`
    ApiBaseUrl       string `json:"apiBaseUrl" bson:"apiBaseUrl" example:"https://api.example.com" description:"API base URL"`
    Status           bool   `json:"status" bson:"status" example:"true" description:"Provider status"`
}

