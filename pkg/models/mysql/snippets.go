package mysql

import (
	"database/sql"
	"errors"

	"ptodd.org/snippetbox/pkg/models"
)

// SnippetModel wrapps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	// Insert SQL to add a row into the snippets table
	stmt := `INSERT INTO snippets (title, content, created, expires)
	        	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Execute the insert
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Then get the returned ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	// Select SQL to retreive a row from the snippets table
	stmt := `SELECT id, title, content, created, expires
				FROM snippets
				WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Initialize structure to hold the returned data
	s := &models.Snippet{}

	// Use row.Scan() to copy attributes returned to their corresponding fields
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil && errors.Is(err, sql.ErrNoRows) { // No records error
		return nil, models.ErrNoRecord
	}
	if err != nil { // All other errors
		return nil, err
	}

	return s, nil
}

// Latest returns the 10 most recently created snippits
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
