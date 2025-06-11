package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/t-okuji/go-openapi-todo-demo/ent"
	"github.com/t-okuji/go-openapi-todo-demo/handlers"
)

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
		log.Printf("Failed to load .env file: %v", err)
	}

	// データベース接続情報を環境変数から取得
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	log.Printf("Connecting to database: %s", dsn)
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

	// OpenAPI ドキュメント用エンドポイント
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi-ui.html")
	})

	// OpenAPI仕様ファイルを提供するエンドポイント
	r.Handle("/openapi/*", http.StripPrefix("/openapi/", http.FileServer(http.Dir("openapi"))))

	// Todo API エンドポイント
	r.Get("/todos", handlers.GetTodosHandler(client))
	r.Post("/todos", handlers.CreateTodoHandler(client))
	r.Get("/todos/{todoId}", handlers.GetTodoByIDHandler(client))
	r.Put("/todos/{todoId}", handlers.UpdateTodoHandler(client))
	r.Delete("/todos/{todoId}", handlers.DeleteTodoHandler(client))

	log.Printf("Starting server: http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
