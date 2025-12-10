package utils

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// GetFirebaseIDTokenForTest tạo Firebase ID token cho test
// Sử dụng Firebase Admin SDK để tạo custom token, sau đó exchange thành ID token
func GetFirebaseIDTokenForTest(projectID, credentialsPath, testUID string) (string, error) {
	// Khởi tạo Firebase Admin SDK
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: projectID,
	}, opt)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Firebase app: %v", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get Firebase Auth client: %v", err)
	}

	// Tạo custom token cho test user
	// Lưu ý: testUID phải là UID hợp lệ trong Firebase Authentication
	customToken, err := authClient.CustomToken(context.Background(), testUID)
	if err != nil {
		return "", fmt.Errorf("failed to create custom token: %v", err)
	}

	// Custom token cần được exchange thành ID token bằng Firebase REST API
	// Hoặc sử dụng Firebase Client SDK để sign in với custom token
	// Tạm thời trả về custom token, cần exchange thành ID token
	return customToken, nil
}

// GetFirebaseIDTokenFromEnv lấy Firebase ID token từ environment variable
// Hoặc tạo mới nếu có config Firebase
func GetFirebaseIDTokenFromEnv() (string, error) {
	// Kiểm tra xem đã có token trong env chưa
	token := os.Getenv("TEST_FIREBASE_ID_TOKEN")
	if token != "" {
		return token, nil
	}

	// Nếu không có, thử tạo từ Firebase config
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	testUID := os.Getenv("TEST_FIREBASE_UID")

	if projectID == "" || credentialsPath == "" || testUID == "" {
		return "", fmt.Errorf("missing Firebase config: need FIREBASE_PROJECT_ID, FIREBASE_CREDENTIALS_PATH, and TEST_FIREBASE_UID")
	}

	// Tạo custom token
	customToken, err := GetFirebaseIDTokenForTest(projectID, credentialsPath, testUID)
	if err != nil {
		return "", err
	}

	// TODO: Exchange custom token thành ID token
	// Cần sử dụng Firebase REST API hoặc Client SDK
	return customToken, nil
}
