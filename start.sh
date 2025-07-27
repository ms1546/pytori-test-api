docker compose down -v
GOARCH=arm64 GOOS=linux go build -o ./cmd/repo-summary/bootstrap ./cmd/repo-summary
docker compose up -d
source .env && go run ./scripts/setup.go
