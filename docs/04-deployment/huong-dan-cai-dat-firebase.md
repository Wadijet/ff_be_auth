# HƯỚNG DẪN CÀI ĐẶT FIREBASE DEPENDENCIES

## 1. CÀI ĐẶT GO DEPENDENCIES

Chạy các lệnh sau trong thư mục `api/`:

```bash
cd api
go get firebase.google.com/go/v4
go get google.golang.org/api/option
go mod tidy
```

Hoặc thêm vào `go.mod` và chạy `go mod tidy`:

```go
require (
    firebase.google.com/go/v4 v4.14.0
    google.golang.org/api v0.177.0
)
```

---

## 2. TẠO THƯ MỤC FIREBASE

```bash
mkdir -p api/config/firebase
```

---

## 3. THÊM SERVICE ACCOUNT FILE

1. Tải Service Account JSON từ Firebase Console
2. Lưu vào: `api/config/firebase/service-account.json`
3. Đảm bảo file này đã được thêm vào `.gitignore`

---

## 4. CẤU HÌNH ENVIRONMENT

Thêm vào `api/config/env/development.env`:

```env
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_CREDENTIALS_PATH=config/firebase/service-account.json
FIREBASE_API_KEY=your_api_key_here
FRONTEND_URL=http://localhost:3000
```

---

## 5. TEST

Sau khi cài đặt xong, chạy server và test endpoint:

```bash
POST http://localhost:8080/api/v1/auth/login/firebase
Content-Type: application/json

{
  "idToken": "firebase_id_token_from_frontend",
  "hwid": "test-device-hwid"
}
```

---

**Sau khi cài đặt dependencies, các lỗi linter sẽ biến mất! ✅**

