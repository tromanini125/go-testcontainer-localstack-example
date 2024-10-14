package persistence

import (
	"context"

	"github.com/tromanini125/go-testcontainer-localstack-example/adapter/output/persistence/model"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/domain"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/output"
	databaseconfig "github.com/tromanini125/go-testcontainer-localstack-example/configuration/database_config"
)

type persistence struct {
}

func NewCardRepository() output.CardPersister {
	return &persistence{}
}

func (p *persistence) CreateCard(ctx context.Context, cardDomain *domain.Card) error {
	db, err := databaseconfig.GetConnection()
	if err != nil {
		return err
	}
	cardEntity := mapDomainToEntity(cardDomain)

	insertedData := db.Create(cardEntity)
	if insertedData.Error != nil {
		return insertedData.Error
	}

	insertedId := int64(cardEntity.ID)
	cardDomain.CardId = &insertedId

	return nil
}

func (p *persistence) FindCardByNumber(ctx context.Context, cardNumber string) (*domain.Card, error) {
	db, err := databaseconfig.GetConnection()
	if err != nil {
		return nil, err
	}
	var cardEntity *model.Card
	db.Where("card_number = ?", cardNumber).First(&cardEntity)
	// SELECT * FROM cards WHERE card_number = 'cardNumber' ORDER BY id LIMIT 1;
	return mapEntityToDomain(cardEntity), nil
}

func mapDomainToEntity(domain *domain.Card) *model.Card {
	return &model.Card{
		CardHolderName: domain.CardHolderName,
		CardNumber:     domain.CardNumber,
		CVV:            domain.CVV,
		ExpiryDate:     domain.ExpiryDate,
	}
}

func mapEntityToDomain(entity *model.Card) *domain.Card {
	Id := int64(entity.ID)
	return &domain.Card{
		CardId:         &Id,
		CardHolderName: entity.CardHolderName,
		CardNumber:     entity.CardNumber,
		CVV:            entity.CVV,
		ExpiryDate:     entity.ExpiryDate,
	}
}
