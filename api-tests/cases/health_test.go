package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Đợi server khởi động
	time.Sleep(2 * time.Second)

	t.Run("🏥 Kiểm tra Health Check API", func(t *testing.T) {
		// Thực hiện request
		resp, err := http.Get("http://localhost:8080/api/v1/system/health")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi health check API: %v", err)
		}
		defer resp.Body.Close()

		// Kiểm tra status code
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code phải là 200")

		// Parse response
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err, "Phải parse được JSON response")

		// Kiểm tra kết quả
		assert.Equal(t, "healthy", result["status"], "Status phải là 'healthy'")
		assert.NotNil(t, result["timestamp"], "Phải có trường timestamp")

		// In kết quả test
		fmt.Printf("✅ Health Check thành công:\n")
		fmt.Printf("   - Status: %v\n", result["status"])
		fmt.Printf("   - Timestamp: %v\n", result["timestamp"])
	})
}
