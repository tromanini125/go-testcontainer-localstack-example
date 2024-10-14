package input

import (
	"context"

	"github.com/tromanini125/go-testcontainer-localstack-example/application/domain"
)

type SaveNewCardUseCase interface {
	Execute(ctx context.Context, newCard *domain.Card) error
}
