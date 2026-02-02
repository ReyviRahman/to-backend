package models

import (
	"time"
)

type Question struct {
	ID                  int64           `json:"id"`
	Category            string          `json:"category"`
	QuestionText        string          `json:"question_text"`
	QuestionImageURL    *string         `json:"question_image_url"`
	Options             QuestionOptions `json:"options"`
	Explanation         string          `json:"explanation"`
	ExplanationImageURL *string         `json:"explanation_image_url"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}
