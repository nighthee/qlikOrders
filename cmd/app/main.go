package main

import (
	"fmt"
	"log"
	"qlikOrders/internal/collections"
	"qlikOrders/internal/server"
)

func main() {
	// Create a new instance of collections
	orderCollections := &collections.OrderCollection{}

	// Inject the collections
	srv := server.NewServer(orderCollections)
	fmt.Println("Starting server on port 8080...")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
