package main

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/openware/pkg/ika"
	"golang.org/x/crypto/ssh/terminal"
)

type Config struct {
	PlatformBaseUrl  string `env:"OPENDAX_BASE_URL"`
	OpendaxApiKey    string `env:"OPENDAX_API_KEY"`
	OpendaxApiSecret string `env:"OPENDAX_API_SECRET"`
	BinanceApiKey    string `env:"BINANCE_API_KEY"`
	BinanceSecret    string `env:"BINANCE_SECRET"`
}

func readConfig() *Config {
	config := &Config{}
	err := ika.ReadConfig("", config)
	if err != nil {
		panic(err)
	}
	return config
}

func fetchBinanceKey(config *Config) {
	if config.BinanceApiKey == "" {
		fmt.Print("Enter BINANCE_API_KEY: ")
		apiKey, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		fmt.Println("")

		config.BinanceApiKey = strings.TrimSpace(string(apiKey))
	}

	if config.BinanceSecret == "" {
		fmt.Print("Enter BINANCE_SECRET: ")
		secretKey, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		fmt.Println("")

		config.BinanceSecret = strings.TrimSpace(string(secretKey))
	}
}
