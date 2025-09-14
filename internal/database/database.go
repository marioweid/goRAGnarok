package database

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// SearchResult represents a row returned from similarity search.
// Imported from models.go

// SimilaritySearch finds the top N most similar vectors in the database.
func SimilaritySearch(ctx context.Context, db *sql.DB, input []float32, topN int) ([]SearchResult, error) {
	// Convert Go slice to PostgreSQL array literal
	vectorStr := fmt.Sprintf("[%s]", float32SliceToString(input))

	query := `
	    SELECT id, url, title, content, embedding
	    FROM public.doc_sections
	    ORDER BY embedding <-> $1::vector
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
		var pgVector []byte // Postgres returns json string bytes
		if err := rows.Scan(&r.ID, &r.URL, &r.Title, &r.Content, &pgVector); err != nil {
			return nil, err
		}
		r.Embedding, err = parsePGVector(pgVector)
		if err != nil {
			return nil, err
		}
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
func parsePGVector(pg []byte) ([]float32, error) {
	s := string(bytes.TrimSpace(pg))
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	if s == "" {
		return []float32{}, nil
	}

	parts := strings.Split(s, ",")
	vec := make([]float32, len(parts)) // allocate memory for size of list

	// parse every single entry to float32
	for i, p := range parts {
		f, err := strconv.ParseFloat(strings.TrimSpace(p), 32)
		if err != nil {
			return nil, fmt.Errorf("invalid float %q: %w", p, err)
		}
		vec[i] = float32(f)
	}
	return vec, nil
}
