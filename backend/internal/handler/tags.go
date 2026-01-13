package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Komura-Taichi/nipopo/backend/internal/usecase"
)

func ListTags(lister usecase.TagsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		cursor := r.URL.Query().Get("cursor")

		limit, err := parseLimit(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		page, err := lister.List(r.Context(), q, limit, cursor)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		response := TagsPage{
			Items:      make([]Tag, 0, len(page.Items)),
			NextCursor: page.NextCursor,
		}
		for _, tag := range page.Items {
			response.Items = append(response.Items, Tag{ID: tag.ID, Name: tag.Name})
		}

		writeJSON(w, http.StatusOK, response)
	}
}

func parseLimit(limS string) (int, error) {
	limSNoSpace := strings.TrimSpace(limS)
	if limSNoSpace == "" {
		return 5, nil
	}

	limit, err := strconv.Atoi(limSNoSpace)
	// 整数じゃない場合のテストに対応
	if err != nil {
		return 0, err
	}

	// 正の数または0じゃない場合のテストに対応
	if limit < 0 {
		return 0, errors.New("limit must be non-negative")
	}

	// limitが多すぎるとフェッチが重いため、最大値を設定
	if limit > 30 {
		limit = 30
	}

	return limit, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		log.Printf("error happened while encoding JSON: %v", err)
	}
}
