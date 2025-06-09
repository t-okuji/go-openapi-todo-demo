# ADR-001: Category API追加計画

## 概要
OpenAPI Todo仕様にCategory管理機能を追加する計画書。

## 背景
現在のTodo APIにカテゴリ機能を追加することで、Todoの分類・整理を可能にする。

## 決定事項

### 実装計画

#### 1. **高優先度** - スキーマ設計
- **Categoryスキーマ**: カテゴリのメインデータ構造を定義
  - 必須フィールド: id, name, createdAt
  - オプション: description, color, updatedAt
- **CategoryInputスキーマ**: 作成・更新時の入力用スキーマを定義
  - 必須: name のみ

#### 2. **中優先度** - APIエンドポイント実装
- **`/categories`パス**: 
  - `GET` - カテゴリ一覧取得
  - `POST` - カテゴリ作成
- **`/categories/{categoryId}`パス**:
  - `GET` - カテゴリ詳細取得
  - `PUT` - カテゴリ更新
  - `DELETE` - カテゴリ削除

#### 3. **中優先度** - メインファイル更新
- `openapi.yml`にカテゴリパスを追加

#### 4. **低優先度** - 既存スキーマ拡張検討
- Todoスキーマにカテゴリ関連フィールド追加を検討
  - categoryId: カテゴリとの関連付け

## 設計原則
- 既存のTodo APIと同様の構造を維持
- モジュール化された設計を継続
- 日本語対応を継続
- RESTful APIの原則に従う

## ファイル構成
```
components/
  schemas/
    category.yml  # 新規追加
    todo.yml      # categoryId追加検討
paths/
  categories.yml     # 新規追加
  categories-id.yml  # 新規追加
openapi.yml          # パス追加
```

## 想定される影響
- 新規機能追加のため、既存APIへの影響は最小限
- Todo-Category間の関連性により、データ整合性の考慮が必要

## 代替案
1. Todoスキーマ内にcategoryフィールドを直接追加（正規化なし）
2. 別途Category管理サービスとして分離

## 状況
- **ステータス**: 承認済み
- **決定日**: 2025-01-09
- **決定者**: 開発チーム