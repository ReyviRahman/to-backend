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

type MetaData struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

func (s *QuestionStore) GetQuestions(ctx context.Context, qq PaginatedQuestionQuery) ([]models.Question, MetaData, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// 1. Query Pertama: Hitung Total Data (Tanpa Limit/Offset)
	var totalItems int
	countQuery := `SELECT COUNT(id) FROM questions`

	// Jika nanti ada filter (misal WHERE category = 'TIU'),
	// pastikan countQuery juga pakai WHERE yang sama.
	err := s.db.QueryRowContext(ctx, countQuery).Scan(&totalItems)
	if err != nil {
		return nil, MetaData{}, err
	}

	// 2. Query Kedua: Ambil Data Sebenarnya (Pakai Limit/Offset)
	query := `
        SELECT id, category, question_text, options, explanation, created_at, updated_at
        FROM questions
				WHERE ($1 = '' OR question_text ILIKE '%' || $1 || '%')
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := s.db.QueryContext(ctx, query, qq.Search, qq.Limit, qq.Offset)
	if err != nil {
		return nil, MetaData{}, err
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
			return nil, MetaData{}, err
		}
		questions = append(questions, q)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	// 3. Hitung Kalkulasi Metadata
	totalPages := 0
	if qq.Limit > 0 {
		// Rumus total page: ceil(totalItems / limit)
		// Cara integer di Go: (total + limit - 1) / limit
		totalPages = (totalItems + qq.Limit - 1) / qq.Limit
	}

	currentPage := 1
	if qq.Limit > 0 {
		currentPage = (qq.Offset / qq.Limit) + 1
	}

	meta := MetaData{
		CurrentPage: currentPage,
		Limit:       qq.Limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}

	return questions, meta, nil
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
