# Start from the official Golang image
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o goRAGnarok ./cmd/main.go

# Use scratch for minimal final image
FROM scratch
WORKDIR /app
COPY --from=builder /app/goRAGnarok .
CMD ["/app/goRAGnarok"]
