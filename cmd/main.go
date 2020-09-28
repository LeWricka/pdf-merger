package main

import (
	pdfmerger "climbingplan.com/pdfmerger"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/joho/godotenv"
	"log"
	"os"
)



func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	funcframework.RegisterHTTPFunction("/", pdfmerger.Merge)
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("framework.Start: %v\n", err)
	}
}
