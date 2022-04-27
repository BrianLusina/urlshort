package quotes

type QuoteDto struct {
	Author string `json:"author" validate:"required" binding:"required"`
	Quote  string `json:"quote" validate:"required" binding:"required"`
}

type QuoteResponseDto struct {
	Identifier string `json:"id"`
	QuoteDto
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type QuoteResponsesDto []QuoteDto
