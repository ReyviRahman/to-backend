package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	// Tambahkan ini agar pesan error menggunakan nama field di JSON (misal: "category"),
	// bukan nama struct Go (misal: "Category")
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func (app *application) parseValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, v := range validationErrors {
			// v.Tag() akan berisi "required", "oneof", "min", dll.
			switch v.Tag() {
			case "required":
				errors[v.Field()] = "wajib diisi"
			case "oneof":
				// v.Param() akan berisi nilai yang diperbolehkan (misal: "TIU TWK TKP")
				errors[v.Field()] = fmt.Sprintf("harus salah satu dari: %s", v.Param())
			case "min":
				errors[v.Field()] = fmt.Sprintf("minimal %s karakter", v.Param())
			case "max":
				errors[v.Field()] = fmt.Sprintf("maksimal %s karakter", v.Param())
			case "len":
				errors[v.Field()] = fmt.Sprintf("panjang harus %s karakter", v.Param())
			case "url":
				errors[v.Field()] = "format URL tidak valid"
			default:
				errors[v.Field()] = fmt.Sprintf("gagal pada aturan '%s'", v.Tag())
			}
		}
	}

	return errors
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		// Kamu bisa menambahkan pengecekan tipe error di sini jika ingin lebih spesifik
		return err
	}

	// Pastikan tidak ada data tambahan setelah JSON selesai di-decode
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}

// Tambahkan parameter 'headers ...any' di akhir
func (app *application) jsonResponse(w http.ResponseWriter, status int, message string, data any, meta ...any) error {

	// Kita buat struct baru yang punya field Meta dengan tag "omitempty"
	// omitempty artinya: kalau nil, field ini gak bakal muncul di JSON
	env := struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
		Meta    any    `json:"meta,omitempty"`
	}{
		Message: message,
		Data:    data,
	}

	// Cek apakah ada parameter meta yang dikirim?
	if len(meta) > 0 {
		env.Meta = meta[0]
	}

	return writeJSON(w, status, env)
}
