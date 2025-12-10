# Phân Tích Log

Hướng dẫn phân tích log để tìm và xử lý vấn đề.

## 📋 Tổng Quan

Log được ghi vào `api/logs/app.log` với format text.

## 📝 Log Format

```
[LEVEL] [TIMESTAMP] [FUNCTION] [FILE:LINE] MESSAGE
```

**Ví dụ:**
```
[INFO] [2024-01-01 10:00:00.000] [main] [main.go:100] Server starting on port: 8080
[ERROR] [2024-01-01 10:00:01.000] [LoginWithFirebase] [service.auth.user.go:50] Invalid Firebase token
```

## 🔍 Log Levels

- **DEBUG**: Chi tiết debug info
- **INFO**: Thông tin chung
- **WARN**: Cảnh báo
- **ERROR**: Lỗi

## 📊 Phân Tích Log

### Tìm Lỗi

```bash
# Tìm tất cả ERROR
grep "ERROR" api/logs/app.log

# Tìm lỗi trong 1 giờ qua
grep "ERROR" api/logs/app.log | grep "$(date +%Y-%m-%d)"
```

### Tìm Request Cụ Thể

```bash
# Tìm request theo user ID
grep "user-id-here" api/logs/app.log
```

### Thống Kê

```bash
# Đếm số ERROR
grep -c "ERROR" api/logs/app.log

# Đếm số request
grep -c "Request" api/logs/app.log
```

## 📚 Tài Liệu Liên Quan

- [Debug Guide](debug.md)
- [Lỗi Thường Gặp](loi-thuong-gap.md)

