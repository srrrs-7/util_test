package utilhttp

import (
	"encoding/json"
	"net/http"
)

func GetRequestJsonBody[T any](r *http.Request) (T, error) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return body, err
	}
	return body, nil
}
