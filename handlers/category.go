package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/t-okuji/go-openapi-todo-demo/ent"
	"github.com/t-okuji/go-openapi-todo-demo/ent/category"
	"github.com/t-okuji/go-openapi-todo-demo/types"
	"github.com/t-okuji/go-openapi-todo-demo/utils"
)

// GetCategories は全カテゴリの一覧を取得するハンドラー
func GetCategories(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// データベースから全カテゴリを取得
		categories, err := client.Category.
			Query().
			Order(ent.Asc(category.FieldCreatedAt)).
			All(context.Background())
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch categories")
			return
		}

		// レスポンス用に変換
		responses := make([]types.CategoryResponse, len(categories))
		for i, cat := range categories {
			responses[i] = utils.ConvertToCategoryResponse(cat)
		}

		utils.SendJSONResponse(w, http.StatusOK, responses)
	}
}

// CreateCategory は新規カテゴリを作成するハンドラー
func CreateCategory(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストボディをパース
		var input types.CategoryInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}

		// バリデーション
		if input.Name == "" {
			utils.SendErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
			return
		}
		if len(input.Name) > 50 {
			utils.SendErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name must be 50 characters or less")
			return
		}
		if input.Description != nil && len(*input.Description) > 255 {
			utils.SendErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Description must be 255 characters or less")
			return
		}

		// デフォルト値の設定
		defaultColor := "#6c757d"
		if input.Color == nil {
			input.Color = &defaultColor
		}

		// カラーコードのバリデーション（簡易版）
		if len(*input.Color) != 7 || (*input.Color)[0] != '#' {
			utils.SendErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Color must be in #RRGGBB format")
			return
		}

		// カテゴリを作成
		createBuilder := client.Category.Create().
			SetName(input.Name).
			SetColor(*input.Color)

		if input.Description != nil && *input.Description != "" {
			createBuilder.SetDescription(*input.Description)
		}

		category, err := createBuilder.Save(context.Background())
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to create category")
			return
		}

		response := utils.ConvertToCategoryResponse(category)
		utils.SendJSONResponse(w, http.StatusCreated, response)
	}
}

// GetCategoryByID は特定のカテゴリの詳細を取得するハンドラー
func GetCategoryByID(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLパラメータからIDを取得
		categoryIDStr := chi.URLParam(r, "categoryId")
		categoryID, ok := utils.ParseUUID(w, categoryIDStr)
		if !ok {
			return
		}

		// カテゴリを取得
		category, err := client.Category.Get(context.Background(), categoryID)
		if err != nil {
			if ent.IsNotFound(err) {
				utils.SendErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
			} else {
				utils.SendErrorResponse(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to fetch category")
			}
			return
		}

		response := utils.ConvertToCategoryResponse(category)
		utils.SendJSONResponse(w, http.StatusOK, response)
	}
}

// UpdateCategory はカテゴリ情報を更新するハンドラー
func UpdateCategory(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLパラメータからIDを取得
		categoryIDStr := chi.URLParam(r, "categoryId")
		categoryID, ok := utils.ParseUUID(w, categoryIDStr)
		if !ok {
			return
		}

		// リクエストボディをパース
		var input types.CategoryInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}

		// バリデーション
		if input.Name == "" {
			utils.SendErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name is required")
			return
		}
		if len(input.Name) > 50 {
			utils.SendErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name must be 50 characters or less")
			return
		}
		if input.Description != nil && len(*input.Description) > 255 {
			utils.SendErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Description must be 255 characters or less")
			return
		}
		if input.Color != nil {
			if len(*input.Color) != 7 || (*input.Color)[0] != '#' {
				utils.SendErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Color must be in #RRGGBB format")
				return
			}
		}

		// カテゴリの存在確認
		exists, err := client.Category.Query().Where(category.ID(categoryID)).Exist(context.Background())
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to check category existence")
			return
		}
		if !exists {
			utils.SendErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
			return
		}

		// カテゴリを更新
		updateBuilder := client.Category.UpdateOneID(categoryID).SetName(input.Name)

		if input.Description != nil {
			if *input.Description == "" {
				updateBuilder.ClearDescription()
			} else {
				updateBuilder.SetDescription(*input.Description)
			}
		}

		if input.Color != nil {
			updateBuilder.SetColor(*input.Color)
		}

		category, err := updateBuilder.Save(context.Background())
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to update category")
			return
		}

		response := utils.ConvertToCategoryResponse(category)
		utils.SendJSONResponse(w, http.StatusOK, response)
	}
}

// DeleteCategory はカテゴリを削除するハンドラー
func DeleteCategory(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLパラメータからIDを取得
		categoryIDStr := chi.URLParam(r, "categoryId")
		categoryID, ok := utils.ParseUUID(w, categoryIDStr)
		if !ok {
			return
		}

		// カテゴリを削除（関連するTodoのcategory_idは自動的にNULLになる）
		err := client.Category.DeleteOneID(categoryID).Exec(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				utils.SendErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
			} else {
				utils.SendErrorResponse(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to delete category")
			}
			return
		}

		// 204 No Content を返す
		w.WriteHeader(http.StatusNoContent)
	}
}