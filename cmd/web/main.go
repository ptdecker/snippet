package main

import (
	"flag"
	"log"
	"net/http"
)

// Config retains passed command-line flags
type Config struct {
	Addr      string
	StaticDir string
}

var cfg *Config

func init() {
	// Retrieve command-line parameters
	cfg = new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()
}

func main() {

	// Initialize new server mux
	mux := http.NewServeMux()

	// Register home page route
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Create a file server to serve static content
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Launch server
	log.Printf("Starting server on %s\n", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, mux))
}
