package main

import (
	"log"
	"os"
	"welcomeProject/internal/config"
	"welcomeProject/internal/server"
)

func main() {
	dbConfig := config.LoadDBConfig()
	s := server.NewServer(dbConfig)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(s.Start(":" + port))
}
