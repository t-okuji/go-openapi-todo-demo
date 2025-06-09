package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/t-okuji/go-openapi-todo-demo/ent"
)

// TodoResponse は API レスポンス用の Todo エンティティを表す
type TodoResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Completed   bool       `json:"completed"`
	CategoryID  *string    `json:"categoryId,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// convertToTodoResponse は Ent の Todo エンティティを TodoResponse に変換する
func convertToTodoResponse(todo *ent.Todo) TodoResponse {
	response := TodoResponse{
		ID:        todo.ID.String(),
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt,
	}

	// オプションフィールドの処理
	if todo.Description != "" {
		response.Description = &todo.Description
	}

	if todo.CategoryID != nil {
		categoryID := todo.CategoryID.String()
		response.CategoryID = &categoryID
	}

	if !todo.UpdatedAt.IsZero() {
		response.UpdatedAt = &todo.UpdatedAt
	}

	return response
}

// getTodosHandler は GET /todos リクエストを処理する
func getTodosHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// データベースから全ての Todo を取得
		todos, err := client.Todo.Query().All(ctx)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Todo 取得エラー: %v", err)
			return
		}

		// Ent エンティティをレスポンス形式に変換
		todoResponses := make([]TodoResponse, len(todos))
		for i, todo := range todos {
			todoResponses[i] = convertToTodoResponse(todo)
		}

		// コンテンツタイプを設定し JSON レスポンスをエンコード
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(todoResponses); err != nil {
			log.Printf("JSON エンコードエラー: %v", err)
		}
	}
}

// Open は新しいデータベース接続を開く
func Open(databaseUrl string) *ent.Client {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	// `db` から ent.Driver を作成
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

func main() {
	// .envファイルを読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf(".envファイルの読み込みに失敗しました: %v", err)
	}

	// データベース接続情報を環境変数から取得
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	
	log.Printf("データベースに接続します: %s", dsn)
	client := Open(dsn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// Todo API エンドポイント
	r.Get("/todos", getTodosHandler(client))

	log.Printf("サーバーを起動します: http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
