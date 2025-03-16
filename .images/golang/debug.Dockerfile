# デバッグ用ビルドステージ
FROM golang:latest AS debug
WORKDIR /app
COPY . .
RUN go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080 2345
CMD ["dlv", "debug", "--headless", "--listen=:2345", "--api-version=2", "--accept-multiclient", "--continue"]

# 本番用ビルドステージ
FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 最終的な軽量イメージ
FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/main .
EXPOSE 8080
CMD ["./main"]

# 開発環境: docker build --target debug -t myapp:debug .
# 本番環境: docker build --target production -t myapp:prod .