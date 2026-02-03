package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ReyviRahman/to-backend/internal/models"
	"github.com/ReyviRahman/to-backend/internal/store"
	"github.com/go-chi/chi/v5"
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
	qq := store.PaginatedQuestionQuery{
		Limit:  20,
		Offset: 0,
	}

	qq, err := qq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(qq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	questions, meta, err := app.store.Questions.GetQuestions(ctx, qq)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	err = app.jsonResponse(w, http.StatusOK, "Berhasil Mendapatkan Data", questions, meta)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, errors.New("ID tidak valid"))
		return
	}

	var payload CreateQuestionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.validationErrorResponse(w, r, err)
		return
	}

	options := make(models.QuestionOptions, len(payload.Options))
	for i, opt := range payload.Options {
		options[i] = models.Option{
			Code:  opt.Code,
			Text:  opt.Text,
			Score: opt.Score,
		}
	}

	question := &models.Question{
		ID:           id,
		Category:     payload.Category,
		QuestionText: payload.QuestionText,
		Options:      options,
		Explanation:  payload.Explanation,
	}

	ctx := r.Context()
	if err := app.store.Questions.Update(ctx, question); err != nil {
		if err.Error() == "data tidak ditemukan" {
			app.notFoundResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "Question Update successfully", question); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) deleteQuestionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, errors.New("ID tidak valid"))
		return
	}

	ctx := r.Context()
	if err := app.store.Questions.Delete(ctx, id); err != nil {
		if err.Error() == "data tidak ditemukan" {
			app.notFoundResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "Question deleted successfully", nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
