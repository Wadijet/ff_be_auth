# Lỗi Thường Gặp

Tài liệu về các lỗi thường gặp và cách xử lý.

## 🔧 Lỗi Server

### Lỗi: Cannot Connect to MongoDB

**Triệu chứng:**
```
Error: cannot connect to MongoDB: connection refused
```

**Nguyên nhân:**
- MongoDB chưa được khởi động
- Connection string sai
- Firewall chặn kết nối

**Giải pháp:**
1. Kiểm tra MongoDB có đang chạy:
```bash
# Windows
Get-Service MongoDB

# Linux
sudo systemctl status mongod
```

2. Khởi động MongoDB:
```bash
# Windows
net start MongoDB

# Linux
sudo systemctl start mongod
```

3. Kiểm tra connection string trong `.env`:
```env
MONGODB_CONNECTION_URI=mongodb://localhost:27017
```

4. Kiểm tra firewall:
```bash
# Cho phép port 27017
```

### Lỗi: Port Already in Use

**Triệu chứng:**
```
Error: bind: address already in use :8080
```

**Nguyên nhân:**
- Port 8080 đã được sử dụng bởi process khác
- Server đã chạy từ trước

**Giải pháp:**
1. Tìm process đang sử dụng port:
```powershell
# Windows
netstat -ano | findstr :8080

# Linux
lsof -i :8080
```

2. Dừng process:
```powershell
# Windows
taskkill /PID <process-id> /F

# Linux
kill -9 <process-id>
```

3. Hoặc thay đổi port trong `.env`:
```env
ADDRESS=8081
```

### Lỗi: Firebase Initialization Failed

**Triệu chứng:**
```
Error: Firebase initialization failed: open config/firebase/service-account.json: no such file or directory
```

**Nguyên nhân:**
- File service account không tồn tại
- Đường dẫn sai
- Quyền truy cập file

**Giải pháp:**
1. Kiểm tra file có tồn tại:
```bash
ls api/config/firebase/service-account.json
```

2. Kiểm tra đường dẫn trong `.env`:
```env
FIREBASE_CREDENTIALS_PATH=config/firebase/service-account.json
```

3. Kiểm tra quyền truy cập file

4. Tải lại service account key từ Firebase Console

## 🔐 Lỗi Authentication

### Lỗi: Invalid Firebase Token

**Triệu chứng:**
```
Error: Invalid Firebase token
```

**Nguyên nhân:**
- Token đã hết hạn
- Token không hợp lệ
- Firebase project ID sai

**Giải pháp:**
1. Kiểm tra token có còn hạn:
   - Firebase ID token có thời hạn 1 giờ
   - Cần refresh token trước khi gửi

2. Kiểm tra Firebase project ID:
```env
FIREBASE_PROJECT_ID=your-project-id
```

3. Verify token với Firebase:
```javascript
// Frontend
const idToken = await user.getIdToken(true); // Force refresh
```

### Lỗi: Unauthorized

**Triệu chứng:**
```
HTTP 401 Unauthorized
```

**Nguyên nhân:**
- Không có token trong header
- Token không hợp lệ
- Token đã hết hạn

**Giải pháp:**
1. Kiểm tra header Authorization:
```javascript
headers: {
  'Authorization': `Bearer ${token}`
}
```

2. Kiểm tra token có còn hạn:
   - JWT token có thời hạn 24 giờ
   - Cần đăng nhập lại nếu hết hạn

3. Verify token:
```go
// Backend sẽ tự động verify token trong middleware
```

### Lỗi: Forbidden

**Triệu chứng:**
```
HTTP 403 Forbidden
```

**Nguyên nhân:**
- User không có quyền truy cập
- Role không được gán permission

**Giải pháp:**
1. Kiểm tra user có role phù hợp:
```http
GET /api/v1/auth/roles
```

2. Kiểm tra role có permission:
```http
GET /api/v1/role/:id/permissions
```

3. Gán role hoặc permission cho user

## 🗄️ Lỗi Database

### Lỗi: Duplicate Key Error

**Triệu chứng:**
```
Error: E11000 duplicate key error collection: users index: firebaseUid_1 dup key
```

**Nguyên nhân:**
- User với firebaseUid đã tồn tại
- Unique index bị vi phạm

**Giải pháp:**
1. Kiểm tra user đã tồn tại:
```javascript
// Tìm user theo firebaseUid
```

2. Sử dụng upsert thay vì create:
```go
// Tự động xử lý trong service
```

### Lỗi: Index Not Found

**Triệu chứng:**
```
Error: index not found
```

**Nguyên nhân:**
- Index chưa được tạo
- Database chưa được khởi tạo

**Giải pháp:**
1. Khởi tạo database:
```go
// Chạy init script
```

2. Tạo index thủ công:
```javascript
// MongoDB shell
db.users.createIndex({ firebaseUid: 1 }, { unique: true })
```

## 📝 Lỗi Validation

### Lỗi: Invalid Input

**Triệu chứng:**
```
Error: validation failed
```

**Nguyên nhân:**
- Input không đúng format
- Thiếu required fields
- Value không hợp lệ

**Giải pháp:**
1. Kiểm tra request body:
```json
{
  "idToken": "string",
  "hwid": "string"
}
```

2. Kiểm tra validation rules trong DTO

3. Xem error message chi tiết trong response

## 🔍 Debug Tips

### Xem Log

```bash
# Xem log real-time
tail -f api/logs/app.log

# Windows PowerShell
Get-Content api/logs/app.log -Wait -Tail 50
```

### Enable Debug Mode

Log level mặc định là Debug. Nếu không thấy log, kiểm tra:
```go
// cmd/server/main.go
logrus.SetLevel(logrus.DebugLevel)
```

### Test Endpoint

```bash
# Health check
curl http://localhost:8080/api/v1/system/health

# Với authentication
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/auth/profile
```

## 📚 Tài Liệu Liên Quan

- [Debug Guide](debug.md)
- [Phân Tích Log](phan-tich-log.md)
- [Performance Issues](performance.md)
- [Cài Đặt và Cấu Hình](../01-getting-started/cai-dat.md)

