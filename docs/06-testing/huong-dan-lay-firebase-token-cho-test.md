# HƯỚNG DẪN LẤY FIREBASE ID TOKEN CHO TEST

Tài liệu này hướng dẫn cách lấy Firebase ID token để chạy test.

---

## PHƯƠNG PHÁP 1: TỪ FIREBASE CONSOLE (Đơn giản nhất)

### Bước 1: Tạo User Test trong Firebase Console
1. Đăng nhập vào [Firebase Console](https://console.firebase.google.com/)
2. Chọn project của bạn
3. Vào **Authentication** > **Users**
4. Click **Add user**
5. Nhập email và password (hoặc tạo user bằng email/password)
6. Copy **UID** của user vừa tạo

### Bước 2: Lấy ID Token bằng Firebase Client SDK
Sử dụng script Node.js hoặc Python để lấy ID token:

**Node.js:**
```javascript
const admin = require('firebase-admin');
const firebase = require('firebase/app');
require('firebase/auth');

// Khởi tạo Firebase Admin SDK
const serviceAccount = require('./path/to/service-account.json');
admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
});

// Khởi tạo Firebase Client SDK
const firebaseConfig = {
  apiKey: "YOUR_API_KEY",
  authDomain: "YOUR_PROJECT_ID.firebaseapp.com",
  projectId: "YOUR_PROJECT_ID"
};
firebase.initializeApp(firebaseConfig);

// Đăng nhập với email/password
firebase.auth().signInWithEmailAndPassword('test@example.com', 'password123')
  .then((userCredential) => {
    // Lấy ID token
    return userCredential.user.getIdToken();
  })
  .then((idToken) => {
    console.log('Firebase ID Token:', idToken);
    // Copy token này và set vào environment variable
  })
  .catch((error) => {
    console.error('Error:', error);
  });
```

**Python:**
```python
import firebase_admin
from firebase_admin import credentials, auth
import pyrebase

# Khởi tạo Firebase Admin SDK
cred = credentials.Certificate("path/to/service-account.json")
firebase_admin.initialize_app(cred)

# Khởi tạo Firebase Client SDK
config = {
    "apiKey": "YOUR_API_KEY",
    "authDomain": "YOUR_PROJECT_ID.firebaseapp.com",
    "projectId": "YOUR_PROJECT_ID"
}
firebase = pyrebase.initialize_app(config)
auth_client = firebase.auth()

# Đăng nhập với email/password
user = auth_client.sign_in_with_email_and_password("test@example.com", "password123")
id_token = user['idToken']
print(f"Firebase ID Token: {id_token}")
```

---

## PHƯƠNG PHÁP 2: TỪ FIREBASE ADMIN SDK (Go)

### Bước 1: Tạo Custom Token
Sử dụng Firebase Admin SDK để tạo custom token:

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "firebase.google.com/go/v4"
    "firebase.google.com/go/v4/auth"
    "google.golang.org/api/option"
)

func main() {
    ctx := context.Background()
    
    // Khởi tạo Firebase Admin SDK
    opt := option.WithCredentialsFile("path/to/service-account.json")
    app, err := firebase.NewApp(ctx, &firebase.Config{
        ProjectID: "your-project-id",
    }, opt)
    if err != nil {
        panic(err)
    }
    
    authClient, err := app.Auth(ctx)
    if err != nil {
        panic(err)
    }
    
    // Tạo custom token cho test user
    // Lưu ý: UID phải tồn tại trong Firebase Authentication
    testUID := "test_user_uid_here"
    customToken, err := authClient.CustomToken(ctx, testUID)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Custom Token:", customToken)
    // Custom token cần được exchange thành ID token
}
```

### Bước 2: Exchange Custom Token thành ID Token
Sử dụng Firebase REST API để exchange:

```bash
curl -X POST \
  https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=YOUR_API_KEY \
  -H 'Content-Type: application/json' \
  -d '{
    "token": "CUSTOM_TOKEN_HERE",
    "returnSecureToken": true
  }'
```

Response sẽ có `idToken` trong response.

---

## PHƯƠNG PHÁP 3: SỬ DỤNG SCRIPT HELPER (Khuyến nghị)

Tạo script helper để tự động lấy token:

### Script PowerShell (`get-firebase-token.ps1`):
```powershell
# Script lấy Firebase ID token cho test
param(
    [string]$Email = "test@example.com",
    [string]$Password = "Test@123",
    [string]$ApiKey = "",
    [string]$ProjectId = ""
)

# Sử dụng Firebase REST API để đăng nhập và lấy ID token
$body = @{
    email = $Email
    password = $Password
    returnSecureToken = $true
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=$ApiKey" `
    -Method Post `
    -Body $body `
    -ContentType "application/json"

$idToken = $response.idToken
Write-Host "Firebase ID Token:" -ForegroundColor Green
Write-Host $idToken
Write-Host ""
Write-Host "Set environment variable:" -ForegroundColor Yellow
Write-Host "`$env:TEST_FIREBASE_ID_TOKEN=`"$idToken`""
```

**Sử dụng:**
```powershell
.\get-firebase-token.ps1 -Email "test@example.com" -Password "Test@123" -ApiKey "YOUR_API_KEY"
```

---

## PHƯƠNG PHÁP 4: TẠO USER VÀ LẤY TOKEN TỰ ĐỘNG

Tạo script Go để tự động tạo user và lấy token:

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

type SignInRequest struct {
    Email             string `json:"email"`
    Password          string `json:"password"`
    ReturnSecureToken bool   `json:"returnSecureToken"`
}

type SignInResponse struct {
    IDToken      string `json:"idToken"`
    Email        string `json:"email"`
    RefreshToken string `json:"refreshToken"`
    ExpiresIn    string `json:"expiresIn"`
    LocalId      string `json:"localId"`
}

func main() {
    apiKey := os.Getenv("FIREBASE_API_KEY")
    if apiKey == "" {
        fmt.Println("Error: FIREBASE_API_KEY environment variable not set")
        os.Exit(1)
    }

    email := "test@example.com"
    password := "Test@123"

    // Đăng nhập và lấy ID token
    reqBody := SignInRequest{
        Email:             email,
        Password:          password,
        ReturnSecureToken: true,
    }

    jsonData, _ := json.Marshal(reqBody)
    url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    var signInResp SignInResponse
    json.Unmarshal(body, &signInResp)

    if signInResp.IDToken == "" {
        fmt.Printf("Error: Failed to get ID token. Response: %s\n", string(body))
        os.Exit(1)
    }

    fmt.Printf("Firebase ID Token:\n%s\n\n", signInResp.IDToken)
    fmt.Printf("Set environment variable:\n")
    fmt.Printf("$env:TEST_FIREBASE_ID_TOKEN=\"%s\"\n", signInResp.IDToken)
}
```

---

## CÁCH SỬ DỤNG TOKEN

### PowerShell:
```powershell
# Set token từ output của script
$env:TEST_FIREBASE_ID_TOKEN="your_token_here"

# Chạy test
go test ./api-tests/cases -v
```

### Bash:
```bash
# Set token
export TEST_FIREBASE_ID_TOKEN="your_token_here"

# Chạy test
go test ./api-tests/cases -v
```

---

## LƯU Ý

1. **Token có thời gian hết hạn**: Firebase ID token thường hết hạn sau 1 giờ
2. **Refresh token**: Có thể sử dụng refresh token để lấy ID token mới
3. **Test user**: Nên tạo user riêng cho test, không dùng user production
4. **Bảo mật**: Không commit token vào git, chỉ dùng trong test

---

## KHUYẾN NGHỊ

**Cách đơn giản nhất**: Sử dụng Firebase REST API với email/password:
1. Tạo user test trong Firebase Console
2. Sử dụng REST API để đăng nhập và lấy ID token
3. Set token vào environment variable
4. Chạy test

**Script helper**: Tạo script tự động để không phải làm thủ công mỗi lần.

---

**Chúc bạn test thành công! ✅**

