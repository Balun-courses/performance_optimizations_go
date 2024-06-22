package gormi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type service struct{}

func (t *service) ZeroArity() string {
	return "Zero arity"
}

func (t *service) OnlyContext(ctx context.Context) string {
	return "Only context"
}

func (t *service) WithoutContext(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

type PointerRequest struct {
	Value  string
	Nested struct {
		Nested bool
	}
	Float float64
	Int   int
}

type PointerResponse struct {
	Stat string
}

func (p PointerRequest) Stat() string {
	return fmt.Sprintf("Value: %s, Nested: %v, Float: %v, Int: %v",
		p.Value, p.Nested, p.Float, p.Int)
}

func (t *service) Pointer0(ctx context.Context, p PointerRequest) PointerResponse {
	return PointerResponse{
		Stat: p.Stat(),
	}
}

func (t *service) Pointer1(p *PointerRequest) *PointerResponse {
	return &PointerResponse{
		Stat: p.Stat(),
	}
}

func (t *service) Pointer2(ctx context.Context, p **PointerRequest) **PointerResponse {
	response := &PointerResponse{
		Stat: (*p).Stat(),
	}

	return &response
}

func (t *service) ErrorMethod(ctx context.Context) (string, error) {
	return "hello", fmt.Errorf("hello")
}

func (t *service) Variadic(ctx context.Context, s string, a ...float64) string {
	return s + strconv.Itoa(len(a))
}

func TestGoRmi(t *testing.T) {
	t.Parallel()

	stubProvider := NewRmiStubProvider()

	server := httptest.NewServer(stubProvider.CreateObjectStub(&service{}))
	t.Cleanup(func() {
		server.Close()
	})

	serverUrl, err := url.Parse(server.URL)
	require.NoError(t, err)

	client := NewRmiClient(http.DefaultClient)

	ctx := context.Background()

	type testCase struct {
		TestName  string
		RmiMethod string
		Request   []any
		Response  string
		Error     require.ErrorAssertionFunc
	}

	trivialCases := []testCase{
		{
			TestName:  "Zero arity",
			RmiMethod: "ZeroArity",
			Response:  "Zero arity",
			Error:     require.NoError,
		},
		{
			TestName:  "Only context",
			RmiMethod: "OnlyContext",
			Response:  "Only context",
			Error:     require.NoError,
		},
		{
			TestName:  "Without context",
			RmiMethod: "WithoutContext",
			Request:   []any{12.12345},
			Response:  "12.12",
			Error:     require.NoError,
		},
		{
			TestName:  "Error method",
			RmiMethod: "ErrorMethod",
			Error:     require.Error,
		},
		{
			TestName:  "Test variadic",
			RmiMethod: "Variadic",
			Request:   []any{"123", 1, 2, 3, 4, 5},
			Response:  "1235",
			Error:     require.NoError,
		},
	}

	for _, tt := range trivialCases {
		// 1.21-
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			var response string
			err = client.Invoke(ctx, serverUrl, tt.RmiMethod, tt.Request, &response)

			tt.Error(t, err)
			require.Equal(t, tt.Response, response)
		})
	}
}

func TestPointers(t *testing.T) {
	t.Parallel()

	stubProvider := NewRmiStubProvider()

	server := httptest.NewServer(stubProvider.CreateObjectStub(&service{}))
	t.Cleanup(func() {
		server.Close()
	})

	serverUrl, err := url.Parse(server.URL)
	require.NoError(t, err)

	client := NewRmiClient(http.DefaultClient)

	ctx := context.Background()

	request := PointerRequest{
		Value:  "123",
		Nested: struct{ Nested bool }{Nested: true},
		Float:  3.1415,
		Int:    12345,
	}

	response := PointerResponse{
		Stat: request.Stat(),
	}

	type testCase struct {
		TestName  string
		RmiMethod string
	}

	testCases := []testCase{
		{
			TestName:  "Pointer 0",
			RmiMethod: "Pointer0",
		},
		{
			TestName:  "Pointer 1",
			RmiMethod: "Pointer1",
		},
		{
			TestName:  "Pointer 2",
			RmiMethod: "Pointer2",
		},
	}

	for _, tt := range testCases {
		// 1.21-
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			var result PointerResponse
			err = client.Invoke(ctx, serverUrl, tt.RmiMethod, []any{request}, &result)

			require.NoError(t, err)
			require.Equal(t, response.Stat, result.Stat)
		})
	}
}
