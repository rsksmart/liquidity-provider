package types

type GlobalProvider struct {
	Id           uint64 `json:"id" example:"1" description:"Provider Id"`
	Provider     string `json:"provider" example:"0x0" description:"Provider Address"`
	Name         string `json:"name" example:"New Provider" description:"Provider Name"`
	ApiBaseUrl   string `json:"apiBaseUrl" example:"https://api.example.com" description:"API base URL"`
	Status       bool   `json:"status" example:"true" description:"Provider status"`
	ProviderType string `json:"providerType" example:"pegin" description:"Provider type"`
}
