package database

// SearchResult represents a row returned from similarity search.
//
// ID: Unique identifier for the article.
// Vector: Embedding vector stored in the database.
// Score: Similarity score (lower is more similar).
// Title, Content: Example additional fields.
type SearchResult struct {
	ID        int       `json:"id" db:"id"`
	URL       string    `json:"url" db:"url"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Embedding []float32 `json:"embedding" db:"embedding"`
}
