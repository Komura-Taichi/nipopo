package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Komura-Taichi/nipopo/backend/internal/usecase"
)

type createRecordRequest struct {
	Date   string   `json:"date"`
	Effort int      `json:"effort"`
	TagIDs []string `json:"tag_ids"`
	Body   string   `json:"body"`
}

func CreateRecord(creator usecase.RecordCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createRecordRequest
		// json形式が正しくない場合のテストに対応
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			writeErrorJSON(w, http.StatusBadRequest, "invalid json",
				map[string]any{"reason": err.Error()},
			)

			return
		}

		// 日付の形式が正しくない (YYYY-mm-ddじゃない) 場合のテストに対応
		dateStr := strings.TrimSpace(req.Date)
		if dateStr == "" {
			writeErrorJSON(w, http.StatusBadRequest, "empty date",
				map[string]any{"field": "date"},
			)

			return
		}
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			writeErrorJSON(w, http.StatusBadRequest, "invalid date",
				map[string]any{"field": "date", "reason": err.Error()},
			)

			return
		}

		// 頑張り度が範囲外 (0...5の外) 場合のテストに対応
		if req.Effort < 0 || req.Effort > 5 {
			writeErrorJSON(w, http.StatusBadRequest, "invalid effort",
				map[string]any{"field": "effort", "value": req.Effort},
			)

			return
		}

		// 本文が抜けてる場合や空白の場合のテストに対応
		body := strings.TrimSpace(req.Body)
		if body == "" {
			writeErrorJSON(w, http.StatusBadRequest, "body is required",
				map[string]any{"field": "body"},
			)

			return
		}

		// タグID一覧が抜けてる場合やタグ一覧が渡されなかった場合のテストに対応
		if len(req.TagIDs) == 0 {
			writeErrorJSON(w, http.StatusBadRequest, "tag_ids is required",
				map[string]any{"field": "tag_ids"},
			)

			return
		}
		// タグID一覧が空白を含む場合や重複を含む場合のテストに対応
		tagIDs := make([]string, 0, len(req.TagIDs))
		seen := map[string]struct{}{}
		for _, tid := range req.TagIDs {
			tidNoSpace := strings.TrimSpace(tid)
			if tidNoSpace == "" {
				writeErrorJSON(w, http.StatusBadRequest, "tag_ids contains empty",
					map[string]any{"field": "tag_ids"},
				)

				return
			}

			if _, isDuplicated := seen[tidNoSpace]; isDuplicated {
				writeErrorJSON(w, http.StatusBadRequest, "tag_ids contains duplication",
					map[string]any{"field": "tag_ids", "value": tidNoSpace},
				)

				return
			}
			seen[tidNoSpace] = struct{}{}
			tagIDs = append(tagIDs, tidNoSpace)
		}

		// ユーザIDの取得 (r.Context() に埋め込まれている前提)
		userID, ok := requireUserID(w, r)
		if !ok {
			return
		}

		// usecase層の呼び出し
		createdRecord, err := creator.Create(r.Context(), userID, usecase.CreateRecordInput{
			Date:   date,
			Effort: req.Effort,
			TagIDs: tagIDs,
			Body:   body,
		})

		if err != nil {
			// 指定した日の記録が既に存在している場合のテストに対応
			var alreadyExistsErr *usecase.RecordAlreadyExistsError
			if errors.As(err, &alreadyExistsErr) {
				writeErrorJSON(w, http.StatusConflict, "record of that day already exists",
					map[string]any{"existing_id": alreadyExistsErr.ExistingID},
				)

				return
			}

			// その他のエラーの場合は500
			writeErrorJSON(w, http.StatusInternalServerError, "internal server error",
				map[string]any{"reason": err.Error()},
			)

			return
		}

		// レスポンスの形式に変換
		response := Record{
			ID:        createdRecord.ID,
			Date:      createdRecord.Date.Format("2006-01-02"),
			CreatedAt: createdRecord.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: createdRecord.UpdatedAt.UTC().Format(time.RFC3339),
			Effort:    createdRecord.Effort,
			// TODO: IDだけ返すようにしてるが、実際にはrepoにアクセスして、名前もとる必要あり。
			Tags: toSimpleTags(createdRecord.TagIDs),
			Body: createdRecord.Body,
		}

		writeJSON(w, http.StatusCreated, response)
	}
}

func toSimpleTags(tagIDs []string) []SimpleTag {
	tags := make([]SimpleTag, 0, len(tagIDs))

	for _, tid := range tagIDs {
		tags = append(tags, SimpleTag{ID: tid})
	}

	return tags
}
