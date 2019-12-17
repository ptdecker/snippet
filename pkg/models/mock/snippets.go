package mock

import (
	"time"

	"ptodd.org/snippetbox/pkg/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

// SnippetModel is a mock structure for the snippet model
type SnippetModel struct{}

// Insert is a mock insert handler
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

// Get is a mock get handler
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// Latest is a mock latest handler
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
