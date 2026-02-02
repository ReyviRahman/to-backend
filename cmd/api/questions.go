package main

import (
	"net/http"

	"github.com/ReyviRahman/to-backend/internal/models"
)

type CreateQuestionPayload struct {
	Category     string                 `json:"category"`
	QuestionText string                 `json:"question_text"`
	Options      models.QuestionOptions `json:"options"`
	Explanation  string                 `json:"explanation"`
}

func (app *application) createQuestionHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateQuestionPayload

	// 1. Baca JSON (Gunakan helper readJSON yang biasa kamu pakai)
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// 2. Validasi (Gunakan validator v10)
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// 3. Mapping: Payload -> Model
	// Di sini kita pindahkan data dari struct "input" ke struct "database"
	question := &models.Question{
		Category:     payload.Category,
		QuestionText: payload.QuestionText,
		Options:      payload.Options,
		Explanation:  payload.Explanation,
	}

	// 4. Simpan ke Database via Store
	ctx := r.Context()
	if err := app.store.Questions.Create(ctx, question); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// 5. Kirim Response (balikan object 'question' yang sudah ada ID-nya)
	if err := app.jsonResponse(w, http.StatusCreated, question); err != nil {
		app.internalServerError(w, r, err)
	}
}
