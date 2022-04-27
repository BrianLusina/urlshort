package entity

import (
	"quote/api/app/pkg"
	"quote/api/app/pkg/identifier"
	"strings"
)

type Quote struct {
	identifier.ID[Quote]
	Author string
	Quote  string
	BaseEntity
}

func NewQuote(author, quote string) (*Quote, error) {
	if len(quote) == 0 {
		return nil, pkg.ErrInvalidQuote
	}

	id := identifier.New[Quote]()

	if author == "" {
		author = "Unknown"
	}

	author = strings.Trim(author, " ")

	return &Quote{
		ID:         id,
		Quote:      quote,
		Author:     author,
		BaseEntity: NewBaseEntity(),
	}, nil
}

func (q Quote) Prefix() string {
	return "quote"
}
