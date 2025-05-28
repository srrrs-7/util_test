package server

import "net/http"

type Create struct{}

func (c *Create) Create(w http.ResponseWriter, r *http.Request) {}
