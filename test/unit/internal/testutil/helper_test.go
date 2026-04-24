package testutil_test

import (
	"testing"

	"github.com/ryo-arima/cmn-core/test/unit/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTestDataPath(t *testing.T) {
	path := testutil.GetTestDataPath()
	assert.NotEmpty(t, path, "Test data path should not be empty")
	assert.Contains(t, path, "testdata", "Path should contain 'testdata'")
}

func TestLoadJSONFile(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		shouldError bool
	}{
		{
			name:        "Load valid user JSON",
			filePath:    "entity/user.json",
			shouldError: false,
		},
		{
			name:        "Load valid group JSON",
			filePath:    "entity/group.json",
			shouldError: false,
		},
		{
			name:        "Load non-existent file",
			filePath:    "entity/nonexistent.json",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data map[string]interface{}
			err := testutil.LoadJSONFile(tt.filePath, &data)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, data)
			}
		})
	}
}

func TestLoadYAMLFile(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		shouldError bool
	}{
		{
			name:        "Load valid config YAML",
			filePath:    "config/app.yaml",
			shouldError: false,
		},
		{
			name:        "Load minimal config YAML",
			filePath:    "config/app_minimal.yaml",
			shouldError: false,
		},
		{
			name:        "Load non-existent file",
			filePath:    "config/nonexistent.yaml",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := testutil.LoadYAMLFile(tt.filePath)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, data)
			}
		})
	}
}

func TestLoadTextFile(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		shouldError bool
	}{
		{
			name:        "Load invalid YAML as text",
			filePath:    "config/app_invalid.yaml",
			shouldError: false,
		},
		{
			name:        "Load non-existent file",
			filePath:    "nonexistent.txt",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := testutil.LoadTextFile(tt.filePath)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, data)
			}
		})
	}
}

func TestGetFilePath(t *testing.T) {
	path := testutil.GetFilePath("config/app.yaml")
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "testdata")
	assert.Contains(t, path, "config/app.yaml")
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		exists   bool
	}{
		{
			name:     "Existing file",
			filePath: "config/app.yaml",
			exists:   true,
		},
		{
			name:     "Non-existent file",
			filePath: "nonexistent.txt",
			exists:   false,
		},
		{
			name:     "Existing JSON file",
			filePath: "entity/user.json",
			exists:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := testutil.FileExists(tt.filePath)
			assert.Equal(t, tt.exists, exists)
		})
	}
}
