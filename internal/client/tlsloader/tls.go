package tlsloader

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"google.golang.org/grpc/credentials"
)

func LoadTLS(ca, cert, key string) (credentials.TransportCredentials, error) {
	caPem, err := os.ReadFile(ca)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, err
	}
	var clientCert tls.Certificate
	clientCert, err = tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	conf := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}
	return credentials.NewTLS(conf), nil
}
