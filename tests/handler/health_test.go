package handler_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Test root health endpoint
	t.Run("Root health check", func(t *testing.T) {
		maxRetries := 5
		var lastErr error

		for i := 0; i < maxRetries; i++ {
			resp, err := http.Get("http://localhost:3000/health")
			if err == nil {
				defer resp.Body.Close()
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				return
			}
			lastErr = err
			t.Logf("Waiting for server to start... (attempt %d/%d)", i+1, maxRetries)
			time.Sleep(2 * time.Second)
		}
		t.Fatalf("Server did not start in time. Last error: %v", lastErr)
	})

	// Test API health endpoint
	t.Run("API health check", func(t *testing.T) {
		resp, err := http.Get("http://localhost:3000/api/health")
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
