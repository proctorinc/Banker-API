package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/proctorinc/banker/internal/db"
	"github.com/proctorinc/banker/internal/graphql"
)

const defaultPort = "8080"

func main() {
	fmt.Println("Connecting to database..")
	conn, err := db.Open("dbname=chase-data sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	repo := db.NewRepository(conn)

	mux := http.NewServeMux()
	mux.Handle("/", graphql.NewPlaygroundHandler("/query"))
	mux.Handle("/query", graphql.NewHandler(repo))

	port := ":8080"
	fmt.Fprintf(os.Stdout, "🚀 Server ready at http://localhost%s\n", port)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(port, mux))

	// conn, err := db.Open("dbname=gqlgen_sqlc_example_db sslmode=disable")
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()

	// // initialize the repository
	// // repo := db.NewRepository(conn)

	// // run the server
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = defaultPort
	// }

	// srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}}))

	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
}
