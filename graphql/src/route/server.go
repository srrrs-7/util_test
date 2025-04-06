package route

import (
	"graphql/graph"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func NewServer() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{},
			},
		),
	)

	http.Handle("/",
		loggingMiddleware(playground.Handler("GraphQL playground", "/query")),
	)
	http.Handle("/query",
		loggingMiddleware(srv),
	)

	return http.ListenAndServe(":"+port, nil)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request started: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Request completed: %s %s - took %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
