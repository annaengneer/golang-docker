# golang-docker

Go + Gin + MySQL を Docker で構築した Web アプリケーションです。  
Docker を使用することで、ローカル環境に依存せず簡単に起動できます。

## はじめに（重要）

このアプリは Docker を使用して起動します。  
事前に以下をインストールしてください。

- Docker
- Docker Compose

## 環境構築手順

以下の手順に従うことで、誰でも同じ環境でアプリを起動できます。

### ① リポジトリをクローン

ターミナルで以下を実行してください。

    git clone git@github.com:annaengneer/golang-docker.git
    cd golang-docker

### ② .env ファイルを作成

プロジェクト直下に `.env` ファイルを作成し、以下をコピーしてください。

    MYSQL_ROOT_PASSWORD=mysql
    MYSQL_DATABASE=mysql
    MYSQL_USER=mysql
    MYSQL_PASSWORD=mysql

    DB_USER=mysql
    DB_PASSWORD=mysql
    DB_HOST=mysql
    DB_PORT=3306
    DB_NAME=mysql

※ `.env` ファイルがないとアプリは起動しません

### ③ コンテナを起動

以下のコマンドを実行してください。

    docker compose up --build

### ④ 起動確認

以下の URL にアクセスしてください。

http://localhost:8080

「Hello World!」と表示されれば成功です。

## 動作確認

以下のエンドポイントで API を確認できます。

- GET / → Hello World が表示されます
- GET /albums → データ取得 API

## live reload（開発用）

このアプリは air を使用しており、  
ソースコードを変更すると自動で再起動されます。

## 使用技術

- Go
- Gin
- MySQL
- Docker / Docker Compose
- air（live reload）

## 工夫した点

- Docker のマルチステージビルドを使用し、開発環境と本番環境を分離しました
- Docker Compose を使用し、アプリとデータベースを一括で起動できるようにしました
- 初心者でも迷わず起動できるように、環境構築手順を詳細に記載しました
