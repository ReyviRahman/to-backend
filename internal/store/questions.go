package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/ReyviRahman/to-backend/internal/models"
)

type QuestionStore struct {
	db *sql.DB
}

func (s *QuestionStore) GetQuestions(ctx context.Context) ([]models.Question, error) {
	query := `
		SELECT id, category, question_text, options, explanation, created_at, updated_at
		FROM questions
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		err := rows.Scan(
			&q.ID,
			&q.Category,
			&q.QuestionText,
			&q.Options,
			&q.Explanation,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
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
