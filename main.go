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
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/t-okuji/go-openapi-todo-demo/ent"
	"github.com/t-okuji/go-openapi-todo-demo/ent/category"
	"github.com/t-okuji/go-openapi-todo-demo/ent/todo"
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

// TodoInput は API リクエスト用の Todo 入力データを表す
type TodoInput struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
	CategoryID  *string `json:"categoryId,omitempty"`
}

// ErrorResponse は API エラーレスポンスを表す
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
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

// sendJSONResponse は JSON レスポンスを送信する共通関数
func sendJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON エンコードエラー: %v", err)
	}
}

// sendErrorResponse はエラーレスポンスを送信する共通関数
func sendErrorResponse(w http.ResponseWriter, status int, code string, message string) {
	var errResp ErrorResponse
	errResp.Error.Code = code
	errResp.Error.Message = message
	sendJSONResponse(w, status, errResp)
}

// parseUUID は文字列をUUIDとしてパースし、エラーがあればエラーレスポンスを送信する
func parseUUID(w http.ResponseWriter, uuidStr string) (uuid.UUID, bool) {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "INVALID_UUID", "無効なUUID形式です")
		return uuid.UUID{}, false
	}
	return id, true
}

// getTodosHandler は GET /todos リクエストを処理する
func getTodosHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// データベースから全ての Todo を取得
		todos, err := client.Todo.Query().All(ctx)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "データベースエラーが発生しました")
			log.Printf("Todo 取得エラー: %v", err)
			return
		}

		// Ent エンティティをレスポンス形式に変換
		todoResponses := make([]TodoResponse, len(todos))
		for i, todo := range todos {
			todoResponses[i] = convertToTodoResponse(todo)
		}

		sendJSONResponse(w, http.StatusOK, todoResponses)
	}
}

// createTodoHandler は POST /todos リクエストを処理する
func createTodoHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// リクエストボディをパース
		var input TodoInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, "INVALID_JSON", "不正なJSONフォーマットです")
			return
		}

		// タイトルの検証
		if input.Title == "" {
			sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "タイトルは必須です")
			return
		}

		// Todo作成クエリを構築
		createQuery := client.Todo.Create().
			SetTitle(input.Title).
			SetCompleted(false) // デフォルトはfalse

		// オプショナルフィールドの処理
		if input.Description != nil {
			createQuery.SetDescription(*input.Description)
		}

		if input.Completed != nil {
			createQuery.SetCompleted(*input.Completed)
		}

		// カテゴリIDの処理
		if input.CategoryID != nil {
			categoryUUID, ok := parseUUID(w, *input.CategoryID)
			if !ok {
				return
			}

			// カテゴリの存在確認
			exists, err := client.Category.Query().
				Where(category.ID(categoryUUID)).
				Exist(ctx)
			if err != nil {
				sendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "データベースエラーが発生しました")
				log.Printf("カテゴリ存在確認エラー: %v", err)
				return
			}
			if !exists {
				sendErrorResponse(w, http.StatusBadRequest, "CATEGORY_NOT_FOUND", "指定されたカテゴリが存在しません")
				return
			}

			createQuery.SetCategoryID(categoryUUID)
		}

		// Todoを作成
		todo, err := createQuery.Save(ctx)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "Todoの作成に失敗しました")
			log.Printf("Todo作成エラー: %v", err)
			return
		}

		// レスポンスを返却
		response := convertToTodoResponse(todo)
		sendJSONResponse(w, http.StatusCreated, response)
	}
}

// getTodoByIDHandler は GET /todos/{todoId} リクエストを処理する
func getTodoByIDHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// URLパラメータからtodoIdを取得
		todoID := chi.URLParam(r, "todoId")
		todoUUID, ok := parseUUID(w, todoID)
		if !ok {
			return
		}

		// Todoを取得
		todo, err := client.Todo.Get(ctx, todoUUID)
		if err != nil {
			if ent.IsNotFound(err) {
				sendErrorResponse(w, http.StatusNotFound, "TODO_NOT_FOUND", "指定されたTodoが見つかりません")
				return
			}
			sendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "データベースエラーが発生しました")
			log.Printf("Todo取得エラー: %v", err)
			return
		}

		// レスポンスを返却
		response := convertToTodoResponse(todo)
		sendJSONResponse(w, http.StatusOK, response)
	}
}

// updateTodoHandler は PUT /todos/{todoId} リクエストを処理する
func updateTodoHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// URLパラメータからtodoIdを取得
		todoID := chi.URLParam(r, "todoId")
		todoUUID, ok := parseUUID(w, todoID)
		if !ok {
			return
		}

		// リクエストボディをパース
		var input TodoInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, "INVALID_JSON", "不正なJSONフォーマットです")
			return
		}

		// タイトルの検証
		if input.Title == "" {
			sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "タイトルは必須です")
			return
		}

		// 対象Todoの存在確認
		exists, err := client.Todo.Query().Where(todo.ID(todoUUID)).Exist(ctx)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "データベースエラーが発生しました")
			log.Printf("Todo存在確認エラー: %v", err)
			return
		}
		if !exists {
			sendErrorResponse(w, http.StatusNotFound, "TODO_NOT_FOUND", "指定されたTodoが見つかりません")
			return
		}

		// Todo更新クエリを構築
		updateQuery := client.Todo.UpdateOneID(todoUUID).
			SetTitle(input.Title)

		// オプショナルフィールドの処理
		if input.Description != nil {
			if *input.Description == "" {
				updateQuery.ClearDescription()
			} else {
				updateQuery.SetDescription(*input.Description)
			}
		}

		if input.Completed != nil {
			updateQuery.SetCompleted(*input.Completed)
		}

		// カテゴリIDの処理
		if input.CategoryID != nil {
			if *input.CategoryID == "" {
				updateQuery.ClearCategoryID()
			} else {
				categoryUUID, ok := parseUUID(w, *input.CategoryID)
				if !ok {
					return
				}

				// カテゴリの存在確認
				exists, err := client.Category.Query().
					Where(category.ID(categoryUUID)).
					Exist(ctx)
				if err != nil {
					sendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "データベースエラーが発生しました")
					log.Printf("カテゴリ存在確認エラー: %v", err)
					return
				}
				if !exists {
					sendErrorResponse(w, http.StatusBadRequest, "CATEGORY_NOT_FOUND", "指定されたカテゴリが存在しません")
					return
				}

				updateQuery.SetCategoryID(categoryUUID)
			}
		}

		// Todoを更新
		todo, err := updateQuery.Save(ctx)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "Todoの更新に失敗しました")
			log.Printf("Todo更新エラー: %v", err)
			return
		}

		// レスポンスを返却
		response := convertToTodoResponse(todo)
		sendJSONResponse(w, http.StatusOK, response)
	}
}

// deleteTodoHandler は DELETE /todos/{todoId} リクエストを処理する
func deleteTodoHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// URLパラメータからtodoIdを取得
		todoID := chi.URLParam(r, "todoId")
		todoUUID, ok := parseUUID(w, todoID)
		if !ok {
			return
		}

		// Todoを削除
		err := client.Todo.DeleteOneID(todoUUID).Exec(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				sendErrorResponse(w, http.StatusNotFound, "TODO_NOT_FOUND", "指定されたTodoが見つかりません")
				return
			}
			sendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "データベースエラーが発生しました")
			log.Printf("Todo削除エラー: %v", err)
			return
		}

		// 204 No Contentを返却
		w.WriteHeader(http.StatusNoContent)
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

	// CORS設定
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// Todo API エンドポイント
	r.Get("/todos", getTodosHandler(client))
	r.Post("/todos", createTodoHandler(client))
	r.Get("/todos/{todoId}", getTodoByIDHandler(client))
	r.Put("/todos/{todoId}", updateTodoHandler(client))
	r.Delete("/todos/{todoId}", deleteTodoHandler(client))

	log.Printf("サーバーを起動します: http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
