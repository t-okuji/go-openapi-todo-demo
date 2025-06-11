# Todo API 処理フロー図

## 1. GET /todos - Todo一覧取得

```mermaid
flowchart TD
    A[クライアント: GET /todos] --> B[DBからTodo一覧取得]
    B --> C{成功？}
    C -->|はい| D[200 OK: Todo一覧をJSON返却]
    C -->|いいえ| E[500: エラーレスポンス]
```

## 2. POST /todos - Todo作成

```mermaid
flowchart TD
    A[クライアント: POST /todos] --> B[リクエストボディ解析]
    B --> C{解析成功？}
    C -->|いいえ| D[400: 不正なリクエスト]
    C -->|はい| E[タイトル必須チェック]
    E --> F{カテゴリID指定？}
    F -->|はい| G[UUID形式検証]
    G --> H{有効？}
    H -->|いいえ| I[400: 無効なUUID]
    F -->|いいえ| J[DB保存]
    H -->|はい| J[DB保存]
    J --> K{成功？}
    K -->|はい| L[201: 作成したTodoを返却]
    K -->|いいえ| M[500: エラーレスポンス]
```

## 3. GET /todos/{todoId} - 特定Todo取得

```mermaid
flowchart TD
    A["クライアント: GET /todos/{todoId}"] --> B[todoId取得]
    B --> C[UUID形式検証]
    C --> D{有効？}
    D -->|いいえ| E["400: 無効なUUID"]
    D -->|はい| F[DBからTodo取得]
    F --> G{存在する？}
    G -->|はい| H["200: Todoを返却"]
    G -->|いいえ| I["404: Todo未発見"]
```

## 4. PUT /todos/{todoId} - Todo更新

```mermaid
flowchart TD
    A["クライアント: PUT /todos/{todoId}"] --> B[todoId検証]
    B --> C{有効なUUID？}
    C -->|いいえ| D["400: 無効なUUID"]
    C -->|はい| E[リクエストボディ解析]
    E --> F{解析成功？}
    F -->|いいえ| G["400: 不正なリクエスト"]
    F -->|はい| H[更新フィールド設定]
    H --> I[DB更新実行]
    I --> J{Todo存在？}
    J -->|はい| K["200: 更新後のTodo返却"]
    J -->|いいえ| L["404: Todo未発見"]
```

## 5. DELETE /todos/{todoId} - Todo削除

```mermaid
flowchart TD
    A["クライアント: DELETE /todos/{todoId}"] --> B[todoId検証]
    B --> C{有効なUUID？}
    C -->|いいえ| D["400: 無効なUUID"]
    C -->|はい| E[DB削除実行]
    E --> F{Todo存在？}
    F -->|はい| G["204: No Content"]
    F -->|いいえ| H["404: Todo未発見"]
```

## HTTPステータスコード一覧

| ステータス | 説明 | 使用場面 |
|---------|------|---------|
| 200 | OK | 取得・更新成功 |
| 201 | Created | 作成成功 |
| 204 | No Content | 削除成功 |
| 400 | Bad Request | 不正なリクエスト/UUID |
| 404 | Not Found | Todo未発見 |
| 500 | Internal Server Error | サーバーエラー |