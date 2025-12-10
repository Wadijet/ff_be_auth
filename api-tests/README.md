# API TESTS

Test suite cho API backend sử dụng Firebase Authentication.

---

## CÁCH CHẠY TEST

### 1. Lấy Firebase ID Token

**Cách đơn giản nhất - Sử dụng script:**

```powershell
# PowerShell
.\api-tests\scripts\get-firebase-token.ps1 -Email "test@example.com" -Password "Test@123" -ApiKey "YOUR_FIREBASE_API_KEY"
```

```bash
# Bash
chmod +x api-tests/scripts/get-firebase-token.sh
./api-tests/scripts/get-firebase-token.sh -e "test@example.com" -p "Test@123" -k "YOUR_FIREBASE_API_KEY"
```

**Lưu ý:** 
- Cần tạo user test trong Firebase Console trước
- Firebase API Key có thể lấy từ `api/config/env/development.env` (FIREBASE_API_KEY)

### 2. Set Environment Variable

**PowerShell:**
```powershell
$env:TEST_FIREBASE_ID_TOKEN="your_firebase_id_token_here"
```

**Bash:**
```bash
export TEST_FIREBASE_ID_TOKEN="your_firebase_id_token_here"
```

### 3. Chạy Test

```bash
# Chạy tất cả test
go test ./api-tests/cases -v

# Chạy test cụ thể
go test ./api-tests/cases -v -run TestAuthFlow
go test ./api-tests/cases -v -run TestErrorHandling
```

---

## TẠO USER TEST TRONG FIREBASE

1. Đăng nhập vào [Firebase Console](https://console.firebase.google.com/)
2. Chọn project của bạn
3. Vào **Authentication** > **Users**
4. Click **Add user**
5. Nhập email và password (ví dụ: `test@example.com` / `Test@123`)
6. Click **Add user**

---

## LẤY FIREBASE API KEY

Firebase API Key có trong file config:
- `api/config/env/development.env` - Tìm `FIREBASE_API_KEY`

Hoặc từ Firebase Console:
1. Vào **Project Settings** > **General**
2. Scroll xuống **Your apps**
3. Copy **Web API Key**

---

## TÀI LIỆU

Chi tiết hơn tại: `docs/huong-dan-lay-firebase-token-cho-test.md`

---

## LƯU Ý

- Firebase ID token có thời gian hết hạn (thường 1 giờ)
- Nếu token hết hạn, chạy lại script để lấy token mới
- Không commit token vào git
- Chỉ dùng user test, không dùng user production

---

**Chúc bạn test thành công! ✅**
