package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log/slog"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// GenerateCertificates creates client certificates for OPC UA authentication.
// It generates client_cert.der, client_cert.pem, client_key.der, and client_key.pem
// in the specified directory. If the directory is empty, it uses the working directory.
func GenerateCertificates(dir string) error {
	if dir == "" {
		// Use the working directory (project root during development)
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
		dir = wd
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create certificate directory: %w", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate RSA key: %w", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("generate serial number: %w", err)
	}

	uri, err := url.Parse("urn:go-opcua-client")
	if err != nil {
		return fmt.Errorf("parse URI: %w", err)
	}

	tmpl := x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "GoOPCUAClient", Organization: []string{"7thSense"}},
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage: x509.KeyUsageDigitalSignature |
			x509.KeyUsageContentCommitment |
			x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDataEncipherment |
			x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		URIs:                  []*url.URL{uri},
	}

	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("create certificate: %w", err)
	}

	keyDER := x509.MarshalPKCS1PrivateKey(key)

	// Write DER files
	certDERPath := filepath.Join(dir, "client_cert.der")
	if err := os.WriteFile(certDERPath, der, 0644); err != nil {
		return fmt.Errorf("write client_cert.der: %w", err)
	}

	keyDERPath := filepath.Join(dir, "client_key.der")
	if err := os.WriteFile(keyDERPath, keyDER, 0600); err != nil {
		return fmt.Errorf("write client_key.der: %w", err)
	}

	// Write PEM files
	certPEMPath := filepath.Join(dir, "client_cert.pem")
	certOut, err := os.Create(certPEMPath)
	if err != nil {
		return fmt.Errorf("create client_cert.pem: %w", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		certOut.Close()
		return fmt.Errorf("encode client_cert.pem: %w", err)
	}
	certOut.Close()

	keyPEMPath := filepath.Join(dir, "client_key.pem")
	keyOut, err := os.OpenFile(keyPEMPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("create client_key.pem: %w", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDER}); err != nil {
		keyOut.Close()
		return fmt.Errorf("encode client_key.pem: %w", err)
	}
	keyOut.Close()

	slog.Info("certificates generated",
		"cert_der", certDERPath,
		"key_der", keyDERPath,
		"cert_pem", certPEMPath,
		"key_pem", keyPEMPath,
	)

	return nil
}
