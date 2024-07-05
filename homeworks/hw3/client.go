package hw3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
)

type RemoteMethodInvocationClient interface {
	Invoke(ctx context.Context, url *url.URL, rmiMethod string, data []any, response any) error
}

var _ RemoteMethodInvocationClient = (*RmiClient)(nil)

type RmiClient struct {
	client *http.Client
}

func NewRmiClient(client *http.Client) *RmiClient {
	return &RmiClient{
		client: client,
	}
}

func (r *RmiClient) Invoke(
	ctx context.Context,
	address *url.URL,
	rmiMethod string,
	arguments []any,
	response any,
) error {
	serializedRequest, err := json.Marshal(arguments)

	if err != nil {
		return fmt.Errorf("can not marshal requestBody, error: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		address.String(),
		bytes.NewReader(serializedRequest),
	)

	if err != nil {
		return fmt.Errorf("can not create http requestBody, error: %w", err)
	}

	request.Header.Set(rmiHttpHeader, rmiMethod)
	request.Header.Set(globalIdHttpHeader, strconv.Itoa(rand.Int()))
	request.Header.Set("Content-Type", "application/json")

	rawResponse, err := r.client.Do(request)

	if err != nil {
		return fmt.Errorf("http client error: %w", err)
	}

	defer func() {
		err = rawResponse.Body.Close()

		if err != nil {
			return
		}
	}()

	if rawResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("got error code %d from server", rawResponse.StatusCode)
	}

	err = json.NewDecoder(rawResponse.Body).Decode(response)

	if err != nil {
		return fmt.Errorf("got response deserialization error: %w", err)
	}

	return nil
}
