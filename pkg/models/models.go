package models

import (
	"errors"
	"time"
)

// ErrNoRecord defines a customer error for no matching record
var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

// Snippet defines the model for the Snippet table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// User defines the model for the users table
// TODO: consider chaning bool "Active" to "Deactivated" time stamp
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
