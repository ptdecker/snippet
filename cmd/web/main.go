/*
 * Basic web site project based upon "Lets Go" book
 */

//TODO: Add password reset
//TODO:  Add a way to contact suppport
//TODO: Add a way to cancel account and delete data
//TODO: Add captha support
//TODO: Add health check and handling for if database is offline
//TODO: Add a redirect from HTTP to HTTPS
//TODO: Add ability to close server with Ctrl-C in a nice way

package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"ptodd.org/snippetbox/pkg/models"
	"ptodd.org/snippetbox/pkg/models/mysql"
)

// Config retains passed command-line flags
type config struct {
	addr      string
	staticDir string
	dsn       string
	secret    string
}

// Application struct is used for application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	snippets interface { // Interface is used here so both mysql and mock models can be used
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	users interface { // Interface is used here so both mysql and mock models can be used
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
	templateCache map[string]*template.Template
}

// ContextKey is used to define our own unique key for storage and retrieval of user details
type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

// Globals
var cfg *config // App configuration driven from command-line and default parameters

func init() {
	// Retrieve command-line parameters
	cfg = new(config)
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "web:snippet@/snippetbox?parseTime=true", "MySQL data source name")
	flag.StringVar(&cfg.secret, "secret", "2pf1tyu8dT19yjHhuNozkSY67KJnR4lG", "Secret key")
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

	// Initialize a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new session manager
	// NOTE:  SameSiteStrictMode means folks who come to the site by following a URL
	//        link will be treated as unauthenticated.
	// TODO: Review parameters https://godoc.org/github.com/golangcollege/sessions#Session
	session := sessions.New([]byte(cfg.secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	// Initialize application dependencies
	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		templateCache: templateCache,
	}

	// Custom TLS settings
	// TODO: Consider restricting to only support strong cipher suites understanding
	// doing so will reduce the range of supported browsers
	// TODO: Consider restricting TLS version supported to 1.2 and 1.3
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Set up our new http.Server leveraging our leveled logging
	srv := &http.Server{
		Addr:         cfg.addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Launch server
	infoLog.Printf("Starting server on %s\n", cfg.addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
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
