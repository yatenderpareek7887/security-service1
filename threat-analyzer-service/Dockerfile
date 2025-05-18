FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY src/ ./src/
WORKDIR /app/src
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o threat-analyzer-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/src/threat-analyzer-service .
EXPOSE 8081
CMD ["./threat-analyzer-service"]