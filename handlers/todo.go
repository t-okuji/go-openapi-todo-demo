package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/t-okuji/go-openapi-todo-demo/ent"
	"github.com/t-okuji/go-openapi-todo-demo/ent/category"
	"github.com/t-okuji/go-openapi-todo-demo/ent/todo"
	"github.com/t-okuji/go-openapi-todo-demo/types"
	"github.com/t-okuji/go-openapi-todo-demo/utils"
)

// GetTodosHandler は GET /todos リクエストを処理する
func GetTodosHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// データベースから全ての Todo を取得
		todos, err := client.Todo.Query().All(ctx)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "Database error occurred")
			log.Printf("Todo fetch error: %v", err)
			return
		}

		// Ent エンティティをレスポンス形式に変換
		todoResponses := make([]types.TodoResponse, len(todos))
		for i, todo := range todos {
			todoResponses[i] = utils.ConvertToTodoResponse(todo)
		}

		utils.SendJSONResponse(w, http.StatusOK, todoResponses)
	}
}

// CreateTodoHandler は POST /todos リクエストを処理する
func CreateTodoHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// リクエストボディをパース
		var input types.TodoInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format")
			return
		}

		// タイトルの検証
		if input.Title == "" {
			utils.SendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Title is required")
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
			categoryUUID, ok := utils.ParseUUID(w, *input.CategoryID)
			if !ok {
				return
			}

			// カテゴリの存在確認
			exists, err := client.Category.Query().
				Where(category.ID(categoryUUID)).
				Exist(ctx)
			if err != nil {
				utils.SendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "Database error occurred")
				log.Printf("Category existence check error: %v", err)
				return
			}
			if !exists {
				utils.SendErrorResponse(w, http.StatusBadRequest, "CATEGORY_NOT_FOUND", "Specified category not found")
				return
			}

			createQuery.SetCategoryID(categoryUUID)
		}

		// Todoを作成
		todo, err := createQuery.Save(ctx)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "Failed to create Todo")
			log.Printf("Todo creation error: %v", err)
			return
		}

		// レスポンスを返却
		response := utils.ConvertToTodoResponse(todo)
		utils.SendJSONResponse(w, http.StatusCreated, response)
	}
}

// GetTodoByIDHandler は GET /todos/{todoId} リクエストを処理する
func GetTodoByIDHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// URLパラメータからtodoIdを取得
		todoID := chi.URLParam(r, "todoId")
		todoUUID, ok := utils.ParseUUID(w, todoID)
		if !ok {
			return
		}

		// Todoを取得
		todo, err := client.Todo.Get(ctx, todoUUID)
		if err != nil {
			if ent.IsNotFound(err) {
				utils.SendErrorResponse(w, http.StatusNotFound, "TODO_NOT_FOUND", "Specified Todo not found")
				return
			}
			utils.SendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "Database error occurred")
			log.Printf("Todo fetch error: %v", err)
			return
		}

		// レスポンスを返却
		response := utils.ConvertToTodoResponse(todo)
		utils.SendJSONResponse(w, http.StatusOK, response)
	}
}

// UpdateTodoHandler は PUT /todos/{todoId} リクエストを処理する
func UpdateTodoHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// URLパラメータからtodoIdを取得
		todoID := chi.URLParam(r, "todoId")
		todoUUID, ok := utils.ParseUUID(w, todoID)
		if !ok {
			return
		}

		// リクエストボディをパース
		var input types.TodoInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format")
			return
		}

		// タイトルの検証
		if input.Title == "" {
			utils.SendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Title is required")
			return
		}

		// 対象Todoの存在確認
		exists, err := client.Todo.Query().Where(todo.ID(todoUUID)).Exist(ctx)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "Database error occurred")
			log.Printf("Todo existence check error: %v", err)
			return
		}
		if !exists {
			utils.SendErrorResponse(w, http.StatusNotFound, "TODO_NOT_FOUND", "指定されたTodoが見つかりません")
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
				categoryUUID, ok := utils.ParseUUID(w, *input.CategoryID)
				if !ok {
					return
				}

				// カテゴリの存在確認
				exists, err := client.Category.Query().
					Where(category.ID(categoryUUID)).
					Exist(ctx)
				if err != nil {
					utils.SendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "Database error occurred")
					log.Printf("Category existence check error: %v", err)
					return
				}
				if !exists {
					utils.SendErrorResponse(w, http.StatusBadRequest, "CATEGORY_NOT_FOUND", "Specified category not found")
					return
				}

				updateQuery.SetCategoryID(categoryUUID)
			}
		}

		// Todoを更新
		todo, err := updateQuery.Save(ctx)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "Failed to update Todo")
			log.Printf("Todo update error: %v", err)
			return
		}

		// レスポンスを返却
		response := utils.ConvertToTodoResponse(todo)
		utils.SendJSONResponse(w, http.StatusOK, response)
	}
}

// DeleteTodoHandler は DELETE /todos/{todoId} リクエストを処理する
func DeleteTodoHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// URLパラメータからtodoIdを取得
		todoID := chi.URLParam(r, "todoId")
		todoUUID, ok := utils.ParseUUID(w, todoID)
		if !ok {
			return
		}

		// Todoを削除
		err := client.Todo.DeleteOneID(todoUUID).Exec(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				utils.SendErrorResponse(w, http.StatusNotFound, "TODO_NOT_FOUND", "Specified Todo not found")
				return
			}
			utils.SendErrorResponse(w, http.StatusInternalServerError, "DB_ERROR", "Database error occurred")
			log.Printf("Todo deletion error: %v", err)
			return
		}

		// 204 No Contentを返却
		w.WriteHeader(http.StatusNoContent)
	}
}
