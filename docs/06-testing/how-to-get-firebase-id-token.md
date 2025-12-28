# Cách Lấy Firebase ID Token Để Test

## 📋 Tổng Quan

Firebase ID Token là token được tạo bởi Firebase Authentication khi user đăng nhập. Token này được dùng để xác thực với backend API.

---

## 🔧 Cách 1: Lấy từ Firebase Console (Test Token)

### Bước 1: Truy cập Firebase Console
1. Vào [Firebase Console](https://console.firebase.google.com/)
2. Chọn project của bạn
3. Vào **Authentication** > **Users**

### Bước 2: Tạo Test User (nếu chưa có)
1. Click **Add user**
2. Nhập email và password (hoặc dùng các phương thức khác)
3. Lưu lại UID của user

### Bước 3: Lấy ID Token từ Firebase Admin SDK

**Lưu ý**: Firebase Console không cung cấp ID token trực tiếp. Bạn cần dùng Firebase Admin SDK hoặc Firebase Client SDK.

---

## 🔧 Cách 2: Lấy từ Web App (Khuyến nghị)

### Tạo file HTML đơn giản để lấy token

Tạo file `get-firebase-token.html`:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Get Firebase ID Token</title>
    <script src="https://www.gstatic.com/firebasejs/10.7.1/firebase-app-compat.js"></script>
    <script src="https://www.gstatic.com/firebasejs/10.7.1/firebase-auth-compat.js"></script>
</head>
<body>
    <h1>Get Firebase ID Token</h1>
    <input type="email" id="email" placeholder="Email" />
    <input type="password" id="password" placeholder="Password" />
    <button onclick="login()">Login & Get Token</button>
    <br><br>
    <textarea id="token" rows="10" cols="80" readonly></textarea>
    <button onclick="copyToken()">Copy Token</button>

    <script>
        // Thay bằng config của bạn
        const firebaseConfig = {
            apiKey: "YOUR_API_KEY",
            authDomain: "YOUR_PROJECT_ID.firebaseapp.com",
            projectId: "YOUR_PROJECT_ID",
            storageBucket: "YOUR_PROJECT_ID.appspot.com",
            messagingSenderId: "YOUR_MESSAGING_SENDER_ID",
            appId: "YOUR_APP_ID"
        };

        firebase.initializeApp(firebaseConfig);
        const auth = firebase.auth();

        async function login() {
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;

            try {
                const userCredential = await auth.signInWithEmailAndPassword(email, password);
                const user = userCredential.user;
                const idToken = await user.getIdToken();
                
                document.getElementById('token').value = idToken;
                console.log('ID Token:', idToken);
            } catch (error) {
                console.error('Error:', error);
                alert('Error: ' + error.message);
            }
        }

        function copyToken() {
            const token = document.getElementById('token').value;
            navigator.clipboard.writeText(token).then(() => {
                alert('Token copied to clipboard!');
            });
        }
    </script>
</body>
</html>
```

### Cách sử dụng:
1. Thay `firebaseConfig` bằng config của project bạn
2. Mở file HTML trong browser
3. Nhập email/password của user Firebase
4. Click "Login & Get Token"
5. Copy token từ textarea

---

## 🔧 Cách 3: Lấy từ Node.js Script

Tạo file `get-token.js`:

```javascript
const admin = require('firebase-admin');

// Khởi tạo Firebase Admin SDK
const serviceAccount = require('./path-to-your-service-account-key.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
});

// Lấy custom token (sau đó đổi sang ID token)
async function getCustomToken(uid) {
  try {
    const customToken = await admin.auth().createCustomToken(uid);
    console.log('Custom Token:', customToken);
    
    // Để lấy ID token, cần dùng client SDK với custom token này
    console.log('\nLưu ý: Cần dùng client SDK để đổi custom token sang ID token');
    console.log('Hoặc dùng cách 2 (Web App) để lấy ID token trực tiếp');
  } catch (error) {
    console.error('Error:', error);
  }
}

// Sử dụng
const uid = 'YOUR_FIREBASE_UID';
getCustomToken(uid);
```

---

## 🔧 Cách 4: Lấy từ Mobile App (React Native / Flutter)

### React Native:
```javascript
import auth from '@react-native-firebase/auth';

async function getIdToken() {
  try {
    const user = auth().currentUser;
    if (user) {
      const idToken = await user.getIdToken();
      console.log('ID Token:', idToken);
      return idToken;
    }
  } catch (error) {
    console.error('Error:', error);
  }
}
```

### Flutter:
```dart
import 'package:firebase_auth/firebase_auth.dart';

Future<String?> getIdToken() async {
  try {
    User? user = FirebaseAuth.instance.currentUser;
    if (user != null) {
      String idToken = await user.getIdToken();
      print('ID Token: $idToken');
      return idToken;
    }
  } catch (e) {
    print('Error: $e');
  }
  return null;
}
```

---

## 🔧 Cách 5: Dùng Firebase CLI (Nhanh nhất cho test)

### Cài đặt Firebase CLI:
```bash
npm install -g firebase-tools
```

### Login:
```bash
firebase login
```

### Lấy ID token từ test user:
```bash
# Tạo test user (nếu chưa có)
firebase auth:export users.json --project YOUR_PROJECT_ID

# Hoặc dùng Firebase Emulator để test
firebase emulators:start --only auth
```

---

## 🔧 Cách 6: Dùng Postman/Insomnia với Firebase REST API

### Bước 1: Lấy API Key từ Firebase Console
1. Vào Firebase Console > Project Settings > General
2. Copy **Web API Key**

### Bước 2: Gọi Firebase Auth REST API

**Request:**
```
POST https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=YOUR_API_KEY
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123",
  "returnSecureToken": true
}
```

**Response:**
```json
{
  "idToken": "eyJhbGciOiJSUzI1NiIsImtpZCI6Ij...",
  "email": "test@example.com",
  "refreshToken": "...",
  "expiresIn": "3600",
  "localId": "..."
}
```

Copy `idToken` từ response.

---

## 🚀 Cách Nhanh Nhất: Dùng Script Helper

Tạo file `scripts/get-firebase-token.js`:

```javascript
const admin = require('firebase-admin');
const readline = require('readline');

// Load service account
const serviceAccount = require('../path-to-service-account.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
});

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

rl.question('Enter Firebase UID: ', async (uid) => {
  try {
    // Tạo custom token
    const customToken = await admin.auth().createCustomToken(uid);
    console.log('\n✅ Custom Token created!');
    console.log('\n📋 Custom Token:');
    console.log(customToken);
    console.log('\n⚠️  Lưu ý: Cần dùng client SDK để đổi custom token sang ID token');
    console.log('   Hoặc dùng cách 2 (Web App) để lấy ID token trực tiếp');
  } catch (error) {
    console.error('❌ Error:', error.message);
  }
  rl.close();
});
```

---

## 📝 Thiết Lập Biến Môi Trường

Sau khi có Firebase ID Token:

### Windows PowerShell:
```powershell
$env:TEST_FIREBASE_ID_TOKEN = "eyJhbGciOiJSUzI1NiIsImtpZCI6Ij..."
```

### Windows CMD:
```cmd
set TEST_FIREBASE_ID_TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6Ij...
```

### Linux/Mac:
```bash
export TEST_FIREBASE_ID_TOKEN="eyJhbGciOiJSUzI1NiIsImtpZCI6Ij..."
```

### Hoặc tạo file `.env`:
```env
TEST_FIREBASE_ID_TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6Ij...
```

---

## ⚠️ Lưu Ý Quan Trọng

1. **Token có thời hạn**: Firebase ID Token có thời hạn (thường 1 giờ). Cần refresh token khi hết hạn.

2. **Bảo mật**: 
   - Không commit token vào Git
   - Không chia sẻ token công khai
   - Token chỉ dùng cho testing

3. **Refresh Token**: Nếu token hết hạn, có thể dùng refresh token để lấy token mới:
   ```javascript
   const idToken = await user.getIdToken(true); // Force refresh
   ```

4. **Test User**: Nên tạo user riêng cho testing, không dùng user production.

---

## 🎯 Khuyến Nghị

**Cách nhanh nhất cho testing:**
1. Dùng **Cách 2 (Web App)** - Tạo file HTML đơn giản
2. Hoặc dùng **Cách 6 (Postman/Insomnia)** với Firebase REST API
3. Lưu token vào biến môi trường
4. Chạy tests

---

## 📚 Tài Liệu Tham Khảo

- [Firebase Authentication](https://firebase.google.com/docs/auth)
- [Firebase REST API](https://firebase.google.com/docs/reference/rest/auth)
- [Firebase Admin SDK](https://firebase.google.com/docs/admin/setup)

