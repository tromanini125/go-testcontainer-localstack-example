package service

import (
	"context"
	"log"

	"github.com/tromanini125/go-testcontainer-localstack-example/application/domain"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/input"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/output"
)

type cardService struct {
	repository output.CardPersister
}

func NewCardService(repository output.CardPersister) input.SaveNewCardUseCase {
	return &cardService{
		repository: repository,
	}
}

func (c *cardService) Execute(ctx context.Context, newCard *domain.Card) error {
	log.Default().Printf("Save new Card %s", newCard.CardNumber)
	err := c.repository.CreateCard(ctx, newCard)
	if err != nil {
		return err
	}
	return nil
}
