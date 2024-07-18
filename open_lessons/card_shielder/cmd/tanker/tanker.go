package main

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"sync"
	"time"
	"unsafe"
)

var c *http.Client

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(size int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, size),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func main() {
	wg := &sync.WaitGroup{}
	semaphore := NewSemaphore(15)

	client := http.DefaultTransport.(*http.Transport).Clone()
	client.DisableKeepAlives = true // TODO
	client.IdleConnTimeout = time.Millisecond * 300
	client.MaxIdleConnsPerHost = client.MaxIdleConns

	c = &http.Client{
		Transport: client,
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		start := time.Now()
		go func() {
			defer wg.Done()
			for time.Since(start) < time.Minute {
				func() {
					semaphore.Acquire()
					defer semaphore.Release()

					cardNumber := generateCardNumber()

					cardToken, err := doSaveCard(cardNumber)

					if err != nil {
						log.Println("doSaveCard error")
						return
					}

					outputCardNumber, err := doGetCardByToken(cardToken)

					if err != nil {
						log.Println("doGetCardByToken error")
						return
					}

					log.Println("CARD: " + outputCardNumber)

					time.Sleep(time.Millisecond * 300)
				}()
			}
		}()
	}

	wg.Wait()
}

func generateCardNumber() string {
	cardNumber := make([]byte, 18)

	for i := 0; i < 18; i++ {
		cardNumber[i] = byte(rand.N(9) + 1 + '0')
	}

	return unsafe.String(unsafe.SliceData(cardNumber), len(cardNumber))
}

func doSaveCard(cardNumber string) (string, error) {
	base, err := url.Parse("http://127.0.0.1:8801/v1/save_card")

	if err != nil {
		panic(err)
	}

	requestUrl := base.JoinPath(cardNumber)

	request := &http.Request{
		Method: http.MethodPost,
		URL:    requestUrl,
	}

	response, err := c.Do(request)

	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusUnprocessableEntity {
		return "", fmt.Errorf("banned card")
	}

	data, err := io.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	return unsafe.String(unsafe.SliceData(data), len(data)), nil
}

func doGetCardByToken(cardToken string) (string, error) {
	base, err := url.Parse("http://127.0.0.1:8801/v1/get_card_by_token")

	if err != nil {
		panic(err)
	}

	requestUrl := base.JoinPath(cardToken)

	request := &http.Request{
		Method: http.MethodGet,
		URL:    requestUrl,
	}

	response, err := c.Do(request)

	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	return unsafe.String(unsafe.SliceData(data), len(data)), nil
}
