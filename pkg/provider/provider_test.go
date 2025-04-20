package provider_test

import (
	"testing"

	"github.com/Reinami/configprovider/pkg/provider"
)

type mockSource map[string]string

func (m mockSource) Get(key string) (string, bool) {
	val, ok := m[key]
	return val, ok
}

type mockDecrypter struct {
	Value string
	Err   error
}

func (m *mockDecrypter) Decrypt(_ string) (string, error) {
	return m.Value, m.Err
}

type mockConfig struct {
	AppName             string          `config:"APP_NAME,required"`
	Port                int             `config:"PORT,required,default=8000"`
	AcceptableErrorRate float32         `config:"ACCEPTABLE.ERROR.RATE,default=0.01"`
	Debug               bool            `config:"DEBUG,required"`
	Tags                []string        `config:"TAGS"`
	FeatureFlags        map[string]bool `config:"FEATURE.FLAGS"`
	SecretKey           string          `config:"SECRET,encrypted"`
	AMissingField       string          `config:"MISSING"`
	unsettable          string
}

func (m *mockConfig) GetUnsettable() string {
	return m.unsettable
}

func TestConfigProvider_Load(t *testing.T) {
	source := mockSource{
		"APP_NAME":      "TestService",
		"DEBUG":         "true",
		"TAGS":          "a,b,c",
		"FEATURE.FLAGS": "{\"featureA\": true, \"featureB\": false}",
		"SECRET":        "encrypted",
	}

	config := mockConfig{
		unsettable: "unsettable",
	}

	err := provider.NewConfigProvider().
		FromSource(source).
		WithDecrypter(&mockDecrypter{Value: "decrypted"}).
		Load(&config)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.AppName != "TestService" {
		t.Errorf("AppName mismatch: expected %v, got %v", "TestService", config.AppName)
	}

	if config.Port != 8000 {
		t.Errorf("Port mismatch: expected %v, got %v", 8000, config.Port)
	}

	if config.AcceptableErrorRate != 0.01 {
		t.Errorf("AcceptableErrorRate mismatch: expected %v, got %v", 0.01, config.AcceptableErrorRate)
	}

	if config.Debug != true {
		t.Errorf("Debug mismatch: expected %v, got %v", true, config.Debug)
	}

	if config.Tags[0] != "a" || config.Tags[1] != "b" || config.Tags[2] != "c" || len(config.Tags) != 3 {
		t.Errorf("Tags mismatch: expected %v, got %v", []string{"a", "b", "c"}, config.Tags)
	}

	if config.FeatureFlags["featureA"] != true || config.FeatureFlags["featureB"] != false {
		t.Errorf("FeatureFlags mismatch: expected %v, got %v", map[string]bool{"featureA": true, "featureB": false}, config.FeatureFlags)
	}

	if config.SecretKey != "decrypted" {
		t.Errorf("SecretKey mismatch: expected %v, got %v", "decrypted", config.SecretKey)
	}

	if config.AMissingField != "" {
		t.Errorf("AMissingField mismatch: expected %v, got %v", "", config.AMissingField)
	}

	if config.GetUnsettable() != "unsettable" {
		t.Errorf("unsettable mismatch: expected %v, got %v", "unsettable", config.GetUnsettable())
	}
}
