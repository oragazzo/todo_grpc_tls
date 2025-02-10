package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
)

// TLSConfig holds the paths for TLS certificates and keys
type TLSConfig struct {
	CertFile   string
	KeyFile    string
	CAFile     string
	ServerName string
}

// NewTLSConfigFromEnv creates a TLSConfig from environment variables
func NewTLSConfigFromEnv() (*TLSConfig, error) {
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")
	caFile := os.Getenv("TLS_CA_FILE")
	serverName := os.Getenv("TLS_SERVER_NAME")

	if certFile == "" || keyFile == "" || caFile == "" {
		return nil, fmt.Errorf("missing required TLS environment variables")
	}

	if serverName == "" {
		serverName = "localhost" // default value
	}

	return &TLSConfig{
		CertFile:   certFile,
		KeyFile:    keyFile,
		CAFile:     caFile,
		ServerName: serverName,
	}, nil
}

// LoadTLSConfig loads TLS configuration from the provided certificate files
func LoadTLSConfig(cfg TLSConfig) *tls.Config {
	// Load certificate of the CA who signed server's certificate
	caCert, err := os.ReadFile(cfg.CAFile)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		log.Fatal("Failed to append CA certificate to pool")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Fatalf("Failed to load client certificate and key: %v", err)
	}

	// Create and return the TLS configuration
	return &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{clientCert},
		ServerName:   cfg.ServerName,
	}
}
