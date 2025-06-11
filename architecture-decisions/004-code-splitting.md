# ADR-004: main.goファイル分割（コード分離）

## ステータス
提案済み

## コンテキスト

現在の`main.go`ファイルは418行となり、以下の問題が発生している：

- **単一責任原則の違反**: HTTP ハンドラー、ユーティリティ関数、型定義、サーバー起動ロジックが混在
- **保守性の低下**: 機能追加時（Category API等）のファイル肥大化
- **テスタビリティの問題**: 関数が private のため単体テストが困難
- **可読性の悪化**: 関連機能が散在し、コード理解が困難

## 決定事項

main.goをエンドポイント単位で分割し、以下の構造に再編成する：

### 新しいファイル構成

```
/
├── main.go                    # サーバー起動・ルーティング・DB接続のみ
├── types/
│   └── types.go              # 共通型定義
├── utils/
│   └── utils.go              # 共通ユーティリティ関数
└── handlers/
    └── todo.go               # Todo API ハンドラー
```

### 責任分離

#### **types/types.go**
- API レスポンス・リクエスト型の定義
- エラー構造体の定義
- 他パッケージに依存しない純粋な型定義層

**含まれる型**:
```go
- TodoResponse    // API レスポンス用 Todo エンティティ
- TodoInput       // API リクエスト用 Todo 入力データ
- ErrorResponse   // API エラーレスポンス
```

#### **utils/utils.go**
- HTTP レスポンス関連の共通ロジック
- データ変換・検証ロジック
- typesパッケージのみに依存

**含まれる関数**:
```go
- SendJSONResponse()        // JSON レスポンス送信
- SendErrorResponse()       // エラーレスポンス送信
- ParseUUID()              // UUID パース・検証
- ConvertToTodoResponse()  // Ent→API レスポンス変換
```

#### **handlers/todo.go**
- Todo API の HTTP ハンドラー実装
- ビジネスロジックの実装
- utils、typesパッケージに依存

**含まれる関数**:
```go
- GetTodosHandler()      // GET /todos
- CreateTodoHandler()    // POST /todos
- GetTodoByIDHandler()   // GET /todos/{todoId}
- UpdateTodoHandler()    // PUT /todos/{todoId}
- DeleteTodoHandler()    // DELETE /todos/{todoId}
```

#### **main.go（リファクタ後）**
- アプリケーション起動処理
- データベース接続設定
- ルーティング設定
- CORS・ミドルウェア設定

**含まれる関数**:
```go
- Open()  // データベース接続
- main()  // アプリケーションエントリーポイント
```

### 依存関係

```
main.go → handlers/todo.go → utils/utils.go → types/types.go
```

- **types**: 最下位層、他パッケージに依存しない
- **utils**: types パッケージのみに依存
- **handlers**: utils、types パッケージに依存
- **main**: handlers パッケージに依存

### 関数の可視性変更

現在の private 関数（小文字開始）を public 関数（大文字開始）に変更：

```go
// 変更前 → 変更後
sendJSONResponse      → SendJSONResponse
sendErrorResponse     → SendErrorResponse
parseUUID            → ParseUUID
convertToTodoResponse → ConvertToTodoResponse
getTodosHandler      → GetTodosHandler
createTodoHandler    → CreateTodoHandler
// etc...
```

## 利点

### 1. 保守性の向上
- 機能別のファイル分離により、変更影響範囲が明確
- 新機能（Category API等）追加時の影響最小化

### 2. テスタビリティの向上
- public 関数化により単体テスト作成が容易
- Mock 作成時の依存関係が明確

### 3. 可読性の向上
- 関連機能の集約により、コード理解が容易
- ファイルサイズの適正化（各ファイル100-200行程度）

### 4. 再利用性の向上
- utils パッケージの関数が他エンドポイントでも利用可能
- types パッケージが共通データ構造として機能

### 5. 将来の拡張性
- Category API実装時の追加が容易
- 新しいエンドポイント追加時のボイラープレート削減

## 代替案

### A. 単一ファイル維持
- **却下理由**: ファイル肥大化が継続、保守性問題が未解決

### B. 機能別ディレクトリ分割
```
/api/
  /todo/
    handler.go
    types.go
  /category/
    handler.go
    types.go
/shared/
  utils.go
```
- **却下理由**: 小規模プロジェクトには過度な分割

### C. レイヤー別分割
```
/models/
/controllers/
/services/
```
- **却下理由**: Go の慣用的なパッケージ構成ではない

## 実装方針

### フェーズ1: Todo API分割（現在対象）
1. types/types.go 作成
2. utils/utils.go 作成  
3. handlers/todo.go 作成
4. main.go リファクタリング

### フェーズ2: Category API実装時
1. Category関連型をtypes/types.goに追加
2. handlers/category.go 作成
3. main.goにルーティング追加

### インポート最適化
- 各ファイルで必要なimportのみ保持
- 循環依存の回避
- モジュール内依存関係の明確化

## 影響

### 正の影響
- コードの整理により開発効率向上
- 新規開発者のコード理解時間短縮
- テストカバレッジ向上の基盤整備

### 潜在的リスク
- 一時的な開発工数増加
- パッケージ間のインターフェース設計要注意

### 軽減策
- 段階的な分割実施
- 十分なテスト実行による動作保証
- import文の慎重な管理

## 実装完了の定義

- [ ] types/types.go 作成完了
- [ ] utils/utils.go 作成完了
- [ ] handlers/todo.go 作成完了
- [ ] main.go リファクタリング完了
- [ ] 全機能の動作確認完了
- [ ] コンパイルエラーなし
- [ ] 既存のTodo API動作保証

この分割により、プロジェクトの長期的な保守性と拡張性が大幅に向上する。