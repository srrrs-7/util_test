package response

import "api/util/utilhttp"

type ErrorRes struct {
	Meg string `json:"message"`
}

func (e *ErrorRes) ResponseJson() ([]byte, error) {
	j, err := utilhttp.Json(e)
	if err != nil {
		return nil, err
	}
	return j, nil
}
