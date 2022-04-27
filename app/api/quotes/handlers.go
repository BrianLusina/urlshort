package quotes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (hdl *quotesRouter) getAllQuotes(ctx *gin.Context) {
	quotes, err := hdl.svc.GetAllQuotes()

	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	response := make([]QuoteResponsesDto, len(quotes))
	for idx, quote := range quotes {
		response[idx] = QuoteResponsesDto{
			QuoteDto{
				Quote:  quote.Quote,
				Author: quote.Author,
			},
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func (hdl *quotesRouter) createQuote(ctx *gin.Context) {
	var request QuoteDto
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	quote, err := hdl.svc.CreateQuote(request.Author, request.Quote)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response := QuoteResponseDto{
		Identifier: quote.ID.String(),
		QuoteDto: QuoteDto{
			Quote:  quote.Quote,
			Author: quote.Author,
		},
		CreatedAt: quote.CreatedAt.Format(time.RFC3339),
		UpdatedAt: quote.UpdatedAt.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, response)
}
