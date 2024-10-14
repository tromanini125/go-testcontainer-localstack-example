package model

type CardCreatedEvent struct {
	CardId         int64  `json:"cardId"`
	CardHolderName string `json:"cardHolderName"`
	CardNumber     string `json:"cardNumber"`
	CVV            string `json:"cvv"`
	ExpiryDate     string `json:"expiryDate"`
}
