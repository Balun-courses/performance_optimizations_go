package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

const (
	smallTag = "small"
	midTag   = "mid"
	highTag  = "high"
)

func getTag(n int) string {
	if n <= 30 {
		return smallTag
	}

	if n <= 10_000 {
		return midTag
	}

	return highTag
}

func generateRequest(method, url, tag, body string) string {
	return strings.Join([]string{method, url, tag, body}, "||") + "\n"
}

func main() {
	f, err := os.Create("/Users/igorwalther/GolandProjects/performance_optimizations_go/open_lessons/card_shielder/pprof/ammo_template.txt")

	if err != nil {
		panic(err)
	}

	defer func() {
		err = f.Close()

		if err != nil {
			log.Println(err)
		}
	}()

	tokenPool := make([]string, 0)
	cardPool := make(map[string]struct{})

	for range 100_000 {
		partition := rand.N(101)

		if partition <= 50 || len(tokenPool) == 0 {
			cardNumber := generateCardNumber(cardPool)
			cardPool[cardNumber] = struct{}{}

			_, err = f.WriteString(generateRequest(http.MethodPost, "/v1/save_card/"+cardNumber, "save", ""))

			if err != nil {
				panic(err)
			}

			token, err := getTokenByCardNumber(cardNumber)

			if err != nil {
				panic(err)
			}

			tokenPool = append(tokenPool, token)
			continue
		}

		cardToken := tokenPool[len(tokenPool)-1]
		tokenPool = tokenPool[:len(tokenPool)-1]

		_, err = f.WriteString(generateRequest(http.MethodGet, "/v1/get_card_by_token/"+cardToken, "get", ""))
	}
}

func getTokenByCardNumber(cardNumber string) (string, error) {
	result := &strings.Builder{}

	for i := 0; i < len(cardNumber); i += 6 {
		if i != 0 {
			result.WriteString("_")
		}

		cur, err := strconv.ParseInt(cardNumber[i:i+6], 10, 64)

		if err != nil {
			return "", err
		}

		result.Write([]byte(strconv.FormatUint(fibonacci(int(cur)), 10)))
	}

	return result.String(), nil
}

func fibonacci(n int) uint64 {
	var (
		prev uint64
		cur  uint64
	)

	prev, cur = 1, 1

	for i := 2; i < n; i++ {
		prev, cur = cur, prev+cur
	}

	return cur
}

func generateCardNumber(pool map[string]struct{}) string {
	for {
		cardNumber := make([]byte, 18)

		for i := 0; i < 18; i++ {
			cardNumber[i] = byte((rand.N(9)+1)%10 + '0')
		}

		result := unsafe.String(unsafe.SliceData(cardNumber), len(cardNumber))

		if _, ok := pool[result]; ok {
			continue
		}

		return result
	}
}
