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

		// Kiểm tra kết quả - Response format: {code, message, data: {status, timestamp, services}, status}
		assert.Equal(t, "success", result["status"], "Status phải là 'success'")
		
		data, ok := result["data"].(map[string]interface{})
		assert.True(t, ok, "Phải có trường data")
		if ok {
			assert.Equal(t, "healthy", data["status"], "Data status phải là 'healthy'")
			assert.NotNil(t, data["timestamp"], "Phải có trường timestamp trong data")
		}

		// In kết quả test
		fmt.Printf("✅ Health Check thành công:\n")
		fmt.Printf("   - Response Status: %v\n", result["status"])
		if ok {
			fmt.Printf("   - Health Status: %v\n", data["status"])
			fmt.Printf("   - Timestamp: %v\n", data["timestamp"])
			if services, ok := data["services"].(map[string]interface{}); ok {
				fmt.Printf("   - Services: %v\n", services)
			}
		}
	})
}
