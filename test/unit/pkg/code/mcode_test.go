package code_test

import (
	"strings"
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/code"
	"github.com/stretchr/testify/assert"
)

func TestMCode_PaddedCode(t *testing.T) {
	tests := []struct {
		name     string
		mcode    code.MCode
		expected func(result string) bool
	}{
		{
			name:  "Short code should be padded",
			mcode: code.MCode{Code: "TEST", Message: "Test message"},
			expected: func(result string) bool {
				return strings.HasPrefix(result, "TEST") && len(result) >= 4
			},
		},
		{
			name:  "Code at max length should not be padded",
			mcode: code.MCode{Code: strings.Repeat("A", code.GetMaxCodeLength()), Message: "Test"},
			expected: func(result string) bool {
				return len(result) == code.GetMaxCodeLength()
			},
		},
		{
			name:  "Code longer than max should return as-is",
			mcode: code.MCode{Code: strings.Repeat("B", code.GetMaxCodeLength()+5), Message: "Test"},
			expected: func(result string) bool {
				return len(result) == code.GetMaxCodeLength()+5
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.mcode.PaddedCode()
			assert.True(t, tt.expected(result), "PaddedCode result does not match expected pattern")
		})
	}
}

func TestGetMaxCodeLength(t *testing.T) {
	maxLen := code.GetMaxCodeLength()
	assert.Greater(t, maxLen, 0, "Max code length should be greater than 0")
	assert.LessOrEqual(t, maxLen, 100, "Max code length should be reasonable")
}

func TestMCodeDefinitions(t *testing.T) {
	tests := []struct {
		name  string
		mcode code.MCode
	}{
		{"SM1", code.SM1},
		{"SM2", code.SM2},
		{"SM3", code.SM3},
		{"MLWC1", code.MLWC1},
		{"MLWC2", code.MLWC2},
		{"MLWC3", code.MLWC3},
		{"MLWC4", code.MLWC4},
		{"MLWC5", code.MLWC5},
		{"CNDBC1", code.CNDBC1},
		{"CNDBC2", code.CNDBC2},
		{"CNDBC3", code.CNDBC3},
		{"RCHK1", code.RCHK1},
		{"RURP1", code.RURP1},
		{"RUCR1", code.RUCR1},
		{"RUUP1", code.RUUP1},
		{"RUDL1", code.RUDL1},
		{"RULS1", code.RULS1},
		{"RUCT1", code.RUCT1},
		{"UUGU1", code.UUGU1},
		{"UUCR1", code.UUCR1},
		{"UUCR2", code.UUCR2},
		{"UUUP1", code.UUUP1},
		{"UUUP2", code.UUUP2},
		{"UUDL1", code.UUDL1},
		{"UULS1", code.UULS1},
		{"UUCT1", code.UUCT1},
		{"GINLOG", code.GINLOG},
		{"UCPCU0", code.UCPCU0},
		{"UCPCU1", code.UCPCU1},
		{"UCPCU2", code.UCPCU2},
		{"UCPCU3", code.UCPCU3},
		{"UCPCU4", code.UCPCU4},
		{"UCPCU5", code.UCPCU5},
		{"UCPCU6", code.UCPCU6},
		{"UCPGU0", code.UCPGU0},
		{"UCPGU1", code.UCPGU1},
		{"UCPGU2", code.UCPGU2},
		{"UCPGU3", code.UCPGU3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.mcode.Code, "Code should not be empty")
			assert.NotEmpty(t, tt.mcode.Message, "Message should not be empty")
		})
	}
}
