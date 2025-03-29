package main

import (
	"log"

	api "github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd"
)

func main() {
	if err := api.Run(); err != nil {
		log.Fatalf("Failed to start server :%v", err)
	}
}
