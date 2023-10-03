// Package: config
package config

import (
	"runtime"

	"google.golang.org/grpc/credentials"
)

// this variable is set by ldflags
var (
	version   = "0.0.3"
	buildDate = "2023-10-03"
	dbPath    = "/tmp/vault.db"
	ca        = "./crypto/ca-cert.pem"
	cert      = "./crypto/client-cert.pem"
	key       = "./crypto/client-key.pem"
)

// ClientConfig is a struct for client configuration
// yaml tags currently not used
type ClientConfig struct {
	StorageAddress    string `yaml:"storage_address"`
	DBPath            string `yaml:"db_path"`
	Secret            string `yaml:"secret"`
	User              string `yaml:"user"`
	Passphrase        string `yaml:"pass_phrase"`
	Token             string `yaml:"token"`
	ClientCA          string `yaml:"client_ca"`
	ClientCert        string `yaml:"client_cert"`
	ClientKey         string `yaml:"client_key"`
	CryptoKey         []byte `yaml:"crypto_key"`
	IsLoggedIn        bool   `yaml:"is_logged_in"`
	LastServerUpdated int64  `yaml:"last_server_updated"`
	TLSConfig         credentials.TransportCredentials
	Version           string `yaml:"version"`
	BuildDate         string `yaml:"build_date"`
}

// NewClientConfig - function of obtaining the client configuration
func NewClientConfig() *ClientConfig {
	if runtime.GOOS == "windows" {
		dbPath = "C:\\Users\\Public\\vault.db"
		ca = "C:\\Users\\Public\\ca-cert.pem"
		cert = "C:\\Users\\Public\\client-cert.pem"
		key = "C:\\Users\\Public\\client-key.pem"
	}

	return &ClientConfig{
		StorageAddress: "localhost:8090",
		DBPath:         dbPath,
		Passphrase:     "",
		IsLoggedIn:     false,
		Token:          "",
		ClientCA:       ca,
		ClientCert:     cert,
		ClientKey:      key,
		TLSConfig:      nil,
		Version:        version,
		BuildDate:      buildDate,
	}
}
