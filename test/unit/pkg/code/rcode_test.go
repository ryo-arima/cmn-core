package code_test

import (
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/code"
	"github.com/stretchr/testify/assert"
)

func TestRCodeDefinitions(t *testing.T) {
	tests := []struct {
		name  string
		rcode code.RCode
	}{
		{"UCPCU001", code.UCPCU001},
		{"UCPCU002", code.UCPCU002},
		{"UCPCU003", code.UCPCU003},
		{"UCPCU004", code.UCPCU004},
		{"UCPCUSUC", code.UCPCUSUC},
		{"UCPGU001", code.UCPGU001},
		{"UCPGUSUC", code.UCPGUSUC},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.rcode.Code, "Code should not be empty")
			assert.NotEmpty(t, tt.rcode.Message, "Message should not be empty")
		})
	}
}

func TestRCodeStructure(t *testing.T) {
	rcode := code.RCode{
		Code:    "TEST_CODE",
		Message: "Test message",
	}

	assert.Equal(t, "TEST_CODE", rcode.Code)
	assert.Equal(t, "Test message", rcode.Message)
}
