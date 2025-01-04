package main

import (
	"github.com/joho/godotenv"
	"log"
	"roman/cmd"
)

func main() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cmd.Start()
}
