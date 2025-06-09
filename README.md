# OpenAPI Demo - Todo管理API

OpenAPI仕様を使用したTodo管理APIのデモプロジェクトです。カテゴリ機能を含む完全なTodo管理システムのAPIを提供します。

## 概要

このプロジェクトは、OpenAPI 3.1.1仕様に基づいたREST APIの実装例です。Todo管理とカテゴリ管理の機能を提供し、モジュール化されたAPI仕様設計のベストプラクティスを示しています。

## 機能

### Todo管理
- Todoの作成、取得、更新、削除（CRUD操作）
- カテゴリによる分類機能
- 完了状態の管理

### カテゴリ管理
- カテゴリの作成、取得、更新、削除
- カテゴリごとの色分け機能
- Todoとの関連付け

## 技術スタック

- **言語**: Go 1.24.4
- **Webフレームワーク**: Chi v5
- **API仕様**: OpenAPI 3.1.1

## プロジェクト構造

```
openapi-demo/
├── main.go                    # メインアプリケーション
├── openapi.yml               # OpenAPIメイン定義
├── paths/                    # APIエンドポイント定義
│   ├── todos.yml            # /todos エンドポイント
│   ├── todos-id.yml         # /todos/{todoId} エンドポイント
│   ├── categories.yml       # /categories エンドポイント
│   └── categories-id.yml    # /categories/{categoryId} エンドポイント
├── components/              # 共通コンポーネント
│   └── schemas/            # データモデル定義
│       ├── todo.yml        # Todoスキーマ
│       └── category.yml    # Categoryスキーマ
└── architecture-decisions/  # アーキテクチャ決定記録
    ├── 001-append-path.md  # API拡張計画
    └── 001-category.md     # カテゴリ機能追加計画
```

## セットアップ

### 必要な環境
- Go 1.24.4以上

### インストール
```bash
# リポジトリのクローン
git clone [repository-url]
cd openapi-demo

# 依存関係のインストール
go mod download
```

## 使用方法

### サーバーの起動
```bash
go run main.go
```

サーバーは `http://localhost:8080` で起動します。

## API仕様

### Todo API

#### エンドポイント一覧
- `GET /todos` - Todo一覧の取得
- `POST /todos` - 新規Todoの作成
- `GET /todos/{todoId}` - 特定のTodoの取得
- `PUT /todos/{todoId}` - Todoの更新
- `DELETE /todos/{todoId}` - Todoの削除

#### Todoデータモデル
```yaml
Todo:
  - id: string (UUID)
  - title: string (必須)
  - description: string
  - completed: boolean
  - categoryId: string (UUID)
  - createdAt: string (date-time)
  - updatedAt: string (date-time)
```

### Category API

#### エンドポイント一覧
- `GET /categories` - カテゴリ一覧の取得
- `POST /categories` - 新規カテゴリの作成
- `GET /categories/{categoryId}` - 特定のカテゴリの取得
- `PUT /categories/{categoryId}` - カテゴリの更新
- `DELETE /categories/{categoryId}` - カテゴリの削除

#### Categoryデータモデル
```yaml
Category:
  - id: string (UUID)
  - name: string (必須)
  - description: string
  - color: string (HEX形式)
  - createdAt: string (date-time)
  - updatedAt: string (date-time)
```

## 今後の拡張予定

### Todo APIの拡張（計画中）
- フィルタリング・検索機能の追加
- 一括操作エンドポイントの実装
- 統計情報・メタデータAPIの追加

詳細は `architecture-decisions/001-append-path.md` を参照してください。

## 開発

### API仕様の確認
OpenAPI仕様ファイルは `openapi.yml` および `paths/`、`components/` ディレクトリ内のファイルで定義されています。
