package main

import (
	"log"
	"net/http"
)

// Home page handler
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func main() {

	// Initialize new server mux
	mux := http.NewServeMux()

	// Register home page route
	mux.HandleFunc("/", home)

	// Launch server
	log.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
