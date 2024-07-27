//go:build optimized_iterative

package conversion

import (
	"card_shielder/internal/entity"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testCard entity.CardNumber = "111122223333444455"
)

func TestCardConversion(t *testing.T) {
	t.Parallel()

	converter := NewCardConverted()

	token, err := converter.GetTokenByCardNumber(testCard)
	require.NoError(t, err)

	cardNumber, err := converter.GetCardNumberByToken(token)
	require.NoError(t, err)

	require.Equal(t, cardNumber, testCard)
}
