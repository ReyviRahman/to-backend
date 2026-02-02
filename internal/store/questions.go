package store

import (
	"context"
	"database/sql"

	"github.com/ReyviRahman/to-backend/internal/models"
)

type QuestionStore struct {
	db *sql.DB
}

func (s *QuestionStore) Create(ctx context.Context, question *models.Question) error {
	query := `
		INSERT INTO questions (category, question_text, options, explanation)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		question.Category,
		question.QuestionText,
		question.Options,
		question.Explanation,
	).Scan(&question.ID, &question.CreatedAt, &question.UpdatedAt)

	return err
}
