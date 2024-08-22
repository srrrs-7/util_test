package utilhttp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func RequestUrlParam[T comparable](req *http.Request, key string) (T, error) {
	t, ok := req.Context().Value(key).(T)
	if !ok {
		return t, errors.New(key + " not found")
	}
	return t, nil
}

func RequestBody[T any](req *http.Request) T {
	var body T

	b, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		panic(err.Error())
	}

	return body
}
