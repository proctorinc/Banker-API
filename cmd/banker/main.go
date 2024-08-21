package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	"github.com/proctorinc/banker/internal/graphql"
)

func main() {
	fmt.Println("Connecting to database..")
	conn, err := db.Open("dbname=chase-data sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	repo := db.NewRepository(conn)

	r := gin.Default()
	r.Use(auth.Middleware(repo))
	r.POST("/query", graphql.GraphqlHandler(repo))
	r.GET("/", graphql.NewPlaygroundHandler())
	r.Run()
}
