package config

type ServerConfig struct {
	LogLevel       string `yaml:"log_level"`
	GRPCAddress    string `yaml:"grpc_address"`
	StorageAddress string `yaml:"storage_address"`
	JWTKey         string `yaml:"jwt_key"`
}
