package repository

import (
	"card_shielder/internal/entity"
	"fmt"
	"sync"
)

var _ CardRepository = (*CardRepositoryImpl)(nil)

type CardRepositoryImpl struct {
	data map[entity.CardNumber]*entity.Card
	mx   sync.RWMutex
}

func NewCardRepository() CardRepository {
	return &CardRepositoryImpl{
		data: make(map[entity.CardNumber]*entity.Card),
	}
}

func (c *CardRepositoryImpl) GetCardByCardNumber(cardNumber entity.CardNumber) (*entity.Card, error) {
	c.mx.RLock() // Need example with Lock for read operation
	defer c.mx.RUnlock()

	if value, ok := c.data[cardNumber]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("card %s not found", cardNumber)
}

func (c *CardRepositoryImpl) SaveCard(card entity.Card) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.data[card.CardNumber]; ok {
		return fmt.Errorf("card %s already exists", card.CardNumber)
	}

	c.data[card.CardNumber] = &card

	return nil
}
