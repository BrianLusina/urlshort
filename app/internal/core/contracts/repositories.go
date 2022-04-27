package contracts

import (
	"quote/api/app/internal/core/domain/entity"
)

type QuoteRepository interface {
	Save(entity.Quote) (entity.Quote, error)
	GetAllQuotes() ([]entity.Quote, error)
	GetQuote(id string) (entity.Quote, error)
	UpdateQuote(entity.Quote) (entity.Quote, error)
}
