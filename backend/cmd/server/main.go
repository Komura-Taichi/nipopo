package main

import (
	"log"
	"net/http"

	"github.com/Komura-Taichi/nipopo/backend/internal/handler"
	"github.com/Komura-Taichi/nipopo/backend/internal/handler/middleware"
	"github.com/Komura-Taichi/nipopo/backend/internal/repository/inmem"
	"github.com/Komura-Taichi/nipopo/backend/internal/usecase"
)

func main() {
	// DI
	tagRepo := inmem.NewTagRepository()
	tagUsecase := usecase.NewTagUsecase(tagRepo)

	// DIが必要なハンドラ
	listTags := middleware.AuthStub("u_1")(handler.ListTags(tagUsecase))
	createTag := middleware.AuthStub("u_1")(handler.CreateTag(tagUsecase))

	// ルーティング
	mux := http.NewServeMux()

	// DIが不要なので、そのまま登録
	mux.HandleFunc("/healthz", handler.Healthz)

	// DIが必要なので、HTTPメソッドで分岐
	mux.HandleFunc("/v1/tags", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listTags.ServeHTTP(w, r)
		case http.MethodPost:
			createTag.ServeHTTP(w, r)
		default:
			http.Error(w, "undefined method", http.StatusMethodNotAllowed)
		}
	})

	address := ":8080"
	log.Println("listening on ", address)
	err := http.ListenAndServe(address, mux)
	if err != nil {
		log.Fatalf("server start error: %v", err)
	}
}
