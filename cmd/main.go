package main

import (
	"context"
	"graphql-backend/app"
	loaders "graphql-backend/data-loader"
	"graphql-backend/graph"
	http_transport "graphql-backend/pkg/http-transport"
	"graphql-backend/store"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	trans "graphql-backend/transport"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ctx := context.Background()

	jwtKeyPair, err := http_transport.LoadRSAKeys()
	if err != nil {
		panic("failed to load RSA keys: " + err.Error())
	}

	jwtHandler := http_transport.NewJWTHandler(jwtKeyPair)
	authMw := http_transport.AuthMiddleware(jwtHandler)

	repo := store.NewRepo(ctx)
	query := app.NewQuery(repo)
	service := app.NewService(repo, jwtHandler)

	api := trans.NewAPI(query, service)
	c := graph.Config{Resolvers: &graph.Resolver{
		Api: api,
	}}
	c.Directives = graph.DirectiveRoot{
		HasRole:          http_transport.HasRole,
		HasAuthenticated: http_transport.HasAuthenticated,
	}

	srv := handler.New(graph.NewExecutableSchema(c))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Middleware for authentication and data loaders
	handler := authMw(srv)
	handler = loaders.Middleware(handler, repo)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", handler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
