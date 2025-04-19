package configloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func parseAndSetValue(field reflect.Value, rawValue string) error {
	if !field.CanSet() {
		return errors.New("field is not settable")
	}

	switch field.Kind() {

	case reflect.String:
		field.SetString(rawValue)
		return nil

	case reflect.Bool:
		parsedValue, err := strconv.ParseBool(rawValue)
		if err != nil {
			return fmt.Errorf("unable to parse bool: %w", err)
		}
		field.SetBool(parsedValue)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bitSize := field.Type().Bits()
		parsedValue, err := strconv.ParseInt(rawValue, 10, bitSize)
		if err != nil {
			return fmt.Errorf("unable to parse int: %w", err)
		}

		field.SetInt(parsedValue)
		return nil

	case reflect.Float32, reflect.Float64:
		bitSize := field.Type().Bits()
		parsedValue, err := strconv.ParseFloat(rawValue, bitSize)
		if err != nil {
			return fmt.Errorf("unable to parse float: %w", err)
		}
		field.SetFloat(parsedValue)
		return nil

	case reflect.Slice:
		return parseAndSetSlice(field, rawValue)

	case reflect.Map:
		return parseAndSetMap(field, rawValue)

	default:
		return fmt.Errorf("unsupported type: %s", field.Kind())
	}
}

func parseAndSetSlice(field reflect.Value, rawValue string) error {
	elemType := field.Type().Elem()
	items := strings.Split(rawValue, ",")
	slice := reflect.MakeSlice(field.Type(), 0, len(items))

	for _, item := range items {
		item := strings.TrimSpace(item)
		elem := reflect.New(elemType).Elem()

		err := parseAndSetValue(elem, item)
		if err != nil {
			return fmt.Errorf("invalid slice element: %w", err)
		}

		slice = reflect.Append(slice, elem)
	}

	field.Set(slice)
	return nil
}

func parseAndSetMap(field reflect.Value, rawValue string) error {
	mapType := field.Type()
	keyType := mapType.Key()
	valueType := mapType.Elem()

	var tmpMap map[string]json.RawMessage
	err := json.Unmarshal([]byte(rawValue), &tmpMap)
	if err != nil {
		return fmt.Errorf("unable to unmarshal JSON map: %w", err)
	}

	result := reflect.MakeMap(mapType)

	for keyString, raw := range tmpMap {
		key := reflect.New(keyType).Elem()
		err = parseAndSetValue(key, keyString)
		if err != nil {
			return fmt.Errorf("unable to convert map key: %s, %w", keyString, err)
		}

		value := reflect.New(valueType).Elem()
		err = json.Unmarshal(raw, value.Addr().Interface())
		if err != nil {
			return fmt.Errorf("unable to unmarshal map for key: %s, %w", keyString, err)
		}

		result.SetMapIndex(key, value)
	}

	field.Set(result)
	return nil
}
