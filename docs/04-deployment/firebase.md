# Firebase Setup

Hướng dẫn cài đặt và cấu hình Firebase cho hệ thống.

## 📋 Tổng Quan

Hệ thống sử dụng Firebase Authentication để xác thực người dùng. Tài liệu này hướng dẫn cách setup Firebase từ đầu.

## 🚀 Bước 1: Tạo Firebase Project

1. Truy cập https://console.firebase.google.com/
2. Click "Add project" hoặc chọn project có sẵn
3. Nhập tên project (ví dụ: `meta-commerce-auth`)
4. Chọn Google Analytics (tùy chọn)
5. Click "Create project"

## 🔐 Bước 2: Bật Authentication

1. Vào **Authentication** > **Get started**
2. Chọn tab **Sign-in method**
3. Bật các providers:
   - **Email/Password**: Bật và lưu
   - **Google**: Bật và cấu hình OAuth consent screen
   - **Facebook**: Bật và cấu hình App ID và App Secret
   - **Phone**: Bật (cần verify domain)

## 🔑 Bước 3: Tạo Service Account

1. Vào **Project Settings** > **Service Accounts**
2. Click **Generate new private key**
3. Lưu file JSON vào `api/config/firebase/service-account.json`

**Lưu ý:** Không commit file này vào git!

## 🔑 Bước 4: Lấy API Key

1. Vào **Project Settings** > **General**
2. Scroll xuống phần **Your apps**
3. Copy **Web API Key**

## ⚙️ Bước 5: Cấu Hình Environment

Thêm vào file `.env`:

```env
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_CREDENTIALS_PATH=config/firebase/service-account.json
FIREBASE_API_KEY=your-api-key
```

## ✅ Kiểm Tra

1. Khởi động server
2. Kiểm tra log xem Firebase đã được khởi tạo chưa
3. Test đăng nhập bằng Firebase

## 📚 Tài Liệu Liên Quan

- [Hướng Dẫn Đăng Ký Firebase](../huong-dan-dang-ky-firebase.md)
- [Hướng Dẫn Cài Đặt Firebase](../huong-dan-cai-dat-firebase.md)
- [Firebase Authentication với Database](../firebase-auth-voi-database.md)

