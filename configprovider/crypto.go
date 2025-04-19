package configprovider

import (
	"fmt"
)

type Encrypter interface {
	Encrypt(plainText string) (string, error)
}

type Decrypter interface {
	Decrypt(cipherText string) (string, error)
}

func decryptValue(key string, value string, decrypter Decrypter) (string, error) {
	if decrypter == nil {
		return "", fmt.Errorf("no decrypter is provided")
	}

	plainText, err := decrypter.Decrypt(value)
	if err != nil {
		return "", fmt.Errorf("decryption failed for %s: %w", key, err)
	}

	return plainText, nil
}
