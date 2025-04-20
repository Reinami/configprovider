package sources

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTmpProperties(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.properties")

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write temp properties file: %v", err)
	}

	return filePath
}

func TestNewPropertiesFileSource_FileValid(t *testing.T) {
	content := `
# This is a comment
DEBUG=true
PORT=8080
NAME=TestApp

; another comment
SECRET = super-secret
`

	path := writeTmpProperties(t, content)

	source, err := NewPropertiesFileSource(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	tests := map[string]string{
		"DEBUG":  "true",
		"PORT":   "8080",
		"NAME":   "TestApp",
		"SECRET": "super-secret",
	}

	for key, expected := range tests {
		got, ok := source.Get(key)
		if !ok {
			t.Errorf("expected key %q to exist", key)
			continue
		}
		if got != expected {
			t.Errorf("key %q: expected %q, got %q", key, expected, got)
		}
	}
}

func TestNewPropertiesFileSource_MalformedLine(t *testing.T) {
	content := `
GOOD=okay
BAD_LINE_NO_EQUALS
ANOTHER=entry
`
	path := writeTmpProperties(t, content)

	_, err := NewPropertiesFileSource(path)
	if err == nil {
		t.Fatalf("expected error for malformed line, got none")
	}
}

func TestNewPropertiesFileSource_NotFound(t *testing.T) {
	_, err := NewPropertiesFileSource("nonexistent/path.properties")
	if err == nil {
		t.Fatalf("expected file not found error, got none")
	}
}
