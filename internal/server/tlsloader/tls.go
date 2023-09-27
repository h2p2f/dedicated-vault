package tlsloader

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func LoadTLS() (*tls.Config, error) {
	caPem, err := os.ReadFile("./crypto/ca-cert.pem")
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, err
	}
	serverCert, err := tls.LoadX509KeyPair("./crypto/server-cert.pem", "./crypto/server-key.pem")
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
