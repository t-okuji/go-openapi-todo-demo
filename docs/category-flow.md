# Category API 処理フロー図

## 1. GET /categories - カテゴリ一覧取得

```mermaid
flowchart TD
    A[クライアント: GET /categories] --> B[DBからカテゴリ一覧取得]
    B --> C{成功？}
    C -->|はい| D[200 OK: カテゴリ一覧をJSON返却]
    C -->|いいえ| E[500: エラーレスポンス]
```

## 2. POST /categories - カテゴリ作成

```mermaid
flowchart TD
    A[クライアント: POST /categories] --> B[リクエストボディ解析]
    B --> C{解析成功？}
    C -->|いいえ| D[400: 不正なリクエスト]
    C -->|はい| E[name必須チェック]
    E --> F{nameあり？}
    F -->|いいえ| G[400: name必須]
    F -->|はい| H[color形式チェック]
    H --> I{有効な16進数カラー？}
    I -->|いいえ| J[400: 無効なカラー形式]
    I -->|はい| K[DB保存]
    K --> L{成功？}
    L -->|はい| M[201: 作成したカテゴリを返却]
    L -->|いいえ| N[500: エラーレスポンス]
```

## 3. GET /categories/{categoryId} - 特定カテゴリ取得

```mermaid
flowchart TD
    A["クライアント: GET /categories/{categoryId}"] --> B[categoryId取得]
    B --> C[UUID形式検証]
    C --> D{有効？}
    D -->|いいえ| E["400: 無効なUUID"]
    D -->|はい| F[DBからカテゴリ取得]
    F --> G{存在する？}
    G -->|はい| H["200: カテゴリを返却"]
    G -->|いいえ| I["404: カテゴリ未発見"]
```

## 4. PUT /categories/{categoryId} - カテゴリ更新

```mermaid
flowchart TD
    A["クライアント: PUT /categories/{categoryId}"] --> B[categoryId検証]
    B --> C{有効なUUID？}
    C -->|いいえ| D["400: 無効なUUID"]
    C -->|はい| E[リクエストボディ解析]
    E --> F{解析成功？}
    F -->|いいえ| G["400: 不正なリクエスト"]
    F -->|はい| H[name存在チェック]
    H --> I{nameあり？}
    I -->|いいえ| J["400: name必須"]
    I -->|はい| K[color形式チェック]
    K --> L{有効な16進数カラー？}
    L -->|いいえ| M["400: 無効なカラー形式"]
    L -->|はい| N[DB更新実行]
    N --> O{カテゴリ存在？}
    O -->|はい| P["200: 更新後のカテゴリ返却"]
    O -->|いいえ| Q["404: カテゴリ未発見"]
```

## 5. DELETE /categories/{categoryId} - カテゴリ削除

```mermaid
flowchart TD
    A["クライアント: DELETE /categories/{categoryId}"] --> B[categoryId検証]
    B --> C{有効なUUID？}
    C -->|いいえ| D["400: 無効なUUID"]
    C -->|はい| E[関連Todo確認]
    E --> F{関連Todoあり？}
    F -->|はい| G["400: 関連Todoあり削除不可"]
    F -->|いいえ| H[DB削除実行]
    H --> I{カテゴリ存在？}
    I -->|はい| J["204: No Content"]
    I -->|いいえ| K["404: カテゴリ未発見"]
```

## HTTPステータスコード一覧

| ステータス | 説明 | 使用場面 |
|---------|------|---------|
| 200 | OK | 取得・更新成功 |
| 201 | Created | 作成成功 |
| 204 | No Content | 削除成功 |
| 400 | Bad Request | 不正なリクエスト/UUID/カラー形式/関連Todo存在 |
| 404 | Not Found | カテゴリ未発見 |
| 500 | Internal Server Error | サーバーエラー |

## バリデーションルール

### name (カテゴリ名)
- 必須フィールド
- 1文字以上50文字以下
- 空文字列不可

### color (カラーコード)
- オプショナルフィールド
- 6桁の16進数形式（例: #FF5733, #00AA00）
- 先頭の#は必須
- 大文字小文字は区別しない

### categoryId
- UUID v4形式
- 例: 123e4567-e89b-12d3-a456-426614174000

## エラーレスポンス形式

```json
{
  "error": "エラーメッセージ"
}
```