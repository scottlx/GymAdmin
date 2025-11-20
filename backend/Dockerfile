FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o gym-admin ./cmd/server

FROM alpine:latest

WORKDIR /app

# Copy binary and config
COPY --from=builder /app/gym-admin .
COPY --from=builder /app/internal/config/config.yaml ./internal/config/

# Create logs directory
RUN mkdir -p logs

EXPOSE 8080

CMD ["./gym-admin"]
