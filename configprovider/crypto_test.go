package configprovider

import (
	"errors"
	"testing"
)

// MockDecrypter for testing
type MockDecrypter struct {
	Fail   bool
	Output string
}

func (m *MockDecrypter) Decrypt(cipherText string) (string, error) {
	if m.Fail {
		return "", errors.New("mock decryption error")
	}
	return m.Output, nil
}

func TestDecryptValue_NoDecrypter(t *testing.T) {
	_, err := decryptValue("SECRET", "ENC:xyz", nil)
	if err == nil || err.Error() != "no decrypter is provided" {
		t.Errorf("expected 'no decrypter is provided', got: %v", err)
	}
}

func TestDecryptValue_DecryptionFails(t *testing.T) {
	mock := &MockDecrypter{Fail: true}
	_, err := decryptValue("SECRET", "ENC:xyz", mock)
	if err == nil {
		t.Errorf("expected wrapped decryption error, got: %v", err)
	}
}

func TestDecryptValue_Success(t *testing.T) {
	mock := &MockDecrypter{Output: "decrypted"}
	plain, err := decryptValue("SECRET", "ENC:xyz", mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plain != "decrypted" {
		t.Errorf("expected 'decrypted', got: %q", plain)
	}
}
