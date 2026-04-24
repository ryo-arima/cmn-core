package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// GetTestDataPath returns the absolute path to the testdata directory
func GetTestDataPath() string {
	_, filename, _, _ := runtime.Caller(0)
	// Navigate from test/unit/internal/testutil to test/unit/testdata
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "..", "..", "testdata")
}

// LoadJSONFile loads a JSON file from testdata directory and unmarshals it into the provided interface
func LoadJSONFile(relativePath string, v interface{}) error {
	fullPath := filepath.Join(GetTestDataPath(), relativePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// LoadYAMLFile loads a YAML file from testdata directory
func LoadYAMLFile(relativePath string) ([]byte, error) {
	fullPath := filepath.Join(GetTestDataPath(), relativePath)
	return os.ReadFile(fullPath)
}

// LoadTextFile loads a text file from testdata directory
func LoadTextFile(relativePath string) (string, error) {
	fullPath := filepath.Join(GetTestDataPath(), relativePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetFilePath returns the full path to a file in testdata directory
func GetFilePath(relativePath string) string {
	return filepath.Join(GetTestDataPath(), relativePath)
}

// FileExists checks if a file exists in testdata directory
func FileExists(relativePath string) bool {
	fullPath := filepath.Join(GetTestDataPath(), relativePath)
	_, err := os.Stat(fullPath)
	return err == nil
}
