package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/abuabdillatief/sample/db"
	"github.com/abuabdillatief/sample/handler"
	"github.com/joho/godotenv"
)


func init() {
	godotenv.Load()
}

func main() {
	// Initialize MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := db.ConnectMongo()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect MongoDB client: %v", err)
		}
	}()

	http.HandleFunc("/users", handler.GetUsersHandler)
	http.HandleFunc("/user", handler.CreateUserHandler)

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

