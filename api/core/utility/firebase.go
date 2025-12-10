package utility

import (
	"context"
	"fmt"
	"os"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	firebaseApp  *firebase.App
	firebaseAuth *auth.Client
)

// InitFirebase khởi tạo Firebase Admin SDK
func InitFirebase(projectID, credentialsPath string) error {
	// Kiểm tra file credentials tồn tại
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		return fmt.Errorf("firebase credentials file not found: %s", credentialsPath)
	}

	// Tạo Firebase app
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: projectID,
	}, opt)

	if err != nil {
		return fmt.Errorf("failed to initialize Firebase app: %v", err)
	}

	firebaseApp = app

	// Tạo Auth client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get Firebase Auth client: %v", err)
	}

	firebaseAuth = authClient
	return nil
}

// GetFirebaseAuth trả về Firebase Auth client
func GetFirebaseAuth() *auth.Client {
	return firebaseAuth
}

// VerifyIDToken verify Firebase ID token và trả về user info
func VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if firebaseAuth == nil {
		return nil, fmt.Errorf("firebase auth not initialized")
	}

	token, err := firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %v", err)
	}

	return token, nil
}

// GetUserByUID lấy thông tin user từ Firebase bằng UID
func GetUserByUID(ctx context.Context, uid string) (*auth.UserRecord, error) {
	if firebaseAuth == nil {
		return nil, fmt.Errorf("firebase auth not initialized")
	}

	user, err := firebaseAuth.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

