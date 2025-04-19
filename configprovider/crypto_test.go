package configprovider

import (
	"errors"
	"testing"
)

// Mocks
type mockCryptoTestDecrypter struct {
	ShouldFail     bool
	ExpectedOutput string
}

func (m *mockCryptoTestDecrypter) Decrypt(cipherText string) (string, error) {
	if m.ShouldFail {
		return "", errors.New("mock decryption error")
	}
	return m.ExpectedOutput, nil
}

// Tests

func TestDecryptValue_NoDecrypter(t *testing.T) {
	_, err := decryptValue("SECRET", "key", nil)
	if err == nil || err.Error() != "no decrypter is provided" {
		t.Errorf("expected 'no decrypter is provided', got: %v", err)
	}
}

func TestDecryptValue_DecryptionFails(t *testing.T) {
	mock := &mockCryptoTestDecrypter{ShouldFail: true}

	_, err := decryptValue("SECRET", "key", mock)
	if err == nil {
		t.Errorf("expected wrapped decryption error, got: %v", err)
	}
}

func TestDecryptValue_Success(t *testing.T) {
	mock := &mockCryptoTestDecrypter{ExpectedOutput: "decrypted"}

	plain, err := decryptValue("SECRET", "key", mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plain != "decrypted" {
		t.Errorf("expected 'decrypted', got: %q", plain)
	}
}
