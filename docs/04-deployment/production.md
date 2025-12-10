# Triển Khai Production

Hướng dẫn triển khai hệ thống lên môi trường production.

## 📋 Tổng Quan

Tài liệu này hướng dẫn cách deploy hệ thống FolkForm Auth Backend lên môi trường production.

## 🔒 Bảo Mật Production

### Environment Variables

**KHÔNG BAO GIỜ** commit file `.env` chứa secret keys vào git!

**Sử dụng:**
- Environment variables của hệ điều hành
- Secret management service (AWS Secrets Manager, HashiCorp Vault, etc.)
- Docker secrets (nếu dùng Docker)

### JWT Secret

Sử dụng secret key mạnh (ít nhất 32 ký tự, ngẫu nhiên):

```bash
# Generate random secret
openssl rand -base64 32
```

### CORS Configuration

**KHÔNG** sử dụng `CORS_ORIGINS=*` trong production!

```env
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
CORS_ALLOW_CREDENTIALS=true
```

### HTTPS

Luôn sử dụng HTTPS trong production. Cấu hình reverse proxy (Nginx, Caddy) với SSL certificate.

## 🚀 Build Application

### Build Binary

```bash
cd api
go build -o server -ldflags="-s -w" cmd/server/main.go
```

**Flags:**
- `-s`: Strip symbol table
- `-w`: Omit DWARF symbol table

### Build cho Linux từ Windows

```bash
# Install cross-compiler
go install github.com/mitchellh/gox@latest

# Build
gox -osarch="linux/amd64" -output="server" ./cmd/server
```

## 🐳 Docker Deployment

### Dockerfile

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/config ./config
CMD ["./server"]
```

### Docker Compose

```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ADDRESS=8080
      - MONGODB_CONNECTION_URI=mongodb://mongo:27017
    depends_on:
      - mongo
    restart: unless-stopped

  mongo:
    image: mongo:latest
    volumes:
      - mongo-data:/data/db
    restart: unless-stopped

volumes:
  mongo-data:
```

## 🔄 Systemd Service

Xem [Systemd Service](systemd.md) để biết cách cấu hình systemd service.

## 📊 Monitoring

### Health Check

Endpoint: `GET /api/v1/system/health`

Sử dụng monitoring service (Prometheus, Grafana) để theo dõi:
- Health status
- Response time
- Error rate
- Database connection

### Logging

- Log level: `Info` hoặc `Warn` (không dùng `Debug` trong production)
- Log rotation: Sử dụng logrotate hoặc tương tự
- Centralized logging: Gửi log đến ELK stack hoặc tương tự

## 🔧 Performance Tuning

### MongoDB

- Connection pooling: Cấu hình max connections
- Indexes: Đảm bảo tất cả indexes đã được tạo
- Replica set: Sử dụng replica set cho high availability

### Application

- Rate limiting: Cấu hình phù hợp
- Caching: Sử dụng cache cho permissions
- Connection timeout: Cấu hình timeout phù hợp

## 📚 Tài Liệu Liên Quan

- [Cấu Hình Server](cau-hinh-server.md)
- [MongoDB Setup](mongodb.md)
- [Firebase Setup](firebase.md)
- [Systemd Service](systemd.md)

