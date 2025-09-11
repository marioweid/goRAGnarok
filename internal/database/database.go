package database

import (
	"context"
	"database/sql"
	"fmt"
)

// SearchResult represents a row returned from similarity search.
// Imported from models.go

// SimilaritySearch finds the top N most similar vectors in the database.
func SimilaritySearch(ctx context.Context, db *sql.DB, input []float32, topN int) ([]SearchResult, error) {
	// Convert Go slice to PostgreSQL array literal
	vectorStr := fmt.Sprintf("ARRAY[%s]", float32SliceToString(input)) // TODO  Continue here

	query := `
        SELECT id, vector, (vector <-> $1::vector) AS score, title, content
        FROM articles
        ORDER BY vector <-> $1::vector
        LIMIT $2
    `
	rows, err := db.QueryContext(ctx, query, vectorStr, topN)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		var pgVector []byte
		if err := rows.Scan(&r.ID, &pgVector, &r.Embedding, &r.Title, &r.Content); err != nil {
			return nil, err
		}
		r.Embedding = parsePGVector(pgVector)
		results = append(results, r)
	}
	return results, rows.Err()
}

// float32SliceToString converts a []float32 to a comma-separated string.
func float32SliceToString(vec []float32) string {
	out := ""
	for i, v := range vec {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%f", v)
	}
	return out
}

// parsePGVector parses PostgreSQL vector byte slice to []float32.
// Implement this based on your pgvector driver or use a library.
// For now, returns nil.
func parsePGVector(pg []byte) []float32 {
	// TODO: Implement parsing if needed
	return nil
}
