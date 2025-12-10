#!/bin/bash
# Script lấy Firebase ID token cho test
# Sử dụng: ./api-tests/scripts/get-firebase-token.sh -e "test@example.com" -p "Test@123" -k "YOUR_API_KEY"

EMAIL="test@example.com"
PASSWORD="Test@123"
API_KEY=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--email)
            EMAIL="$2"
            shift 2
            ;;
        -p|--password)
            PASSWORD="$2"
            shift 2
            ;;
        -k|--api-key)
            API_KEY="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Kiểm tra API Key
if [ -z "$API_KEY" ]; then
    # Thử lấy từ environment variable
    API_KEY=$FIREBASE_API_KEY
    if [ -z "$API_KEY" ]; then
        echo "[ERROR] Firebase API Key không được cung cấp"
        echo "Sử dụng: -k 'YOUR_API_KEY' hoặc set FIREBASE_API_KEY environment variable"
        exit 1
    fi
fi

echo "========================================"
echo "  LẤY FIREBASE ID TOKEN CHO TEST"
echo "========================================"
echo ""

# Tạo request body
BODY=$(cat <<EOF
{
    "email": "$EMAIL",
    "password": "$PASSWORD",
    "returnSecureToken": true
}
EOF
)

# Gọi Firebase REST API
URL="https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=$API_KEY"

echo "Đang đăng nhập với Firebase..."
RESPONSE=$(curl -s -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d "$BODY")

# Kiểm tra response
ID_TOKEN=$(echo "$RESPONSE" | grep -o '"idToken":"[^"]*' | cut -d'"' -f4)

if [ -z "$ID_TOKEN" ]; then
    echo "[ERROR] Không nhận được ID token từ response"
    echo "Response: $RESPONSE"
    exit 1
fi

echo "[OK] Đăng nhập thành công!"
echo ""
echo "Firebase ID Token:"
echo "$ID_TOKEN"
echo ""
echo "========================================"
echo "  SET ENVIRONMENT VARIABLE"
echo "========================================"
echo ""
echo "Bash:"
echo "export TEST_FIREBASE_ID_TOKEN=\"$ID_TOKEN\""
echo ""
echo "PowerShell:"
echo "\$env:TEST_FIREBASE_ID_TOKEN=\"$ID_TOKEN\""
echo ""

# Hỏi có muốn set tự động không
read -p "Bạn có muốn set token vào environment variable ngay bây giờ? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    export TEST_FIREBASE_ID_TOKEN="$ID_TOKEN"
    echo "[OK] Đã set TEST_FIREBASE_ID_TOKEN vào environment variable"
    echo "Lưu ý: Token chỉ có hiệu lực trong session hiện tại"
fi

