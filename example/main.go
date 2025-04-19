package main

import (
	"fmt"
	"log"

	"github.com/Reinami/configloader"
	"github.com/Reinami/configloader/cryptography"
	"github.com/Reinami/configloader/sources"
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

	source, err := sources.FromPropertiesFile("./testconfig.properties")
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}

	// This should not be hard coded or its not secure!
	const secretKey string = "12345678901234567890123456789012"

	crypto := cryptography.NewAESGCMCrypto(secretKey)
	err = configloader.Load(&config, source, crypto)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Println("Loaded Config:")
	fmt.Printf("  Debug:       %v\n", config.Debug)
	fmt.Printf("  Port:        %v\n", config.Port)
	fmt.Printf("  Name:        %v\n", config.Name)
	fmt.Printf("  SecretKey:   %v\n", config.SecretKey)
	fmt.Printf("  Tags:        %v\n", config.Tags)
	fmt.Printf("  Flags:       %v\n", config.Flags)
}
