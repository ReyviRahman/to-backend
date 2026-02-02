package main

import (
	"net/http"

	"github.com/ReyviRahman/to-backend/internal/models"
)

// internal/handler/questions.go (atau di mana kamu mendefinisikan payload)

type CreateQuestionPayload struct {
	// Wajib diisi, dan harus salah satu dari TIU, TWK, atau TKP
	Category string `json:"category" validate:"required,oneof=TIU TWK TKP"`

	// Wajib diisi, minimal 10 karakter, maksimal 500 karakter
	QuestionText string `json:"question_text" validate:"required,min=10,max=500"`

	// Field ini opsional (boleh kosong), tapi jika diisi harus berupa URL yang valid
	QuestionImageURL string `json:"question_image_url" validate:"omitempty,url"`

	// Wajib ada, minimal harus kirim 2 opsi, dan 'dive' berarti validasi setiap item di dalamnya
	Options []OptionPayload `json:"options" validate:"required,min=2,dive"`

	// Wajib diisi
	Explanation string `json:"explanation" validate:"required"`
}

type OptionPayload struct {
	// Wajib diisi, panjang karakter harus tepat 1 (misal: A, B, C)
	Code string `json:"code" validate:"required,len=1"`
	Text string `json:"text" validate:"required"`
	// Skor minimal 0, maksimal 5
	Score int `json:"score" validate:"min=0,max=5"`
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
		app.validationErrorResponse(w, r, err)
		return
	}

	// 3. Mapping: Payload -> Model
	// Di sini kita pindahkan data dari struct "input" ke struct "database"

	options := make(models.QuestionOptions, len(payload.Options))
	for i, opt := range payload.Options {
		options[i] = models.Option{
			Code:  opt.Code,
			Text:  opt.Text,
			Score: opt.Score,
		}
	}
	question := &models.Question{
		Category:     payload.Category,
		QuestionText: payload.QuestionText,
		Options:      options,
		Explanation:  payload.Explanation,
	}

	// 4. Simpan ke Database via Store
	ctx := r.Context()
	if err := app.store.Questions.Create(ctx, question); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// 5. Kirim Response (balikan object 'question' yang sudah ada ID-nya)
	if err := app.jsonResponse(w, http.StatusCreated, "Question created successfully", question); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getQuestionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	questions, err := app.store.Questions.GetQuestions(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	err = app.jsonResponse(w, http.StatusOK, "Berhasil Mendapatkan Data", questions)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
