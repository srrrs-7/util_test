package main

import (
	"graphql/route"
)

func main() {
	if err := route.NewServer(); err != nil {
		panic(err)
	}
}
