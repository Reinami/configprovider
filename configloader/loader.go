package configloader

import (
	"fmt"
	"reflect"
)

func Load(config any, source Source, decrypter Decrypter) error {
	reflectValue := reflect.ValueOf(config)

	if reflectValue.Kind() != reflect.Ptr || reflectValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("load expects a pointer to a struct and got %T", config)
	}

	structValue := reflectValue.Elem()
	return assignFields(structValue, source, decrypter)
}
