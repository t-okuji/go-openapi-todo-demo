# アーキテクチャ決定記録: Category API実装

## ステータス
承認済み

## コンテキスト
現在、Todo管理APIシステムにはTodo APIが実装されているが、Category APIは未実装である。OpenAPI仕様書では両方のAPIが定義されており、データベーススキーマとEntスキーマも既に存在している。システムの完全性を保つため、Category APIの実装が必要である。

## 決定事項

### 1. 実装アプローチ
既存のTodo API実装パターンに従い、一貫性のあるコードベースを維持する。

### 2. データ構造

#### CategoryResponse構造体
```go
type CategoryResponse struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    Description *string    `json:"description,omitempty"`
    Color       string     `json:"color"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}
```

#### CategoryInput構造体
```go
type CategoryInput struct {
    Name        string  `json:"name"`
    Description *string `json:"description,omitempty"`
    Color       *string `json:"color,omitempty"`
}
```

### 3. 実装するエンドポイント

1. **GET /categories**
   - 全カテゴリの一覧を取得
   - レスポンス: CategoryResponse配列

2. **POST /categories**
   - 新規カテゴリを作成
   - リクエスト: CategoryInput
   - レスポンス: CategoryResponse（201 Created）

3. **GET /categories/{categoryId}**
   - 特定カテゴリの詳細を取得
   - レスポンス: CategoryResponse

4. **PUT /categories/{categoryId}**
   - カテゴリ情報を更新
   - リクエスト: CategoryInput
   - レスポンス: CategoryResponse

5. **DELETE /categories/{categoryId}**
   - カテゴリを削除
   - レスポンス: 204 No Content

### 4. 実装の詳細

#### 共通関数の利用
- `sendJSONResponse`: JSON形式のレスポンス送信
- `sendErrorResponse`: エラーレスポンスの統一的な送信
- `parseUUID`: UUID文字列の検証とパース

#### エラーハンドリング
- 400 Bad Request: 不正なリクエスト形式、UUID形式エラー
- 404 Not Found: 指定されたカテゴリが存在しない
- 500 Internal Server Error: データベースエラー

#### バリデーション
- カテゴリ名: 必須、1文字以上50文字以下
- 説明: オプション、最大255文字
- 色: HEX形式（#RRGGBB）、デフォルト値 #6c757d

### 5. データベース制約の考慮

#### カテゴリ削除時の影響
- `ON DELETE SET NULL`制約により、カテゴリ削除時に関連するTodoのcategory_idはNULLに設定される
- 明示的なTodo更新処理は不要（データベースが自動処理）

#### トランザクション
- 単一のCRUD操作のため、Entのデフォルトトランザクション処理で十分
- 複雑なトランザクションは現時点では不要

### 6. 実装順序

1. データ構造の定義（CategoryResponse、CategoryInput）
2. 変換関数の実装（convertToCategoryResponse）
3. 各ハンドラー関数の実装
4. main.goへのルーティング追加
5. 動作テスト

## 結果

### 利点
- 既存コードとの一貫性維持
- OpenAPI仕様への完全準拠
- 明確なエラーハンドリング
- 保守性の高いコード構造

### 欠点
- 特になし（標準的な実装アプローチ）

### リスク
- カテゴリ削除時のTodoへの影響をユーザーが理解していない可能性
  - 対策: APIドキュメントに明記

## 今後の検討事項
- カテゴリごとのTodo数を返すエンドポイントの追加
- カテゴリのソート機能
- ページネーション（カテゴリ数が増加した場合）