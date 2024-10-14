package model

type CardCreatedEvent struct {
	CardHolderName string `json:"cardHolderName"`
	CardNumber     string `json:"cardNumber"`
	CVV            string `json:"cvv"`
	ExpiryDate     string `json:"expiryDate"`
}
