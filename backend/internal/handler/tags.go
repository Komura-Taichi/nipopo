package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Komura-Taichi/nipopo/backend/internal/usecase"
)

type createTagRequest struct {
	Name string `json:"name"`
}

func ListTags(lister usecase.TagsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		cursor := r.URL.Query().Get("cursor")

		limitStr := r.URL.Query().Get("limit")
		limit, err := parseLimit(limitStr)
		if err != nil {
			writeErrorJSON(w, http.StatusBadRequest, "invalid limit",
				map[string]any{
					"field":  "limit",
					"value":  limitStr,
					"reason": err.Error(),
				})

			return
		}

		// ユーザIDの取得 (r.Context() に埋め込まれている前提)
		userID, ok := requireUserID(w, r)
		if !ok {
			return
		}

		page, err := lister.List(r.Context(), userID, q, limit, cursor)
		if err != nil {
			writeErrorJSON(w, http.StatusInternalServerError, "internal server error",
				map[string]any{"reason": err.Error()},
			)

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

func CreateTag(creator usecase.TagCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTagRequest
		// json形式が正しくない場合のテストに対応
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			writeErrorJSON(w, http.StatusBadRequest, "invalid json",
				map[string]any{"reason": err.Error()},
			)

			return
		}

		// タグ名が空だった場合のテストに対応
		name := strings.TrimSpace(req.Name)
		if name == "" {
			writeErrorJSON(w, http.StatusBadRequest, "empty tag name",
				map[string]any{"field": "name"},
			)

			return
		}

		// ユーザIDの取得 (r.Context() に埋め込まれている前提)
		userID, ok := requireUserID(w, r)
		if !ok {
			return
		}

		createdTag, err := creator.Create(r.Context(), userID, name)
		if err != nil {
			writeErrorJSON(w, http.StatusInternalServerError, "internal server error",
				map[string]any{"reason": err.Error()},
			)

			return
		}

		status := http.StatusOK
		if createdTag.Created {
			status = http.StatusCreated
		}

		response := Tag{ID: createdTag.Tag.ID, Name: createdTag.Tag.Name}
		writeJSON(w, status, response)
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
