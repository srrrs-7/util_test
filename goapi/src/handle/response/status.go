package response

import "api/util/utilhttp"

type StatusRes struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func (s *StatusRes) ResponseJson() ([]byte, error) {
	j, err := utilhttp.Json(s)
	if err != nil {
		return nil, err
	}
	return j, nil
}
