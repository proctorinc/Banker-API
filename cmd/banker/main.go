package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/proctorinc/banker/internal/router"
)

func main() {
	godotenv.Load(".env")
	router := router.NewRouter()
	router.Run(fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
}
