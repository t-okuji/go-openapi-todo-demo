# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

OpenAPIを使用したTodo管理APIのデモプロジェクト。Go言語で実装され、カテゴリ機能付きのTodo管理システムを構築している。

## 開発コマンド

### サーバー起動
```bash
go run main.go
```
- ポート8080で起動
- http://localhost:8080

## アーキテクチャ

### OpenAPI仕様の構造
モジュール化されたOpenAPI 3.1.1仕様を採用：
- `openapi.yml`: メインのAPI定義（他ファイルを$refで参照）
- `paths/`: エンドポイント定義を個別ファイルに分離
- `components/schemas/`: データモデル定義を個別ファイルに分離

### API構成
1. **Todo API** (`/todos`, `/todos/{todoId}`)
   - 基本的なCRUD操作
   - カテゴリとの連携（categoryIdフィールド）

2. **Category API** (`/categories`, `/categories/{categoryId}`)
   - カテゴリのCRUD操作
   - Todo管理のためのカテゴリ分類機能

### 実装上の注意点
- HTTPサーバーはChi v5フレームワークを使用
- 現在のmain.goは基本的なウェルカムページのみ実装
- API実装時はOpenAPI仕様に準拠すること
- 日本語でのコメントとドキュメント作成を推奨