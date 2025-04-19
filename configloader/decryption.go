package configloader

import "fmt"

type CryptoAlgorithm interface {
	Encrypt(plainText string) (string, error)
	Decrypt(cipherText string) (string, error)
}

func decryptValue(key string, value string, cryptoAlgo CryptoAlgorithm) (string, error) {
	if cryptoAlgo == nil {
		return "", fmt.Errorf("field %s is marked as encrypted but no cryptoAlgorithm is provided", key)
	}

	plainText, err := cryptoAlgo.Decrypt(value)
	if err != nil {
		return "", fmt.Errorf("decryption failed for %s: %w", key, err)
	}

	return plainText, nil
}
