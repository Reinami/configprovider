package cryptography

import (
	"strings"
	"testing"
)

const testKey string = "12345678901234567890123456789012"

func TestAESGCM_EncryptionAndDecryption(t *testing.T) {
	crypto, err := NewAESGCMCrypto(testKey)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

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

func TestAESGCM_EncryptionProducesUnique(t *testing.T) {
	crypto, err := NewAESGCMCrypto(testKey)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

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

func TestAESGCM_DecryptionInvalidBase64(t *testing.T) {
	crypto, err := NewAESGCMCrypto(testKey)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = crypto.Decrypt("something that isn't base64")
	if err == nil || !strings.Contains(err.Error(), "base64 decode failed") {
		t.Errorf("expected base64 decode error, got: %v", err)
	}
}

func TestAESGCM_DecryptionShortNonce(t *testing.T) {
	crypto, err := NewAESGCMCrypto(testKey)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	shortInput := "dGVzdA==" // base64 for "test"
	_, err = crypto.Decrypt(shortInput)

	if err == nil || !strings.Contains(err.Error(), "ciphertext too short") {
		t.Errorf("expected short ciphertext error, got: %v", err)
	}
}

func TestAESGCM_InvalidKeyLength(t *testing.T) {
	_, err := NewAESGCMCrypto("invalid-ley")
	if err == nil {
		t.Errorf("expected error for invalid key length")
	}
}
