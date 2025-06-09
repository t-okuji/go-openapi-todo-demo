# ADR-001: Todo API パス拡張計画

## ステータス
提案済み

## 背景
現在のTodo APIは基本的なCRUD操作のみを提供している。実用的なTodo管理アプリケーションには、フィルタリング、検索、一括操作、統計情報などの機能が必要である。

## 現在のAPIエンドポイント
1. `GET /todos` - Todo一覧取得
2. `POST /todos` - Todo作成
3. `GET /todos/{todoId}` - Todo詳細取得
4. `PUT /todos/{todoId}` - Todo更新
5. `DELETE /todos/{todoId}` - Todo削除

## 決定事項

### 1. フィルタリング・検索機能の拡張
既存の `GET /todos` エンドポイントに以下のクエリパラメータを追加:

- `completed` (boolean): 完了状態でフィルタ
- `search` (string): タイトル・説明でキーワード検索
- `limit` (integer): 取得件数 (1-100, default: 20)
- `offset` (integer): オフセット (minimum: 0, default: 0)
- `sort` (enum): ソート項目 [createdAt, updatedAt, title] (default: createdAt)
- `order` (enum): ソート順 [asc, desc] (default: desc)

**使用例:**
```
GET /todos?completed=false&limit=10
GET /todos?search=買い物&sort=title&order=asc
```

### 2. 一括操作エンドポイントの追加
新しいエンドポイント `/todos/bulk` を追加:

#### `PATCH /todos/bulk` - 一括更新
```yaml
requestBody:
  type: object
  properties:
    todoIds:
      type: array
      items:
        type: string
    updates:
      type: object
      properties:
        completed:
          type: boolean
```

#### `DELETE /todos/bulk` - 一括削除
```yaml
requestBody:
  type: object
  properties:
    todoIds:
      type: array
      items:
        type: string
```

#### `POST /todos/bulk` - 一括作成
```yaml
requestBody:
  type: array
  items:
    $ref: "#/components/schemas/TodoInput"
```

### 3. 統計・メタデータエンドポイントの追加

#### `GET /todos/stats` - 統計情報取得
レスポンス:
```yaml
type: object
properties:
  total:
    type: integer
    description: Todo総数
  completed:
    type: integer
    description: 完了済み数
  pending:
    type: integer
    description: 未完了数
  createdToday:
    type: integer
    description: 今日作成された数
  completedToday:
    type: integer
    description: 今日完了した数
```

#### `GET /todos/count` - 件数取得
クエリパラメータ:
- `completed` (boolean): 完了状態でフィルタ

レスポンス:
```yaml
type: object
properties:
  count:
    type: integer
```

## 最終的なAPI構成
合計10のエンドポイント:

**基本CRUD (既存):**
1. `GET /todos` (拡張)
2. `POST /todos`
3. `GET /todos/{todoId}`
4. `PUT /todos/{todoId}`
5. `DELETE /todos/{todoId}`

**一括操作 (新規):**
6. `PATCH /todos/bulk`
7. `DELETE /todos/bulk`
8. `POST /todos/bulk`

**統計・メタデータ (新規):**
9. `GET /todos/stats`
10. `GET /todos/count`

## 実装ファイル構成
- `paths/todos.yml` - 既存ファイル（パラメータ拡張）
- `paths/todos-id.yml` - 既存ファイル（変更なし）
- `paths/todos-bulk.yml` - 新規作成
- `paths/todos-stats.yml` - 新規作成
- `openapi.yml` - 新規パス追加

## 利点
1. **実用性**: 実際のTodoアプリに必要な機能を網羅
2. **パフォーマンス**: 一括操作により効率的な処理が可能
3. **ユーザビリティ**: フィルタリング・検索により使いやすさ向上
4. **モニタリング**: 統計情報により使用状況を把握可能
5. **スケーラビリティ**: ページネーションにより大量データに対応

## 影響
- 既存エンドポイントへの後方互換性は維持
- 新機能は段階的に実装可能
- フロントエンド側で豊富な機能を提供可能

## 関連事項
- 実装時はセキュリティ（認証・認可）の考慮が必要
- レート制限の検討が必要
- 将来的にはタグ機能やカテゴリ機能の追加も検討