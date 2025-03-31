package main

import (
	"log"
	"os"
	"welcomeProject/internal/server"
)

func main() {
	s := server.NewServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(s.Start(":" + port))
}
