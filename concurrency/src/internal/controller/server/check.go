package server

import "net/http"

type Check struct{}

func (c *Check) Check(w http.ResponseWriter, r *http.Request) {}
