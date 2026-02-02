package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		question.Category,
		question.QuestionText,
		question.Options,
		question.Explanation,
	).Scan(&question.ID, &question.CreatedAt, &question.UpdatedAt)

	return err
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

func (s *QuestionStore) Update(ctx context.Context, question *models.Question) error {
	query := `
		UPDATE questions
		SET category = $1, question_text = $2, options = $3, explanation = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		question.Category,
		question.QuestionText,
		question.Options,
		question.Explanation,
		question.ID,
	).Scan(&question.UpdatedAt)

	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return errors.New("data tidak ditemukan")
		default:
			return err
		}
	}

	return nil
}

func (s *QuestionStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM questions WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("data tidak ditemukan")
	}

	return nil
}
