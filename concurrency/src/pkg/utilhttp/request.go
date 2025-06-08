package utilhttp

import (
	"net/http"

	"github.com/goccy/go-json"
)

func GetRequestJsonBody[T any](r *http.Request) (T, error) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return body, err
	}
	return body, nil
}

func GetRequestQueryParam[T any](r *http.Request, key string) (T, error) {
	var value T
	param := r.URL.Query().Get(key)
	if param == "" {
		return value, http.ErrNoLocation
	}
	if err := json.Unmarshal([]byte(param), &value); err != nil {
		return value, err
	}
	return value, nil
}

func GetRequestPathParam[T any](r *http.Request, key string) (T, error) {
	var value T
	param := r.Context().Value(key)
	if param == nil {
		return value, http.ErrNoLocation
	}
	if err := json.Unmarshal([]byte(param.(string)), &value); err != nil {
		return value, err
	}
	return value, nil
}

func GetRequestHeader[T any](r *http.Request, key string) (T, error) {
	var value T
	header := r.Header.Get(key)
	if header == "" {
		return value, http.ErrNoLocation
	}
	if err := json.Unmarshal([]byte(header), &value); err != nil {
		return value, err
	}
	return value, nil
}

func GetRequestCookie[T any](r *http.Request, key string) (T, error) {
	var value T
	cookie, err := r.Cookie(key)
	if err != nil {
		return value, err
	}
	if cookie.Value == "" {
		return value, http.ErrNoLocation
	}
	if err := json.Unmarshal([]byte(cookie.Value), &value); err != nil {
		return value, err
	}
	return value, nil
}

func GetRequestFormValue[T any](r *http.Request, key string) (T, error) {
	var value T
	if err := r.ParseForm(); err != nil {
		return value, err
	}
	param := r.FormValue(key)
	if param == "" {
		return value, http.ErrNoLocation
	}
	if err := json.Unmarshal([]byte(param), &value); err != nil {
		return value, err
	}
	return value, nil
}
