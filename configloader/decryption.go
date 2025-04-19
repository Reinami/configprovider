package configloader

import (
	"fmt"
	"reflect"
)

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

func assignFields(target reflect.Value, source Source, cryptoAlgo CryptoAlgorithm) error {
	targetType := target.Type()

	for i := range target.NumField() {
		var finalValue string

		field := target.Field(i)
		fieldType := targetType.Field(i)

		if !field.CanSet() {
			continue
		}

		tagOpts := parseTag(fieldType)
		if tagOpts.Key == "" {
			continue
		}

		finalValue, found := source.Get(tagOpts.Key)
		if !found {
			if tagOpts.Default != "" {
				finalValue = tagOpts.Default
			} else if tagOpts.IsRequired {
				return fmt.Errorf("required key %s is missing", tagOpts.Key)
			} else {
				continue
			}
		}

		if tagOpts.IsEncrypted {
			decryptedValue, err := decryptValue(tagOpts.Key, finalValue, cryptoAlgo)
			if err != nil {
				return err
			}
			finalValue = decryptedValue
		}

		err := parseAndSetValue(field, finalValue)
		if err != nil {
			return err
		}
	}

	return nil
}
