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
	ShouldFail     bool
	ExpectedOutput string
}

func (m *mockParseTestDecrypter) Decrypt(cipherText string) (string, error) {
	if m.ShouldFail {
		return "", errors.New("mock decryption error")
	}
	return m.ExpectedOutput, nil
}

type mockParseTestConfig struct {
	TestField       bool   `config:"TEST_FIELD"`
	DefaultField    string `config:"DEFAULT_FIELD,default=somethingDefault"`
	RequiredField   string `config:"REQUIRED_FIELD,required"`
	EncryptedField  string `config:"ENCRYPTED_FIELD,encrypted"`
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

	err := assignFields(value, &mockParseTestSource{"REQUIRED_FIELD": "required", "ENCRYPTED_FIELD": "encrypted"}, &mockParseTestDecrypter{ShouldFail: true})
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
		},
		&mockCryptoTestDecrypter{ExpectedOutput: "encrypted"},
	)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if config.TestField != true ||
		config.RequiredField != "required" ||
		config.EncryptedField != "encrypted" ||
		config.DefaultField != "somethingDefault" {
		t.Errorf("got incorrect fields on config %v", config)
	}
}
