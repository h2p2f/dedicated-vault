package config

import "crypto/x509"

type ClientConfig struct {
	StorageAddress string `yaml:"storage_address"`
	DBPath         string `yaml:"db_path"`
	Secret         string `yaml:"secret"`
	User           string `yaml:"user"`
	Passphrase     string `yaml:"pass_phrase"`
	Token          string `yaml:"token"`
	CryptoKey      []byte `yaml:"crypto_key"`
	IsLoggedIn     bool   `yaml:"is_logged_in"`
	Cert           *x509.CertPool
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		StorageAddress: "localhost:8090",
		DBPath:         "/tmp/vault.db",
		Passphrase:     "",
		IsLoggedIn:     false,
		Token:          "",
		Cert:           nil,
	}
}

func (c ClientConfig) SetPassphrase(passphrase string) {
	c.Passphrase = passphrase
}
