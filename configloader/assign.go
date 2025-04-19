package configloader

import (
	"fmt"
	"reflect"
)

func assignFields(target reflect.Value, source Source, cryptoAlgo CryptoAlgorithm) error {
	targetType := target.Type()

	for i := 0; i < target.NumField(); i++ {
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
