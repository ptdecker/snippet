package main

import (
	"log"
	"net/http"
)

func main() {

	// Initialize new server mux
	mux := http.NewServeMux()

	// Register home page route
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Launch server
	log.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
