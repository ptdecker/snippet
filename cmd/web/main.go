package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

	// Set-up leveled logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize new server mux
	mux := http.NewServeMux()

	// Register home page route
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Create a file server to serve static content
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Set up our new http.Server leveraging our leveled logging
	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Launch server
	infoLog.Printf("Starting server on %s\n", cfg.Addr)
	errorLog.Fatal(srv.ListenAndServe())
}
