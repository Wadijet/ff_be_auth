# Debug Guide

Hướng dẫn debug và xử lý lỗi trong hệ thống.

## 📋 Tổng Quan

Tài liệu này hướng dẫn cách debug và xử lý các vấn đề trong hệ thống.

## 🔍 Debug Techniques

### 1. Xem Logs

```bash
# Xem log real-time
tail -f api/logs/app.log

# Windows PowerShell
Get-Content api/logs/app.log -Wait -Tail 50
```

### 2. Enable Debug Mode

Log level mặc định là Debug. Nếu không thấy log, kiểm tra:
```go
// cmd/server/main.go
logrus.SetLevel(logrus.DebugLevel)
```

### 3. Test Endpoints

```bash
# Health check
curl http://localhost:8080/api/v1/system/health

# Với authentication
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/auth/profile
```

### 4. Database Queries

```bash
# MongoDB shell
mongosh

# Kiểm tra collections
show collections

# Query data
db.users.find()
```

## 🐛 Common Issues

### Issue 1: Server Không Khởi Động

**Triệu chứng:** Server không start hoặc crash ngay lập tức

**Debug:**
1. Kiểm tra MongoDB có chạy không
2. Kiểm tra port có bị chiếm không
3. Xem log file để tìm lỗi

### Issue 2: Authentication Fail

**Triệu chứng:** Không thể đăng nhập

**Debug:**
1. Kiểm tra Firebase credentials
2. Verify Firebase token
3. Kiểm tra JWT secret

### Issue 3: Permission Denied

**Triệu chứng:** 403 Forbidden

**Debug:**
1. Kiểm tra user có role không
2. Kiểm tra role có permission không
3. Verify permission trong middleware

## 📚 Tài Liệu Liên Quan

- [Lỗi Thường Gặp](loi-thuong-gap.md)
- [Phân Tích Log](phan-tich-log.md)
- [Performance Issues](performance.md)

