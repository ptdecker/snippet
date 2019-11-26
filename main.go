package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	// Extract expected 'id' parameter from query string
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// createSnippet handler
func createSnippet(w http.ResponseWriter, r *http.Request) {

	// Guard against non-POST calls to this endpoint
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

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
