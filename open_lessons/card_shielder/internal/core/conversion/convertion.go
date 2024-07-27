//go:build recursive || optimized_recursive || iterative || optimized_iterative

package conversion

import (
	"card_shielder/internal/entity"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const BatchSize = 6
const delimiter = "_"

type CardConverter interface {
	GetTokenByCardNumber(cardNumber entity.CardNumber) (entity.CardToken, error)
	GetCardNumberByToken(cardToken entity.CardToken) (entity.CardNumber, error)
}

var _ CardConverter = (*CardConverterImpl)(nil)

type CardConverterImpl struct {
}

func NewCardConverted() CardConverter {
	return &CardConverterImpl{}
}

func (c *CardConverterImpl) GetTokenByCardNumber(cardNumber entity.CardNumber) (entity.CardToken, error) {
	result := &strings.Builder{}

	for i := 0; i < len(cardNumber); i += BatchSize {
		if i != 0 {
			result.WriteString(delimiter)
		}

		cur, err := strconv.ParseInt(string(cardNumber[i:i+BatchSize]), 10, 64)

		if err != nil {
			return "", err
		}

		result.WriteString(strconv.FormatUint(Fibonacci(int(cur)), 10))
	}

	return entity.CardToken(result.String()), nil
}

func (c *CardConverterImpl) GetCardNumberByToken(cardToken entity.CardToken) (entity.CardNumber, error) {
	cardNumber := &strings.Builder{}
	parts := strings.Split(string(cardToken), delimiter)

	for i := 0; i < len(parts); i++ {
		partNumber, err := strconv.ParseUint(parts[i], 10, 64)

		if err != nil {
			return "", fmt.Errorf("can not parse part of token %s, error - %w", parts[i], err)
		}

		number := optimizedGetFibSerialByNumber(partNumber)

		cardNumber.WriteString(strconv.FormatInt(int64(number), 10))
	}

	return entity.CardNumber(cardNumber.String()), nil
}

func getFibSerialByNumber(fibNumber uint64) (int, error) {
	phi := (1 + math.Sqrt(5)) / 2
	i := int(math.Log(float64(fibNumber)*math.Sqrt(5)+0.5)/math.Log(phi)) - 5

	for ; i <= math.MaxInt; i++ {
		if Fibonacci(i) == fibNumber {
			return i, nil
		}
	}

	return 0, fmt.Errorf("not found")
}

func optimizedGetFibSerialByNumber(fibNumber uint64) int {
	var (
		prev uint64
		cur  uint64
	)

	prev, cur = 1, 1

	for i := 2; i < math.MaxInt; i++ {
		if cur == fibNumber {
			return i
		}

		prev, cur = cur, prev+cur
	}

	return 1
}
