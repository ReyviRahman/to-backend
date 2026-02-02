package store

import (
	"context"
	"database/sql"

	"github.com/ReyviRahman/to-backend/internal/models"
)

type Storage struct {
	Questions interface {
		Create(ctx context.Context, question *models.Question) error
		GetQuestions(ctx context.Context) ([]models.Question, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Questions: &QuestionStore{db},
	}
}
