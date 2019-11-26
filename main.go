package main

import (
	"log"
	"net/http"
)

// Home page handler
func home(w http.ResponseWriter, r *http.Request) {

	// Site-wide 404 handler
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

// showSnippet handler
func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

// createSnippet handler
func createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}

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
