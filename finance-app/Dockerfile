FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose ./goose
COPY --from=builder /app/migrations ./migrations
COPY docker-entrypoint.sh .

RUN chmod +x docker-entrypoint.sh

EXPOSE 8080

CMD ["./docker-entrypoint.sh"]