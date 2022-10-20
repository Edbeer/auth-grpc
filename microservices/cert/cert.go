package cert

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"google.golang.org/grpc/credentials"
)

// Load TLS Credentials for Server
func LoadTLSCredentialsServer() (credentials.TransportCredentials, error) {
	// Load certificate of the certificate authority (CA) who signed server's certificate
	pemServerCA, err := os.ReadFile("/home/pasha/go/src/auth-grpc/microservices/cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, err
	}
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	cfg := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(cfg), nil
}

// Load TLS Credentials for client
func LoadTLSCredentialsClient() (credentials.TransportCredentials, error) {
	// Load certificate of the certificate authority (CA) who signed server's certificate
	pemServerCA, err := os.ReadFile("/home/pasha/go/src/auth-grpc/microservices/cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, err
	}

	// Create the credentials and return it
	cfg := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(cfg), nil
}
