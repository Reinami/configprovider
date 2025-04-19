package configloader

import (
	"reflect"
	"strings"
)

type tagOptions struct {
	Key         string // The config key to lookup
	Default     string // The default value of the config
	IsRequired  bool   // If the field is IsRequired
	IsEncrypted bool   // If the field is encrypted
}

// Example tag:
// `config:"PORT,default=8000,required,encrypted"`
func parseTag(field reflect.StructField) tagOptions {
	rawTag := field.Tag.Get("config")
	if rawTag == "" {
		return tagOptions{}
	}

	parts := strings.Split(rawTag, ",")
	options := tagOptions{
		Key: strings.TrimSpace(parts[0]),
	}

	for _, part := range parts[1:] {
		trimmedPart := strings.TrimSpace(part)
		switch {
		case trimmedPart == "required":
			options.IsRequired = true
		case trimmedPart == "encrypted":
			options.IsEncrypted = true
		case strings.HasPrefix(trimmedPart, "default="):
			options.Default = strings.TrimPrefix(trimmedPart, "default=")
		}
	}

	return options
}
