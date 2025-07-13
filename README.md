# 📦 Shiritori API - Local Setup Guide

このプロジェクトは、Gitのコミット履歴を利用したしりとり遊びを管理・表示するAPIです。
DynamoDB Local、AWS SAM CLI を使用して、ローカルで完全再現できます。

---

## ✅ 前提条件（インストールされていること）

- Docker Desktop：https://www.docker.com/products/docker-desktop/
- Node.js：https://nodejs.org/
- AWS SAM CLI（Rosettaターミナルでインストール推奨）
  https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html

💡 Apple Silicon (M1/M2/M3) ユーザーへ
Rosetta ターミナルを使って `sam local start-api` を実行してください。

---

## 🚀 起動手順

### 1. このディレクトリに移動

cd pytori-test-api

---

### 2. Node.js 依存パッケージのインストール

npm i

---

### 3. DynamoDB Local の起動

docker compose up -d

---

### 4. テーブルとテストデータの作成

node utils/createTable.js && node utils/setup-test.js

✅ game_commits を作成しました
✅ game_repos を作成しました
✅ テストデータを投入しました

---

### 5. SAM API のローカル起動（Rosettaターミナル）

sam local start-api

起動成功時：

Mounting RepoSummaryFunction at http://127.0.0.1:3000/repo-summary [GET]

---

## 🧪 テストAPI呼び出し

### 単一リポジトリを取得：

curl "http://localhost:3000/repo-summary?repository_id=101"

### 全リポジトリを取得：

curl "http://localhost:3000/repo-summary"

---

### ✅ 単一リポジトリのレスポンス例：

{
  "repository_id": 101,
  "repository_name": "team-a",
  "status": 1,
  "shiritori_count": 1,
  "current_word": "ぬいぐるみ",
  "review_comment": "ちょーすごい",
  "merged_on": "2025-07-10T15:20:00Z"
}

---

### ✅ 全リポジトリのレスポンス例：

[
  {
    "repository_id": 101,
    "repository_name": "team-a",
    "status": 1,
    "shiritori_count": 1,
    "current_word": "ぬいぐるみ",
    "review_comment": "ちょーすごい",
    "merged_on": "2025-07-10T15:20:00Z"
  },
  ...
]

---

## 🧠 補足

- DynamoDB Local は http://localhost:8000 で動作
- SAM CLI は `host.docker.internal` を通じて DynamoDB Local にアクセス

---

## 🔗 参考リンク

- SAM CLI
  https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli.html
- DynamoDB Local
  https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html
- aws-sdk/client-dynamodb
  https://www.npmjs.com/package/@aws-sdk/client-dynamodb
