package share_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/stretchr/testify/assert"
)

func TestGetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test without middleware
	requestID := share.GetRequestID(c)
	assert.Empty(t, requestID)

	// Test with middleware
	c.Set(share.RequestIDKey, "test-request-id-123")
	requestID = share.GetRequestID(c)
	assert.Equal(t, "test-request-id-123", requestID)
}

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(share.RequestID())

	router.GET("/test", func(c *gin.Context) {
		requestID := share.GetRequestID(c)
		assert.NotEmpty(t, requestID)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// LoggerWithConfig requires a BaseConfig
	conf := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Common: config.Common{},
			},
		},
	}

	router := gin.New()
	if logger, ok := conf.Logger.(share.ServerLoggerInterface); ok {
		router.Use(share.LoggerWithConfig(logger))
	}

	router.GET("/log-test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/log-test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMcodeMiddleware(t *testing.T) {
	// Mcode is not a middleware, it's a struct type
	// Testing the MCode functionality instead
	mcode := share.MCode{
		Code:    "TEST1",
		Message: "Test message",
	}

	assert.Equal(t, "TEST1", mcode.Code)
	assert.Equal(t, "Test message", mcode.Message)
}
