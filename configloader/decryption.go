package configloader

import "fmt"

type Decrypter interface {
	Decrypt(cipherText string) (string, error)
}

func decryptValue(key string, value string, decrypter Decrypter) (string, error) {
	if decrypter == nil {
		return "", fmt.Errorf("field %s is marked as encrypted but no decrypter is provided", key)
	}

	plainText, err := decrypter.Decrypt(value)
	if err != nil {
		return "", fmt.Errorf("decryption failed for %s: %w", key, err)
	}

	return plainText, nil
}
