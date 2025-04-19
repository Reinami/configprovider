package decrypters

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

type AESGCMDecrypter struct {
	Key []byte
}

func NewAESGCMDecrypter(key []byte) *AESGCMDecrypter {
	if len(key) != 32 {
		panic("AESGCMDecrypter: Key must be exactly 32 bytes (AES-256)")
	}

	return &AESGCMDecrypter{
		Key: key,
	}
}

func (d *AESGCMDecrypter) Decrypt(base64CipherText string) (string, error) {
	cipherData, err := base64.StdEncoding.DecodeString(base64CipherText)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	aesCipher, err := aes.NewCipher(d.Key)
	if err != nil {
		return "", fmt.Errorf("failed to create aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", fmt.Errorf("failed to create gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(cipherData) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce := cipherData[:nonceSize]
	cipherText := cipherData[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plainText), nil
}
