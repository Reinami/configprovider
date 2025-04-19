package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type AESGCMCrypto struct {
	key []byte
}

func NewAESGCMCrypto(key string) (*AESGCMCrypto, error) {
	if len(key) != 32 {
		return &AESGCMCrypto{}, errors.New("AESGCMDecrypter: key must be exactly 32 bytes (AES-256)")
	}

	return &AESGCMCrypto{
		key: []byte(key),
	}, nil
}

func (c *AESGCMCrypto) Encrypt(plainText string) (string, error) {
	aesCipher, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("failed to create aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", fmt.Errorf("failed to created gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	cipherText := gcm.Seal(nil, nonce, []byte(plainText), nil)
	final := append(nonce, cipherText...)

	return base64.StdEncoding.EncodeToString(final), nil
}

func (c *AESGCMCrypto) Decrypt(base64CipherText string) (string, error) {
	cipherData, err := base64.StdEncoding.DecodeString(base64CipherText)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	aesCipher, err := aes.NewCipher(c.key)
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
