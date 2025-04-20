package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Reinami/configprovider/pkg/cryptography"
	"github.com/Reinami/configprovider/pkg/provider"
)

const version = "v0.1.0"

var supportedAlgorithms = map[string]string{
	"aesgcm": "AES-GCM encryption with 256-bit key",
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	fs := flag.NewFlagSet(command, flag.ExitOnError)

	var algorithm, secret, value string
	fs.StringVar(&algorithm, "c", "", "Crypto Algorithm (run --l to see a list of supported Algorithms)")
	fs.StringVar(&algorithm, "crypto-algorithm", "", "Crypto Algorithm (run --l to see a list of supported Algorithms)")
	fs.StringVar(&secret, "s", "", "Your secret key for encryption/decryption")
	fs.StringVar(&secret, "secret-key", "", "Your secret key for encryption/decryption")
	fs.StringVar(&value, "v", "", "Value to encrypt/decrypt")
	fs.StringVar(&value, "value", "", "Value to encrypt/decrypt")

	err := fs.Parse(os.Args[2:])
	if err != nil {
		panic(err)
	}

	args := fs.Args()
	if len(args) >= 2 {
		if secret == "" {
			secret = strings.TrimSpace(args[0])
		}

		if value == "" {

			value = strings.TrimSpace(args[1])
		}
	}

	if algorithm == "" || secret == "" || value == "" {
		fmt.Println("Error: --crypto-algorithm, secret, and value are required")
		showHelp()
		return
	}

	switch command {
	case "encrypt":
		encryptSecret(algorithm, secret, value)
	case "decrypt":
		decryptSecret(algorithm, secret, value)
	case "--help", "-h", "help":
		showHelp()
	case "--l", "--list-algorithms":
		showAlgorithms()
	case "--version":
		fmt.Println("lockbox version", version)
	default:
		fmt.Println("unknown command ", command)
		showHelp()
	}
}

func encryptSecret(algorithm string, secret string, value string) {
	algo, err := getAlgorithm(algorithm, secret)
	if err != nil {
		printErr(err)
		return
	}

	encryptedValue, err := algo.Encrypt(value)
	if err != nil {
		printErr(err)
		return
	}

	fmt.Println(encryptedValue)
}

func decryptSecret(algorithm string, secret string, value string) {
	algo, err := getAlgorithm(algorithm, secret)
	if err != nil {
		printErr(err)
		return
	}

	decryptedValue, err := algo.Decrypt(value)
	if err != nil {
		printErr(err)
		return
	}

	fmt.Println(decryptedValue)
}

func getAlgorithm(algorithm string, secret string) (provider.CryptoAlgorithm, error) {
	switch algorithm {
	case "aesgcm":
		return cryptography.NewAESGCMCrypto(secret)
	}

	return nil, fmt.Errorf("unsupported algorithm, %v", algorithm)
}

func showHelp() {
	fmt.Println(`Usage:
  lockbox <command> [options] <secret> <value>

Commands:
  encrypt     Encrypt a value
  decrypt     Decrypt a value

Options:
  --c, --crypto-algorithm   Required. The crypto algorithm to use (e.g., aesgcm)
  --s, --secret-key         Optional. Secret key (can be passed positionally)
  --v, --value              Optional. Value to encrypt/decrypt (can be passed positionally)
  --l, --list-algorithms    Optional. Will display a list of currently supported algorithms
  --version                 Optional. Displays current version of lockbox

Examples:
  lockbox encrypt --c=aesgcm mysecret myvalue
  lockbox decrypt --c=aesgcm mysecret myencryptedvalue`)
}

func showAlgorithms() {
	fmt.Println("Supported Algorithms:")
	for name, desc := range supportedAlgorithms {
		fmt.Printf("  %-10s - %s\n", name, desc)
	}
}

func printErr(err error) {
	fmt.Println(err)
	showHelp()
}
