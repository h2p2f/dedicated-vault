package tlsloader

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"os"
)

func LoadTLS() (credentials.TransportCredentials, error) {
	caPem, err := os.ReadFile("./crypto/ca-cert.pem")
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, err
	}
	var clientCert tls.Certificate
	clientCert, err = tls.LoadX509KeyPair("./crypto/client-cert.pem", "./crypto/client-key.pem")
	if err != nil {
		return nil, err
	}
	conf := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}
	return credentials.NewTLS(conf), nil
}
