package configprovider

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// Mocks

type mockParseTestSource map[string]string

func (m mockParseTestSource) Get(key string) (string, bool) {
	val, ok := m[key]
	return val, ok
}

type mockParseTestDecrypter struct {
	Err   error
	Value string
}

func (m *mockParseTestDecrypter) Decrypt(cipherText string) (string, error) {
	return m.Value, m.Err
}

type mockParseTestConfig struct {
	TestBoolField   bool              `config:"TEST_FIELD"`
	DefaultField    string            `config:"DEFAULT_FIELD,default=somethingDefault"`
	RequiredField   string            `config:"REQUIRED_FIELD,required"`
	EncryptedField  string            `config:"ENCRYPTED_FIELD,encrypted"`
	MapField        map[string]string `config:"MAP_FIELD"`
	ListField       []int             `config:"LIST_FIELD"`
	IntField        int               `config:"INT_FIELD"`
	FloatField      float32           `config:"FLOAT_FIELD"`
	unsettableField string
	MissingKey      string
}

// Tests

func TestAssignFields_MissingRequired(t *testing.T) {
	config := mockParseTestConfig{}

	value := reflect.ValueOf(&config).Elem()

	err := assignFields(value, &mockParseTestSource{}, nil)
	if err == nil || !strings.Contains(err.Error(), "required key REQUIRED_FIELD is missing") {
		t.Errorf("expected error for missing field but got %v", err)
	}
}

func TestAssignFields_DecryptionError(t *testing.T) {
	config := mockParseTestConfig{}

	value := reflect.ValueOf(&config).Elem()

	err := assignFields(
		value,
		&mockParseTestSource{
			"REQUIRED_FIELD":  "required",
			"ENCRYPTED_FIELD": "encrypted",
		},
		&mockParseTestDecrypter{Err: errors.New("decryption error")},
	)
	if err == nil || !strings.Contains(err.Error(), "decryption failed for ENCRYPTED_FIELD") {
		t.Errorf("expected error from decrypter but got %v", err)
	}
}

func TestAssignFields_ParseError(t *testing.T) {
	config := mockParseTestConfig{}

	value := reflect.ValueOf(&config).Elem()

	err := assignFields(value, &mockParseTestSource{"TEST_FIELD": "notabool"}, nil)
	if err == nil || !strings.Contains(err.Error(), "unable to parse") {
		t.Errorf("expected error from parser but got %v", err)
	}
}

func TestAssignFields_Success(t *testing.T) {
	config := mockParseTestConfig{}

	value := reflect.ValueOf(&config).Elem()

	err := assignFields(
		value,
		&mockParseTestSource{
			"TEST_FIELD":      "true",
			"REQUIRED_FIELD":  "required",
			"ENCRYPTED_FIELD": "encrypted",
			"MAP_FIELD":       "{\"someKey\": \"someValue\", \"someOtherKey\": \"someOtherValue\"}",
			"LIST_FIELD":      "0, 1,2",
			"INT_FIELD":       "7",
			"FLOAT_FIELD":     "2.3",
		},
		&mockParseTestDecrypter{Value: "encrypted"},
	)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if config.TestBoolField != true ||
		config.RequiredField != "required" ||
		config.EncryptedField != "encrypted" ||
		config.MapField["someKey"] != "someValue" ||
		config.MapField["someOtherKey"] != "someOtherValue" ||
		config.ListField[0] != 0 ||
		config.ListField[1] != 1 ||
		config.ListField[2] != 2 ||
		config.IntField != 7 ||
		config.FloatField != 2.3 ||
		config.DefaultField != "somethingDefault" {
		t.Errorf("got incorrect fields on config %v", config)
	}
}
