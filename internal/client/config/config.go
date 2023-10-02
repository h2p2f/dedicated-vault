package config

import (
	"google.golang.org/grpc/credentials"
	"runtime"
)

type ClientConfig struct {
	StorageAddress    string `yaml:"storage_address"`
	DBPath            string `yaml:"db_path"`
	Secret            string `yaml:"secret"`
	User              string `yaml:"user"`
	Passphrase        string `yaml:"pass_phrase"`
	Token             string `yaml:"token"`
	CryptoKey         []byte `yaml:"crypto_key"`
	IsLoggedIn        bool   `yaml:"is_logged_in"`
	LastServerUpdated int64  `yaml:"last_server_updated"`
	TLSConfig         credentials.TransportCredentials
}

func NewClientConfig() *ClientConfig {
	dbPath := "/tmp/vault.db"
	if runtime.GOOS == "windows" {
		dbPath = "C:\\Users\\Public\\vault.db"
	}

	return &ClientConfig{
		StorageAddress: "localhost:8090",
		DBPath:         dbPath,
		Passphrase:     "",
		IsLoggedIn:     false,
		Token:          "",
		TLSConfig:      nil,
	}
}
