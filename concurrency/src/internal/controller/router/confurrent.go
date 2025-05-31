package router

import "net/http"

type Concurrent struct{}

func (c *Concurrent) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
}
func (c *Concurrent) Check(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
}
