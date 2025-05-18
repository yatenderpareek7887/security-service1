FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY src/ ./src/
WORKDIR /app/src
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o log-ingestor-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/src/log-ingestor-service .
EXPOSE 8080
CMD ["./log-ingestor-service"]