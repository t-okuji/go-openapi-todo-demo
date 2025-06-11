package types

import "time"

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
