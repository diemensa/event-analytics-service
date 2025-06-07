FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o analytics-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/analytics-service .

EXPOSE 8080

CMD ["./analytics-service"]
