package input

import "github.com/tromanini125/go-testcontainer-localstack-example/application/domain"

type SaveNewCardUseCase interface {
	Execute(newCard *domain.Card) error
}
