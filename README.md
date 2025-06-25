# configprovider

A flexible, tag-based configuration loader for Go.  
Supports `.properties` files, encrypted fields, and custom source + crypto integrations.

---

## Features

- Struct-based configuration loading
- currently has `.properties` file support
- Optional field defaults, required fields, and encryption
- Extendable via custom sources or decryption strategies
- Optional CLI helper: [`lockbox`](#-lockbox-cli-optional)

---

## Quick Start

### 1. Define Your Config Struct

```go
type AppConfig struct {
  Port     int     `config:"PORT,default=8080"`
  Debug    bool    `config:"DEBUG,default=false"`
  Timeout  float64 `config:"TIMEOUT,required"`
  Secret   string  `config:"SECRET_KEY,required,encrypted"`
}
```

### 2. Load Your Config

```go
config := AppConfig{}
err := configprovider.New().
  FromPropertiesFile("app.properties").
  WithAESGCMDecrypter("your-32-byte-secret-key").
  Load(&config)
if err != nil {
  log.Fatal(err)
}
```

---

## Custom Source

Implement the `configprovider.Source` interface:

```go
type MyCustomSource struct{}

func (m *MyCustomSource) Get(key string) (string, bool) {
  // fetch from SSM, Vault, env, etc.
  return "value", true
}
```

Then:

```go
configprovider.New().FromSource(&MyCustomSource{})
```

---

## Custom Encrypter / Decrypter

Create a type that implements:

```go
type Decrypter interface {
  Decrypt(cipherText string) (string, error)
}

type Encrypter interface {
  Encrypt(plainText string) (string, error)
}
```

Example:

```go
type MyBase64Decrypter struct{}

func (d *MyBase64Decrypter) Decrypt(input string) (string, error) {
  data, err := base64.StdEncoding.DecodeString(input)
  return string(data), err
}
```

Use it like:

```go
configprovider.New().
  FromFile("config.properties").
  WithDecrypter(&MyBase64Decrypter{}).
  Load(&cfg)
```

---

## `lockbox` CLI (optional)

A helper CLI to encrypt/decrypt values using the same algorithms used by `configprovider`.

### Install (Go required)

```bash
go install github.com/Reinami/configprovider/cmd/lockbox@latest
```

### Usage

```bash
# Encrypt
lockbox encrypt --c=aesgcm mysecret mysecretvalue

# Decrypt
lockbox decrypt --c=aesgcm mysecret ciphertext

# Show available algorithms
lockbox --list-algorithms
```

---

## Field Tags

| Tag           | Description                                             |
|---------------|---------------------------------------------------------|
| `config:"KEY"` | Defines the source key name                            |
| `default=...` | Optional default value if key is missing                |
| `required`    | Fail if the key is missing and no default is provided  |
| `encrypted`   | Decrypt the value using the configured decrypter       |

---

## License

MIT License. See [LICENSE](./LICENSE).

