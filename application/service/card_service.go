package service

import (
	"log"

	"github.com/tromanini125/go-testcontainer-localstack-example/application/domain"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/input"
)

type cardService struct{}

func NewCardService() input.SaveNewCardUseCase {
	return &cardService{}
}

func (c *cardService) Execute(newCard *domain.Card) error {
	log.Default().Printf("Save new Card %s", newCard.CardNumber)
	return nil
}
