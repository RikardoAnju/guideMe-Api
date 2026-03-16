package main

import (
	"guide-me/internal/config"
	"log"
	"net/http"
	"os"
)

func main() {

	config.ConnectDB()
	config.RunMigrations()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Running 🚀"))
	})

	log.Println("Server running on port:", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}