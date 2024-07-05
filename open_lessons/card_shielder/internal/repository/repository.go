package repository

import (
	"card_shielder/internal/entity"
)

type CardRepository interface {
	GetCardByCardNumber(cardNumber entity.CardNumber) (*entity.Card, error)
	SaveCard(card entity.Card) error
}
