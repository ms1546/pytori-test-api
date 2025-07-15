# 📦 Shiritori API - Local Setup Guide

このプロジェクトは、Gitのコミット履歴を利用したしりとり遊びを管理・表示するAPIです。
DynamoDB Local、AWS SAM CLI を使用して、ローカルで完全再現できます。

---

## ✅ 前提条件（インストールされていること）

- Docker Desktop：https://www.docker.com/products/docker-desktop/
- **Go**（1.21+） <https://go.dev/dl/>
- AWS SAM CLI（※Rosettaターミナルでインストールする）
  https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html

💡 Apple Silicon (M1/M2/M3) ユーザーへ
Rosetta ターミナルを使って `sam local start-api` を実行、APIサーバを起動してください
参考:
https://qiita.com/funatsufumiya/items/cec08f1ba3387edc2eed
---

## 🚀 起動手順

### 1. このディレクトリに移動

cd <pytori-test-api>

---

### 2. Go バイナリをビルド
```bash
GOARCH=arm64 GOOS=linux go build -o ./cmd/repo-summary/bootstrap ./cmd/repo-summary
```

上記により、Go ランタイム用 Lambda バイナリ (bootstrap) を作成します。
---

### 3. DynamoDB Local の起動
```bash
docker compose up -d
```

（落とす時 docker compose down -v )

---

### 4. テーブルとテストデータの作成
```bash
source .env && go run ./scripts/setup.go
```

✅ pytori_commits を作成しました

✅ pytori_repos を作成しました

✅ テストデータを投入しました

---
### 5. AWSにログインしてDockerからpublic ECRのイメージ取得
```bash
aws ecr-public get-login-password --region us-east-1 \
  | docker login --username AWS --password-stdin public.ecr.aws
```

---
### 6. SAM API のローカル起動（※Rosettaターミナルで実施する）
```bash
(cd pytori-test-api)
DOCKER_HOST=unix:///Users/$USER/.docker/run/docker.sock sam build   # 変更時は毎回
DOCKER_HOST=unix:///Users/$USER/.docker/run/docker.sock sam local start-api
```

起動成功時：

Mounting RepoSummaryFunction at http://127.0.0.1:3000/repo-summary [GET]

---

## 🧪 テストAPI呼び出し
※  Docker コンテナから ホストマシンのネットワークにアクセスする際の名前解決のため、初回のみ時間がかかる
### 単一リポジトリを取得：
```bash
curl "http://localhost:3000/repo-summary?repository_id=101" | jq .
```

### 全リポジトリを取得：
```bash
curl "http://localhost:3000/repo-summary" | jq .
```

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
