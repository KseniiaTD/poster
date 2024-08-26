package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/KseniiaTD/poster/graph"

	"github.com/KseniiaTD/poster/internal/database"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

var withInMemoryDB = flag.Bool("m", false, "run with in-memory database")

func main() {
	err := godotenv.Load("./config/config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.Parse()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := database.New(*withInMemoryDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		DB: db}}))

	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
