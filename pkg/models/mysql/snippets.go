package mysql

import (
	"database/sql"

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
	return nil, nil
}

// Latest returns the 10 most recently created snippits
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
