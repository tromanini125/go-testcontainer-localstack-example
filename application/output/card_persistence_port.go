package output

import (
	"context"

	"github.com/tromanini125/go-testcontainer-localstack-example/application/domain"
)

type CardPersister interface {
	CreateCard(ctx context.Context, cardDomain *domain.Card) error
	FindCardByNumber(ctx context.Context, cardNumber string) (*domain.Card, error)
}
