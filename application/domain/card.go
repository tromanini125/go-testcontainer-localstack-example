package domain

type Card struct {
	CardId         *int64 `json:"id"`
	CardHolderName string `json:"cardHolderName"`
	CardNumber     string `json:"cardNumber"`
	CVV            string `json:"cvv"`
	ExpiryDate     string `json:"expiryDate"`
}
