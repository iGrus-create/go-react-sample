# Go + React Sample Project

このプロジェクトは Go バックエンドと React フロントエンドのサンプルアプリケーションです。

## 構成

-   `practice/`: Go バックエンド（API + Lambda）
-   `react-todo/`: React フロントエンド

## 環境構築

### 1. データベース環境の起動

```bash
# PostgreSQL + pgwebを起動
docker-compose up -d

# 確認
docker-compose ps
```

### 2. バックエンド（Go）

```bash
cd practice

# 環境変数設定
cp .env.example .env

# 依存関係インストール
go mod download

# アプリケーション起動
GO_ENV=dev go run .
```

### 3. フロントエンド（React）

```bash
cd react-todo

# 環境変数設定
cp .env.example .env

# 依存関係インストール
npm install

# 開発サーバー起動
npm run dev
```

## アクセス情報

-   **フロントエンド**: http://localhost:5173
-   **バックエンド API**: http://localhost:8080
-   **pgweb（DB 管理）**: http://localhost:8081
-   **PostgreSQL**: localhost:5434

## データベース接続情報

-   Host: localhost
-   Port: 5434
-   Database: testdb
-   User: root
-   Password: root
