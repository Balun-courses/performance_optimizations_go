//go:build recursive || optimized_recursive || iterative || optimized_iterative

package api

import (
	"card_shielder/internal/entity"
	"log/slog"
	"net/http"
	"unsafe"
)

func (c *CardApiImpl) GetCardByToken(w http.ResponseWriter, r *http.Request) {
	cardToken := entity.CardToken(r.PathValue("card_token"))

	logger := c.logger.With(
		slog.String("path", r.URL.String()),
		slog.Any("card_token", cardToken),
	)

	cardNumber, err := c.cardConverter.GetCardNumberByToken(cardToken)

	if err != nil {
		logger.LogAttrs(
			r.Context(),
			slog.LevelWarn,
			"can not get card number by token",
			slog.Any("error", err),
		)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	card, err := c.cardRepository.GetCardByCardNumber(cardNumber)

	if err != nil {
		logger.LogAttrs(
			r.Context(),
			slog.LevelError,
			"can not get card info",
			slog.Any("error", err),
		)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = w.Write(unsafe.Slice(unsafe.StringData(string(card.CardNumber)), len(cardNumber)))

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
}
