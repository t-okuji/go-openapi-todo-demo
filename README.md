# Go OpenAPI Todo Demo

OpenAPI仕様ベースのTodo管理REST API - Go言語とEntフレームワークによる実装。カテゴリ機能付きの完全なCRUD操作を提供する本格的なAPIサーバーです。

## 概要

このプロジェクトは、OpenAPI 3.1.1仕様に基づいたREST APIの実装例です。以下の特徴を持つ、プロダクションレディなAPIサーバーの実装パターンを示しています：

- **モジュール化されたOpenAPI仕様**: パス定義とスキーマ定義を分離した保守性の高い構造
- **型安全なORM**: Entフレームワークによる型安全なデータベース操作
- **構造化されたエラーハンドリング**: 統一されたエラーレスポンス形式
- **RESTful設計**: 適切なHTTPステータスコードとレスポンス形式
- **Docker化された開発環境**: PostgreSQLとの統合開発環境

## 実装状況

- ✅ **実装済み**: Todo API完全実装（全CRUD操作）
- ✅ **実装済み**: Category API完全実装（全CRUD操作）

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
- **ORM**: Ent（EntGo）
- **環境変数管理**: godotenv
- **データベース**: PostgreSQL 17 Alpine
- **コンテナ**: Docker/Docker Compose
- **API仕様**: OpenAPI 3.1.1

## プロジェクト構造

```
go-openapi-todo-demo/
├── main.go                    # メインアプリケーション（Chi + Ent + godotenv）
├── .env                       # 環境変数設定（DB接続情報）
├── handlers/                  # HTTPハンドラー実装
│   └── todo.go               # Todo/Categoryハンドラー
├── types/                     # 型定義
│   └── types.go              # APIリクエスト/レスポンス型
├── utils/                     # ユーティリティ関数
│   └── utils.go              # 共通処理（レスポンス、UUID検証等）
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
├── openapi-ui.html          # Stoplight Elements UI
└── architecture-decisions/  # アーキテクチャ決定記録
```

## セットアップ

### 必要な環境
- Go 1.24.4以上
- Docker & Docker Compose

### インストールと環境セットアップ
```bash
# リポジトリのクローン
git clone [repository-url]
cd go-openapi-todo-demo

# 依存関係のインストール
go mod download

# データベース起動
docker compose up -d

# データベース初期化
docker compose exec postgres psql -U user -d demo -f /db/schema.sql

# サンプルデータ投入（オプション）
docker compose exec postgres psql -U user -d demo -f /db/seed.sql
```

### 環境変数設定
`.env`ファイルにデータベース接続情報が設定済みです：
```env
POSTGRES_DB=demo
POSTGRES_PORT=15432
POSTGRES_HOST=localhost
POSTGRES_USER=user
POSTGRES_PASSWORD=password
```

## 使用方法

### サーバーの起動
```bash
go run main.go
```

サーバーは `http://localhost:8080` で起動します。

### API テスト

#### 基本的な動作確認
```bash
# ヘルスチェック
curl http://localhost:8080/

# APIドキュメント（Stoplight Elements）
open http://localhost:8080/docs
```

#### Todo API
```bash
# Todo一覧取得
curl http://localhost:8080/todos

# 新規Todo作成
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "買い物", "description": "牛乳を買う"}'

# Todo更新（完了状態の変更）
curl -X PUT http://localhost:8080/todos/{todoId} \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

#### Category API
```bash
# カテゴリ一覧取得
curl http://localhost:8080/categories

# 新規カテゴリ作成
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "仕事", "color": "#0066cc"}'
```

## API仕様

### Todo API

#### エンドポイント一覧
- ✅ `GET /todos` - Todo一覧の取得（実装済み）
- ✅ `POST /todos` - 新規Todoの作成（実装済み）
- ✅ `GET /todos/{todoId}` - 特定のTodoの取得（実装済み）
- ✅ `PUT /todos/{todoId}` - Todoの更新（実装済み）
- ✅ `DELETE /todos/{todoId}` - Todoの削除（実装済み）

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
- ✅ `GET /categories` - カテゴリ一覧の取得（実装済み）
- ✅ `POST /categories` - 新規カテゴリの作成（実装済み）
- ✅ `GET /categories/{categoryId}` - 特定のカテゴリの取得（実装済み）
- ✅ `PUT /categories/{categoryId}` - カテゴリの更新（実装済み）
- ✅ `DELETE /categories/{categoryId}` - カテゴリの削除（実装済み）

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
OpenAPI仕様ファイルは `openapi/openapi.yml` および `openapi/paths/`、`openapi/components/` ディレクトリ内のファイルで定義されています。

### 開発コマンド
```bash
# コード品質チェック
go vet ./...
go fmt ./...

# アプリケーションのビルド
go build -v .

# 依存関係の整理
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

## アーキテクチャの特徴

### 実装パターン
- **ハンドラー分離**: エンティティごとに独立したハンドラーファイル
- **共通エラー処理**: 構造化されたエラーコードとレスポンス
- **UUID主キー**: PostgreSQLの`gen_random_uuid()`による自動生成
- **自動タイムスタンプ**: トリガーによる`updated_at`の自動更新

### セキュリティとバリデーション
- UUID形式の検証
- 文字列長の制限（カテゴリ名: 50文字、説明: 255文字）
- HEXカラーコードのバリデーション
- 外部キー制約による参照整合性

