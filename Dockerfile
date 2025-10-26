FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/ordersystem

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/internal/infra/database/migrations ./internal/infra/database/migrations

EXPOSE 8000 50051 8080

CMD ["./main"]