// Package config
// configuring the server, logging level, database
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ServerConfig - server configuration structure
type ServerConfig struct {
	LogLevel       string `yaml:"log_level"`
	GRPCAddress    string `yaml:"grpc_address"`
	StorageAddress string `yaml:"storage_address"`
	JWTKey         string `yaml:"jwt_key"`
	DBUser         string `yaml:"db_user"`
	DBPassword     string `yaml:"db_password"`
	ServerCert     string `yaml:"server_cert"`
	ServerKey      string `yaml:"server_key"`
}

// NewServerConfig - function of obtaining the server configuration, processes the yaml file
func NewServerConfig() *ServerConfig {
	config := ServerConfig{}

	yamlFile, err := os.Open("./cmd/server/config/config.yaml")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err2 := yamlFile.Close(); err2 != nil {
			panic(err2)
		}
	}()
	decoder := yaml.NewDecoder(yamlFile)
	err = decoder.Decode(&config)
	return &config
}
