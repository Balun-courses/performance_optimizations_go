//go:build recursive || optimized_recursive || iterative || optimized_iterative

package api

import (
	"card_shielder/internal/entity"
	"log/slog"
	"net/http"
)

func (c *CardApiImpl) SaveCard(w http.ResponseWriter, r *http.Request) {
	cardNumber := entity.CardNumber(r.PathValue("card_number"))

	logger := c.logger.With(
		slog.String("path", r.URL.String()),
	)

	cardToken, err := c.cardConverter.GetTokenByCardNumber(cardNumber)

	if err != nil {
		logger.LogAttrs(
			r.Context(),
			slog.LevelWarn,
			"can not create card token",
			slog.Any("error", err),
		)

		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	card := entity.Card{
		CardNumber: cardNumber,
		CardToken:  cardToken,
	}

	err = c.cardRepository.SaveCard(card)

	if err != nil {
		logger.LogAttrs(
			r.Context(),
			slog.LevelError,
			"can not save card info",
			slog.Any("error", err),
		)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = w.Write([]byte(cardToken))

	if err != nil {
		logger.LogAttrs(
			r.Context(),
			slog.LevelError,
			"can not write response ammo_template.txt",
			slog.Any("error", err),
		)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	logger.Info("OK")

	return
}
