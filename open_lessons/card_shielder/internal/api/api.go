//go:build recursive || optimized_recursive || iterative || optimized_iterative

package api

import (
	"card_shielder/internal/core/conversion"
	"card_shielder/internal/repository"
	"log/slog"
	"net/http"
)

type CardApi interface {
	GetServeMux() *http.ServeMux
}

var _ CardApi = (*CardApiImpl)(nil)

type CardApiImpl struct {
	logger         *slog.Logger
	cardRepository repository.CardRepository
	cardConverter  conversion.CardConverter
}

func NewCardApi(logger *slog.Logger, cardRepository repository.CardRepository) CardApi {
	return &CardApiImpl{
		logger:         logger,
		cardRepository: cardRepository,
		cardConverter:  conversion.NewCardConverted(),
	}
}

func (c *CardApiImpl) GetServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/save_card/{card_number}", c.SaveCard)
	mux.HandleFunc("GET /v1/get_card_by_token/{card_token}", c.GetCardByToken)

	return mux
}
