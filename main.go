package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kocierik/k8s-to-diagram/pkg/api"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	startServer()

}

func startServer() {
	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = "8080"
	}

	fmt.Printf("Starting server on port %s...\n", envPort)
	port, err := strconv.Atoi(envPort)

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	if err := api.StartServer(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
