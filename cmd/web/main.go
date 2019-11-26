/*
 * Basic web site project based upon "Lets Go" book
 */

package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Config retains passed command-line flags
type config struct {
	addr      string
	staticDir string
	dsn       string
}

// Application struct is used for application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

var cfg *config

func init() {
	// Retrieve command-line parameters
	cfg = new(config)
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "web:snippet@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
}

func main() {

	// Set-up leveled logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize database connection pool (MySQL)
	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize application dependencies
	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	// Set up our new http.Server leveraging our leveled logging
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Launch server
	infoLog.Printf("Starting server on %s\n", cfg.addr)
	errorLog.Fatal(srv.ListenAndServe())
}

// openDB is a wrapper for sql.Open()
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
