# ADR-003: Todo API実装計画

## ステータス
- 日付: 2025-01-09
- ステータス: 承認済み

## コンテキスト
現在、Todo管理APIのうち`GET /todos`エンドポイントのみが実装されている。OpenAPI仕様では以下のエンドポイントが定義されているが、未実装である：
- `POST /todos` - Todo作成
- `GET /todos/{todoId}` - Todo詳細取得
- `PUT /todos/{todoId}` - Todo更新
- `DELETE /todos/{todoId}` - Todo削除

これらのエンドポイントを実装し、完全なCRUD機能を提供する必要がある。

## 決定事項

### 1. POST /todos - Todo作成API

#### 実装内容
- リクエストボディから`TodoInput`を受け取る
- 必須フィールド`title`の検証
- オプショナルフィールド（description、completed、categoryId）の処理
- categoryIdが指定された場合、存在確認を行う
- Entを使用してデータベースに保存
- 201ステータスで作成されたTodoを返却

#### エラーハンドリング
- 400 Bad Request: タイトルが空、不正なJSON、存在しないcategoryId
- 500 Internal Server Error: データベースエラー

### 2. GET /todos/{todoId} - Todo詳細取得API

#### 実装内容
- URLパラメータからtodoIdを取得
- UUID形式の検証
- Entを使用してIDでTodoを検索
- 見つかった場合は200で返却

#### エラーハンドリング
- 400 Bad Request: 不正なUUID形式
- 404 Not Found: 指定されたIDのTodoが存在しない
- 500 Internal Server Error: データベースエラー

### 3. PUT /todos/{todoId} - Todo更新API

#### 実装内容
- URLパラメータからtodoIdを取得
- リクエストボディから`TodoInput`を受け取る
- 対象Todoの存在確認
- 更新可能フィールド（title、description、completed、categoryId）の更新
- categoryIdが変更される場合、新しいカテゴリの存在確認
- updated_atは自動更新（データベーストリガーで処理）

#### エラーハンドリング
- 400 Bad Request: 不正なリクエスト（空のtitle、不正なcategoryId等）
- 404 Not Found: 更新対象のTodoが存在しない
- 500 Internal Server Error: データベースエラー

### 4. DELETE /todos/{todoId} - Todo削除API

#### 実装内容
- URLパラメータからtodoIdを取得
- UUID形式の検証
- 対象Todoの存在確認
- Entを使用してTodoを削除
- 204 No Contentを返却

#### エラーハンドリング
- 400 Bad Request: 不正なUUID形式
- 404 Not Found: 削除対象のTodoが存在しない
- 500 Internal Server Error: データベースエラー

### 5. 共通処理の設計

#### エラーレスポンス形式
```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "タイトルは必須です"
  }
}
```

#### 共通バリデーション
- UUID形式の検証関数
- 必須フィールドの検証
- 文字列長の制限（Entスキーマに基づく）

#### 共通ユーティリティ
- JSONレスポンス送信関数
- エラーレスポンス送信関数
- リクエストボディのパース関数

### 6. 実装順序
1. 共通エラーハンドリング関数の作成
2. POST /todos（作成）
3. GET /todos/{todoId}（詳細取得）
4. PUT /todos/{todoId}（更新）
5. DELETE /todos/{todoId}（削除）

## 実装方針

### 技術的な実装方針
- 既存の`GET /todos`の実装パターンに従う
- Chi v5のルーティング機能を活用
- EntのORMメソッドを使用したデータベース操作
- `TodoResponse`構造体を使用した一貫したレスポンス形式（camelCase）

### コード構造
```go
// 共通エラーレスポンス構造体
type ErrorResponse struct {
    Error struct {
        Code    string `json:"code"`
        Message string `json:"message"`
    } `json:"error"`
}

// TodoInput構造体（リクエスト用）
type TodoInput struct {
    Title       string  `json:"title"`
    Description *string `json:"description,omitempty"`
    Completed   *bool   `json:"completed,omitempty"`
    CategoryID  *string `json:"categoryId,omitempty"`
}
```

### バリデーション戦略
- リクエストボディのJSONパース時の基本的な検証
- ビジネスロジックレベルでの詳細な検証
- データベース制約による最終的な検証

## 結果として期待される効果
- 完全なCRUD機能を持つTodo API
- OpenAPI仕様に準拠した一貫性のあるAPI
- 適切なエラーハンドリングによる開発者体験の向上
- 既存コードとの一貫性を保った実装

## 参考資料
- OpenAPI仕様: `/openapi/openapi.yml`
- Entスキーマ定義: `/ent/schema/todo.go`
- 既存実装: `main.go`の`GET /todos`エンドポイント