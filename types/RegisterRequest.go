package types

type ProviderRegisterRequest struct {
	Name         string `json:"name" bson:"name" example:"New Provider" description:"Provider Name"`
	Fee          uint   `json:"fee" bson:"fee" example:"100" description:"Fee in wei"`
	ApiBaseUrl   string `json:"apiBaseUrl" bson:"apiBaseUrl" example:"https://api.example.com" description:"API base URL"`
	ProviderType string `json:"providerType" bson:"providerType" example:"pegin" description:"Provider type must be \"pegin\", \"pegout\" or \"both\""`
	Status       bool   `json:"status" bson:"status" example:"true" description:"Provider status"`
}
