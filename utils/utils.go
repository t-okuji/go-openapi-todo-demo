package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/t-okuji/go-openapi-todo-demo/ent"
	"github.com/t-okuji/go-openapi-todo-demo/types"
)

// SendJSONResponse は JSON レスポンスを送信する共通関数
func SendJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encoding error: %v", err)
	}
}

// SendErrorResponse はエラーレスポンスを送信する共通関数
func SendErrorResponse(w http.ResponseWriter, status int, code string, message string) {
	var errResp types.ErrorResponse
	errResp.Error.Code = code
	errResp.Error.Message = message
	SendJSONResponse(w, status, errResp)
}

// ParseUUID は文字列をUUIDとしてパースし、エラーがあればエラーレスポンスを送信する
func ParseUUID(w http.ResponseWriter, uuidStr string) (uuid.UUID, bool) {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "INVALID_UUID", "Invalid UUID format")
		return uuid.UUID{}, false
	}
	return id, true
}

// ConvertToTodoResponse は Ent の Todo エンティティを TodoResponse に変換する
func ConvertToTodoResponse(todo *ent.Todo) types.TodoResponse {
	response := types.TodoResponse{
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
