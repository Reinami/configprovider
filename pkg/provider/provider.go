package provider

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Reinami/configprovider/pkg/cryptography"
	"github.com/Reinami/configprovider/pkg/sources"
)

type Source interface {
	Get(key string) (string, bool)
}

type configProvider struct {
	source    Source
	decrypter Decrypter
}

// Source options

func (c *configProvider) FromSource(source Source) *configProvider {
	c.source = source
	return c
}

func (c *configProvider) FromFile(path string) *configProvider {
	extension := strings.ToLower(filepath.Ext(path))

	switch extension {
	case ".properties":
		return c.FromPropertiesFile(path)
	}

	panic("Unsupported file type: " + extension)
}

func (c *configProvider) FromPropertiesFile(path string) *configProvider {
	source, err := sources.NewPropertiesFileSource(path)
	if err != nil {
		panic(err)
	}

	c.source = source
	return c
}

// Decrypter options

func (c *configProvider) WithDecrypter(decrypter Decrypter) *configProvider {
	c.decrypter = decrypter
	return c
}

func (c *configProvider) WithAESGCMDecrypter(secretKey string) *configProvider {
	aesGCMDecrypter, err := cryptography.NewAESGCMCrypto(secretKey)
	if err != nil {
		panic(err)
	}

	c.decrypter = aesGCMDecrypter
	return c
}

func (c *configProvider) Load(configStruct any) error {
	reflectValue := reflect.ValueOf(configStruct)

	if reflectValue.Kind() != reflect.Ptr || reflectValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("load expects a pointer to a struct and got %T", configStruct)
	}

	structValue := reflectValue.Elem()
	return assignFields(structValue, c.source, c.decrypter)
}

// Constructor

func NewConfigProvider() *configProvider {
	return &configProvider{}
}
