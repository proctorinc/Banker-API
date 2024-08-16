package main

import (
	"fmt"

	"github.com/proctorinc/banker/router"
)

const PORT = 3000

func main() {
	fmt.Println("Running!")
	router := router.NewRouter()
	router.Run(fmt.Sprintf("localhost:%d", PORT))
}
