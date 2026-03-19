package inmem

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
	"github.com/Komura-Taichi/nipopo/backend/internal/repository"
)

var _ repository.TagRepository = (*TagRepository)(nil) // インタフェースを満たすかコンパイル時にチェック

type TagRepository struct {
	mu sync.RWMutex

	// 作成順を保持
	tags []entity.Tag

	// ユーザID+タグID -> index のマッピング
	uidAndTIDToIdx map[string]map[string]int
	// ユーザID+タグ名 -> index のマッピング
	uidAndNameToIdx map[string]map[string]int
}

func NewTagRepository() *TagRepository {
	return &TagRepository{
		mu:              sync.RWMutex{},
		tags:            make([]entity.Tag, 0),
		uidAndTIDToIdx:  make(map[string]map[string]int),
		uidAndNameToIdx: make(map[string]map[string]int),
	}
}

func (r *TagRepository) FindByIDs(ctx context.Context, userID string, ids []string) ([]entity.Tag, error) {
	_ = ctx

	if len(ids) == 0 {
		return []entity.Tag{}, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// そもそもタグ一つもが無いならエラー
	idToIdx, ok := r.uidAndTIDToIdx[userID]
	if !ok {
		return nil, repository.ErrTagNotFound
	}

	tags := make([]entity.Tag, 0, len(ids))
	for _, id := range ids {
		idx, ok := idToIdx[id]
		// 存在してるタグか？
		if !ok {
			return nil, repository.ErrTagNotFound
		}
		tags = append(tags, r.tags[idx])
	}

	return tags, nil
}

func (r *TagRepository) FindByName(ctx context.Context, userID, name string) (entity.Tag, bool, error) {
	_ = ctx

	key := strings.TrimSpace(name)

	r.mu.RLock()
	defer r.mu.RUnlock()

	nameToIdx, ok := r.uidAndNameToIdx[userID]
	if !ok {
		return entity.Tag{}, false, nil
	}
	idx, ok := nameToIdx[key]
	if !ok {
		return entity.Tag{}, false, nil
	}

	return r.tags[idx], true, nil
}

func (r *TagRepository) Create(ctx context.Context, userID, name string) (entity.Tag, error) {
	_ = ctx

	key := strings.TrimSpace(name)

	r.mu.Lock()
	defer r.mu.Unlock()

	nameToIdx, ok := r.uidAndNameToIdx[userID]
	// そのユーザが初めてタグを作る場合、初期化
	if !ok {
		nameToIdx = map[string]int{}
		r.uidAndNameToIdx[userID] = nameToIdx
	}

	idToIdx, ok := r.uidAndTIDToIdx[userID]
	if !ok {
		idToIdx = map[string]int{}
		r.uidAndTIDToIdx[userID] = idToIdx
	}

	// 多重で作成されないように
	if _, exists := nameToIdx[key]; exists {
		return entity.Tag{}, repository.ErrAlreadyTagExists
	}

	// インメモリでは仮の連番ID
	id := fmt.Sprintf("t_%d", len(r.tags)+1)

	tag := entity.Tag{UserID: userID, ID: id, Name: key}
	r.tags = append(r.tags, tag)
	idx := len(r.tags) - 1

	nameToIdx[key] = idx // nameToIdxとuidAndNameToIdxは結びついているので
	idToIdx[id] = idx    // idToIdxとuidAndTIDToIdxは結びついてる

	return tag, nil
}

func (r *TagRepository) List(ctx context.Context, userID, q string, limit int, cursor string) (entity.TagsPage, error) {
	_ = ctx

	q = strings.TrimSpace(q)
	cursor = strings.TrimSpace(cursor)

	r.mu.RLock()
	defer r.mu.RUnlock()

	// userID, qによるフィルタ
	filtered := make([]entity.Tag, 0, len(r.tags))
	for _, tag := range r.tags {
		if tag.UserID != userID {
			continue
		}
		if q == "" || strings.Contains(tag.Name, q) {
			filtered = append(filtered, tag)
		}
	}

	// cursorをもとに開始位置を決める (空文字列なら最初から。cursorのIDの次から)
	start := 0
	if cursor != "" {
		found := false
		for i, tag := range filtered {
			if tag.ID == cursor {
				start = i + 1
				found = true
				break
			}
		}
		// cursorに対応するタグIDが存在しない場合はエラー
		if !found {
			return entity.TagsPage{}, repository.ErrCursorNotFound
		}
	}

	// limit分切り出す
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	items := filtered[start:end]

	// next_cursorは次がなければ空文字列
	nextCursor := ""
	if len(items) > 0 && end < len(filtered) {
		nextCursor = items[len(items)-1].ID
	}

	return entity.TagsPage{
		Items:      items,
		NextCursor: nextCursor,
	}, nil
}
