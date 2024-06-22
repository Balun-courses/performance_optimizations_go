package gormi

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
)

type RemoteMethodStubProvider interface {
	CreateObjectStub(serverObject any) http.HandlerFunc
}

var _ RemoteMethodStubProvider = (*RmiStubProvider)(nil)

type RmiStubProvider struct {
}

func NewRmiStubProvider() *RmiStubProvider {
	return &RmiStubProvider{}
}

func (r *RmiStubProvider) CreateObjectStub(serverObject any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		method := reflect.ValueOf(serverObject).
			MethodByName(r.Header.Get(rmiHttpHeader))
		methodSignature := method.Type()

		if !method.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		arguments := make([]any, 0)
		err := json.NewDecoder(r.Body).Decode(&arguments)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		methodArguments := make([]reflect.Value, 0, len(arguments))
		contextType := reflect.TypeOf((*context.Context)(nil)).Elem()

		methodArity := methodSignature.NumIn()

		if methodArity > 0 && methodSignature.In(0).Implements(contextType) {
			methodArguments = append(methodArguments, reflect.ValueOf(r.Context()))
		}

		if len(methodArguments)+len(arguments) < methodArity {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !methodSignature.IsVariadic() && len(methodArguments)+len(arguments) != methodArity {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		argumentPointer := len(methodArguments)

		for i := 0; i < len(arguments); i++ {
			methodArgument := methodSignature.In(min(argumentPointer, methodArity-1))
			argumentPointer++

			if methodArgument.Kind() != reflect.Pointer && methodArgument.Kind() != reflect.Struct {
				methodArguments = append(methodArguments, reflect.ValueOf(arguments[i]))
				continue
			}

			serializedArgument, err := json.Marshal(arguments[i])

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			argument := reflect.New(methodArgument).Interface()
			err = json.Unmarshal(serializedArgument, argument)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			methodArguments = append(methodArguments, reflect.ValueOf(argument).Elem())
		}

		returned := method.Call(methodArguments)

		var (
			result      any
			methodError any
		)

		if len(returned) > 0 {
			result = returned[0].Interface()
		}

		if len(returned) > 1 {
			methodError = returned[1].Interface()
		}
		// ignore other arguments

		if methodError != nil {
			http.Error(w, methodError.(error).Error(), http.StatusExpectationFailed)
			return
		}

		err = json.NewEncoder(w).Encode(result)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
