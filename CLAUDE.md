# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

OpenAPIベースのTodo管理API - Go言語とEntフレームワークを使用した、カテゴリ機能付きRESTful APIの実装。Chi v5ルーター、PostgreSQL、Docker Composeによる開発環境構築済み。

## 開発コマンド

### 環境セットアップ
```bash
# 依存関係のインストール
go mod download

# 環境変数ファイルの設定（必要に応じて編集）
# .envファイルにデータベース接続情報が設定済み

# データベース起動（Docker Compose）
docker compose up -d

# データベース初期化
docker compose exec postgres psql -U user -d demo -f /db/schema.sql

# サンプルデータ投入（オプション）
docker compose exec postgres psql -U user -d demo -f /db/seed.sql
```

### サーバー起動
```bash
go run main.go
```
- ポート8080で起動
- http://localhost:8080

### テスト・検証
```bash
# アプリケーションのビルド
go build -v .

# コード品質チェック
go vet ./...
go fmt ./...

# 基本的な動作確認
curl http://localhost:8080/

# Todo API のテスト
curl http://localhost:8080/todos

# Category API のテスト
curl http://localhost:8080/categories

# 依存関係の確認・整理
go list -m all
go mod tidy
go mod verify

# Entコード生成（スキーマ変更時）
go generate ./ent
```

### データベース管理
```bash
# データベース接続
docker compose exec postgres psql -U user -d demo

# データベース停止
docker compose down

# データベース完全削除（ボリューム含む）
docker compose down -v
```

## アーキテクチャ

### 技術スタック
- **言語**: Go 1.24.4
- **Webフレームワーク**: Chi v5
- **ORM**: Ent（EntGo）
- **環境変数管理**: godotenv
- **データベース**: PostgreSQL 17 Alpine
- **コンテナ**: Docker/Docker Compose
- **API仕様**: OpenAPI 3.1.1

### プロジェクト構造
```
go-openapi-todo-demo/
├── main.go                    # メインアプリケーション（Chi + Ent + godotenv）
├── .env                       # 環境変数設定（DB接続情報）
├── ent/                       # Entコード生成済みORM
│   ├── schema/               # Entスキーマ定義
│   ├── todo.go              # Todoエンティティ
│   └── category.go          # Categoryエンティティ
├── openapi/                   # OpenAPI仕様定義
│   ├── openapi.yml           # メインAPI定義
│   ├── paths/                # APIエンドポイント定義
│   └── components/schemas/   # データモデル定義
├── db/                       # データベース関連
│   ├── schema.sql           # 完全スキーマ定義
│   ├── seed.sql             # サンプルデータ
│   └── migrations/          # マイグレーションファイル
├── compose.yml              # Docker Compose設定
└── architecture-decisions/  # アーキテクチャ決定記録
```

### データベース設計
- **PostgreSQL 17使用**: Docker Composeで管理
- **UUID主キー**: gen_random_uuid()で自動生成
- **自動タイムスタンプ**: created_at/updated_atの自動管理
- **外部キー制約**: todos.category_id → categories.id
- **CHECK制約**: データ検証（カテゴリ名長さ、色形式等）
- **インデックス**: 検索性能最適化

### OpenAPI仕様の構造
モジュール化されたOpenAPI 3.1.1仕様を採用：
- `openapi/openapi.yml`: メインのAPI定義（他ファイルを$refで参照）
- `openapi/paths/`: エンドポイント定義を個別ファイルに分離
- `openapi/components/schemas/`: データモデル定義を個別ファイルに分離

### API構成
1. **Todo API** (`/todos`, `/todos/{todoId}`)
   - 基本的なCRUD操作
   - カテゴリとの連携（categoryIdフィールド）

2. **Category API** (`/categories`, `/categories/{categoryId}`)
   - カテゴリのCRUD操作
   - Todo管理のためのカテゴリ分類機能

### 実装パターンと規約

#### ハンドラー実装
- ハンドラーはエンティティごとに`handlers/`ディレクトリに分離
- 各エンドポイントは独立した関数として実装（例: `GetTodos`, `CreateTodo`）
- コンテキストは`context.Background()`を使用
- バリデーションはハンドラー内で実装

#### レスポンス処理
- 共通関数を使用: `utils.SendJSONResponse()`, `utils.SendErrorResponse()`
- エラーコード体系: `DB_ERROR`, `INVALID_UUID`, `NOT_FOUND`等の構造化されたコード
- HTTPステータス: RESTful規約準拠（201 Created、204 No Content等）
- JSONフィールド: OpenAPI仕様に合わせてcamelCase（例: `categoryId`）

#### データベース操作
- Entクライアントは`main.go`で初期化し、ハンドラーにグローバル変数として渡す
- UUID検証は`utils.ParseUUID()`を使用
- 更新・削除前に必ず存在確認を実施
- 空文字列でフィールドをクリア可能（オプショナルフィールド）

#### 型変換
- Entエンティティ→APIレスポンス変換: `utils.TodoToResponse()`, `utils.CategoryToResponse()`
- 入力データ型は`types/types.go`に定義
- オプショナルフィールドはポインタ型で実装

### 環境変数設定
.envファイルに以下の変数を設定：
```
POSTGRES_DB=demo
POSTGRES_PORT=15432
POSTGRES_HOST=localhost
POSTGRES_USER=user
POSTGRES_PASSWORD=password
```

### データベースマイグレーション
- マイグレーションファイル: `db/migrations/`ディレクトリ
- 命名規則: `001_initial_schema.up.sql`, `001_initial_schema.down.sql`
- 自動タイムスタンプ更新: PostgreSQLトリガーで`updated_at`を管理
- 外部キー制約: `ON DELETE SET NULL`でカテゴリ削除時の処理

### 開発環境の要件
- Go 1.24.4以上
- Docker & Docker Compose
- PostgreSQL 17（Dockerで自動セットアップ）