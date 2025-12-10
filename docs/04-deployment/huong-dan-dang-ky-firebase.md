# HƯỚNG DẪN ĐĂNG KÝ FIREBASE - CHỈ NHỮNG GÌ CẦN THIẾT

Tài liệu này hướng dẫn đơn giản cách đăng ký Firebase để sử dụng:
- Phone OTP Authentication
- Email Verification (tạm thời có thể bỏ qua hoặc dùng cách đơn giản)

---

## 1. ĐĂNG KÝ FIREBASE PROJECT

### Bước 1: Tạo Firebase Project

1. Truy cập [Firebase Console](https://console.firebase.google.com/)
2. Đăng nhập bằng tài khoản Google
3. Click **"Add project"** hoặc **"Create a project"**
4. Điền thông tin:
   - **Project name**: `meta-commerce-auth` (hoặc tên bạn muốn)
   - **Project ID**: Tự động tạo (có thể thay đổi)
   - **Google Analytics**: **Tắt** (không cần cho authentication)
5. Click **"Create project"**
6. Chờ vài giây để project được tạo
7. Click **"Continue"**

### Bước 2: Bật Phone Authentication

1. Trong Firebase Console, vào **"Authentication"** → Click **"Get started"**
2. Vào tab **"Sign-in method"**
3. Tìm **"Phone"** → Click vào
4. Click toggle để **"Enable"**
5. (Optional) Thêm số điện thoại test:
   - Click **"Phone numbers for testing"**
   - Thêm số điện thoại và verification code (ví dụ: `+84123456789` với code `123456`)
6. Click **"Save"**

✅ **Xong! Phone OTP đã sẵn sàng**

---

## 2. TẠO SERVICE ACCOUNT (CHO BACKEND)

### Bước 1: Tạo Service Account Key

1. Vào **"Project Settings"** (biểu tượng bánh răng ⚙️ ở góc trên bên trái)
2. Vào tab **"Service accounts"**
3. Click **"Generate new private key"**
4. Click **"Generate key"** trong popup cảnh báo
5. File JSON sẽ được download tự động (tên file: `your-project-firebase-adminsdk-xxxxx.json`)

### Bước 2: Lưu file Service Account

1. Tạo thư mục: `api/config/firebase/`
2. Đổi tên file thành: `service-account.json`
3. Di chuyển file vào: `api/config/firebase/service-account.json`

⚠️ **QUAN TRỌNG**: 
- **KHÔNG commit file này vào Git!**
- File này chứa private key, rất nhạy cảm
- Thêm vào `.gitignore`:
  ```
  api/config/firebase/service-account.json
  ```

### Bước 3: Lấy Project ID

1. Vẫn trong **"Project Settings"**
2. Tab **"General"**
3. Tìm **"Project ID"** → Copy và lưu lại

✅ **Xong! Service Account đã sẵn sàng**

---

## 3. LẤY WEB API KEY (CHO FRONTEND)

### Bước 1: Tạo Web App

1. Vẫn trong **"Project Settings"** → Tab **"General"**
2. Scroll xuống phần **"Your apps"**
3. Click **"Add app"** → Chọn **"Web"** (biểu tượng `</>`)
4. Điền **App nickname**: `Meta Commerce Web`
5. (Optional) Check **"Also set up Firebase Hosting"** nếu cần
6. Click **"Register app"**

### Bước 2: Lấy API Key

1. Sau khi register, bạn sẽ thấy config:
   ```javascript
   const firebaseConfig = {
     apiKey: "AIzaSy...",
     authDomain: "...",
     projectId: "...",
     // ...
   };
   ```
2. **Copy `apiKey`** và lưu lại (cần cho frontend)

✅ **Xong! Web API Key đã sẵn sàng**

---

## 4. CẤU HÌNH VÀO .ENV

Thêm vào file `api/config/env/development.env`:

```env
# Firebase Configuration
FIREBASE_PROJECT_ID=your-project-id-here
FIREBASE_CREDENTIALS_PATH=config/firebase/service-account.json
FIREBASE_API_KEY=your_api_key_here

# Frontend URL (cho redirect sau khi login)
FRONTEND_URL=http://localhost:3000
```

**Ví dụ:**
```env
FIREBASE_PROJECT_ID=meta-commerce-auth-12345
FIREBASE_CREDENTIALS_PATH=config/firebase/service-account.json
FIREBASE_API_KEY=AIzaSyAbc123xyz...
FRONTEND_URL=http://localhost:3000
```

---

## 5. EMAIL VERIFICATION (TÙY CHỌN)

Nếu cần email verification, có 2 cách:

### Cách 1: Tạm thời bỏ qua (Khuyến nghị cho bước đầu)
- Bỏ qua email verification
- Chỉ dùng Phone OTP và OAuth
- Có thể thêm sau

### Cách 2: Dùng Firebase Extensions (Khi cần)
- Cần đăng ký thêm SendGrid hoặc Mailgun
- Cài đặt Firebase Extension "Trigger Email"
- Xem hướng dẫn chi tiết trong file `huong-dan-dang-ky-dich-vu.md` (phần 4)

---

## 6. KIỂM TRA

### Kiểm tra Firebase đã setup đúng:

1. ✅ File service account tồn tại: `api/config/firebase/service-account.json`
2. ✅ Phone Authentication đã enable trong Firebase Console
3. ✅ Project ID đã copy vào `.env`
4. ✅ API Key đã copy vào `.env`
5. ✅ File `.env` đã được thêm vào `.gitignore`

### Test Phone OTP:

1. Frontend sử dụng Firebase SDK để gửi OTP
2. Backend verify ID token từ Firebase
3. Xem code implementation để test

---

## 7. TỔNG KẾT - CHECKLIST

- [ ] Firebase Project đã tạo
- [ ] Phone Authentication đã enable
- [ ] Service Account JSON đã download và lưu vào `api/config/firebase/service-account.json`
- [ ] Project ID đã copy
- [ ] Web API Key đã copy
- [ ] Đã thêm config vào `.env` file
- [ ] File service account đã thêm vào `.gitignore`

---

## 8. LƯU Ý BẢO MẬT

### ⚠️ QUAN TRỌNG:

1. **KHÔNG commit Service Account JSON vào Git:**
   ```
   # Thêm vào .gitignore
   api/config/firebase/service-account.json
   *.json
   ```

2. **KHÔNG commit .env file:**
   ```
   # Thêm vào .gitignore
   api/config/env/*.env
   ```

3. **Giữ bí mật:**
   - Service Account JSON = Private Key
   - API Key = Public (có thể dùng trong frontend)
   - Project ID = Public (có thể dùng trong frontend)

---

## 9. TROUBLESHOOTING

### Lỗi "Permission denied":
- Kiểm tra Service Account có đủ quyền
- Kiểm tra file service account JSON đúng

### Lỗi "Invalid credentials":
- Kiểm tra đường dẫn file service account đúng
- Kiểm tra file JSON không bị corrupt

### Phone OTP không gửi được:
- Kiểm tra Phone Authentication đã enable
- Kiểm tra Firebase project đúng
- Kiểm tra frontend config đúng

---

## 10. CHI PHÍ

### Firebase Authentication:
- **Phone OTP**: Miễn phí (có giới hạn)
- **Email/Password**: Miễn phí
- **OAuth Providers**: Miễn phí

### Free Tier:
- 50,000 MAU (Monthly Active Users) miễn phí
- Sau đó: $0.0055 per verification

### Khi nào cần trả phí:
- Khi có > 50,000 users/tháng
- Rất rẻ: ~$0.0055 per verification

---

## 11. TÀI LIỆU THAM KHẢO

- [Firebase Console](https://console.firebase.google.com/)
- [Firebase Authentication Documentation](https://firebase.google.com/docs/auth)
- [Firebase Phone Authentication](https://firebase.google.com/docs/auth/web/phone-auth)
- [Firebase Admin SDK for Go](https://firebase.google.com/docs/admin/setup)

---

**Chỉ cần Firebase là đủ để bắt đầu! 🚀**

**Các tính năng có thể dùng ngay:**
- ✅ Phone OTP Authentication
- ✅ Email/Password Authentication (nếu cần)
- ✅ Google OAuth (cần đăng ký thêm Google OAuth - xem file khác)
- ✅ Facebook OAuth (cần đăng ký thêm Facebook App - xem file khác)

**Email Verification có thể thêm sau khi cần!**

