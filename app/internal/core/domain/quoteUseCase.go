package domain

import (
	"quote/api/app/internal/core/contracts"
	"quote/api/app/internal/core/domain/entity"
)

type QuotesUseCase struct {
	quoteRepo contracts.QuoteRepository
}

func NewQuotesUseCase(quoteRepo contracts.QuoteRepository) *QuotesUseCase {
	return &QuotesUseCase{quoteRepo: quoteRepo}
}

func (q *QuotesUseCase) CreateQuote(author, quote string) (*entity.Quote, error) {
	newQuote, err := entity.NewQuote(author, quote)
	if err != nil {
		return nil, err
	}

	_, err = q.quoteRepo.Save(*newQuote)
	if err != nil {
		return nil, err
	}
	return newQuote, nil
}

func (q *QuotesUseCase) GetAllQuotes() ([]entity.Quote, error) {
	allQuotes, err := q.quoteRepo.GetAllQuotes()
	if err != nil {
		return nil, err
	}
	return allQuotes, nil
}
