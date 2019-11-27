package models

import (
	"errors"
	"time"
)

// ErrNoRecord defines a customer error for no matching record
var ErrNoRecord = errors.New("models: no matching record found")

// Snippet defines the model for the Snippet table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
