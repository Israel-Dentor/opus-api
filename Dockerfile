# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 配置 Go 代理加速下载（使用国内镜像）
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GOSUMDB=off

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Create logs directory
RUN mkdir -p /app/logs

# Expose port
EXPOSE 3002

# Run server
CMD ["./server"]