package cryptography

import (
	"strings"
	"testing"
)

const testKey string = "12345678901234567890123456789012"

func TestEncryptionAndDecryption(t *testing.T) {
	crypto := NewAESGCMCrypto(testKey)

	original := "super-secret-value"
	encrypted, err := crypto.Encrypt(original)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	decrypted, err := crypto.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("descryption failed: %v", err)
	}

	if decrypted != original {
		t.Errorf("expected decrypted text to be %s, got %s", original, decrypted)
	}
}

func TestEncryptionProducesUnique(t *testing.T) {
	crypto := NewAESGCMCrypto(testKey)

	secret := "super-secret-value"
	encryption1, err1 := crypto.Encrypt(secret)
	encryption2, err2 := crypto.Encrypt(secret)

	if err1 != nil || err2 != nil {
		t.Fatalf("encryption failed: %v, %v", err1, err2)
	}

	if encryption1 == encryption2 {
		t.Errorf("expected different outputs for same input due to nonce generation, got identical outputs: %s, %s", encryption1, encryption2)
	}
}

func TestDecryptionInvalidBase64(t *testing.T) {
	crypto := NewAESGCMCrypto(testKey)

	_, err := crypto.Decrypt("something that isn't base64")
	if err == nil || !strings.Contains(err.Error(), "base64 decode failed") {
		t.Errorf("expected base64 decode error, got: %v", err)
	}
}

func TestDecryptionShortNonce(t *testing.T) {
	crypto := NewAESGCMCrypto(testKey)

	shortInput := "dGVzdA==" // base64 for "test"
	_, err := crypto.Decrypt(shortInput)

	if err == nil || !strings.Contains(err.Error(), "ciphertext too short") {
		t.Errorf("expected short ciphertext error, got: %v", err)
	}
}

func TestInvalidKeyLength(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for invalid key length")
		}
	}()

	NewAESGCMCrypto("invalid-key")
}
