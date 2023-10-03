// Package tlsloader
// in this file loading tls certificates and keys
package tlsloader

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/h2p2f/dedicated-vault/internal/server/config"
)

// LoadTLS - function for loading tls certificates and keys
func LoadTLS(config *config.ServerConfig) (*tls.Config, error) {
	caPem, err := os.ReadFile("./crypto/ca-cert.pem")
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, err
	}
	serverCert, err := tls.LoadX509KeyPair(config.ServerCert, config.ServerKey)
	if err != nil {
		return nil, err
	}
	conf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}
	return conf, nil
}
