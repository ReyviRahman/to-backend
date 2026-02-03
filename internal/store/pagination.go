package store

import (
	"net/http"
	"strconv"
)

type PaginatedQuestionQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Search string `json:"search" validate:"max=100"`
}

func (qq PaginatedQuestionQuery) Parse(r *http.Request) (PaginatedQuestionQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return qq, nil
		}

		qq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return qq, nil
		}

		qq.Offset = l
	}

	search := qs.Get("search")
	if search != "" {
		qq.Search = search
	}

	return qq, nil
}
