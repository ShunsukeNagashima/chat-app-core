# デプロイ用コンテナに含めるバイナリを作成するコンテナ
FROM golang:1.19.3-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app ./cmd

# --------------------------------------------------

# デプロイ用のコンテナ
FROM debian:bullseye-slim as deploy

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=deploy-builder /app/app .

CMD ["./app"]

# --------------------------------------------------
# ローカル開発環境で利用するホットリロード環境

FROM golang:1.19.3 as dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

CMD ["air"]
