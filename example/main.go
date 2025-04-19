package main

import (
	"fmt"

	"github.com/Reinami/configprovider"
)

type AppConfig struct {
	Debug     bool            `config:"DEBUG"`
	Port      int             `config:"PORT,default=8080"`
	Name      string          `config:"NAME,default=defaultName"`
	SecretKey string          `config:"SECRET_KEY,encrypted"`
	Tags      []string        `config:"TAGS"`
	Flags     map[string]bool `config:"FLAGS"`
}

func main() {
	var config AppConfig

	// This should not be hard coded or its not secure!
	const secretKey string = "12345678901234567890123456789012"

	err := configprovider.NewConfigProvider().
		FromPropertiesFile("./testconfig.properties").
		WithAESGCMDecrypter(secretKey).
		Load(&config)

	if err != nil {
		panic(err)
	}

	fmt.Println("Loaded Config:")
	fmt.Printf("  Debug:       %v\n", config.Debug)
	fmt.Printf("  Port:        %v\n", config.Port)
	fmt.Printf("  Name:        %v\n", config.Name)
	fmt.Printf("  SecretKey:   %v\n", config.SecretKey)
	fmt.Printf("  Tags:        %v\n", config.Tags)
	fmt.Printf("  Flags:       %v\n", config.Flags)
}
