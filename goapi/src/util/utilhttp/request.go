package utilhttp

import (
	"encoding/json"
	"io"
	"net/http"
)

func RequestUrlParam(req *http.Request, key string) string {
	return req.URL.Query().Get(key)
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
