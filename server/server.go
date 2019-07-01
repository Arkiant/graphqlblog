package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/arkiant/graphqlblog"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query",
		handler.GraphQL(
			graphqlblog.NewExecutableSchema(
				graphqlblog.Config{
					Resolvers: &graphqlblog.Resolver{},
				},
			),
		),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
