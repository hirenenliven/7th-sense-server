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

// ApplicationURI is the OPC UA application URI in the client certificate. An OPC UA client
// must announce the same URI (opcua.Settings.ApplicationURI), or the controller refuses it.
const ApplicationURI = "urn:go-opcua-client"

// certificateValidity is how long a new client certificate is valid. The NCU trusts this exact
// certificate, so a short life would silently stop every Siemens connection when it ends.
const certificateValidity = 10 * 365 * 24 * time.Hour

// Paths are the client certificate and key files (DER) the CNC DLL and the OPC UA PLC reader
// present to Siemens controllers.
type Paths struct {
	CertificateFile string
	PrivateKeyFile  string
}

// EnsureCertificates makes sure the OPC UA client certificate exists in dir (the working
// directory when empty) and returns its files. An existing, unexpired certificate is kept: the
// controllers trust that exact certificate (NCU Access Administration), so a new one at every
// start would break every connection until it is trusted again. A missing or expired one is
// created: client_cert.der / .pem and client_key.der / .pem.
func EnsureCertificates(dir string) (Paths, error) {
	if dir == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			return Paths{}, fmt.Errorf("get working directory: %w", err)
		}
		dir = workingDir
	}
	paths := Paths{
		CertificateFile: filepath.Join(dir, "client_cert.der"),
		PrivateKeyFile:  filepath.Join(dir, "client_key.der"),
	}

	expiresAt, err := existingCertificateExpiry(paths)
	switch {
	case err == nil && time.Now().Before(expiresAt):
		slog.Info("client certificate kept", "cert_der", paths.CertificateFile, "expiresAt", expiresAt)
		return paths, nil
	case err == nil:
		slog.Warn("client certificate expired, creating a new one: trust it on every Siemens NCU again",
			"cert_der", paths.CertificateFile, "expiredAt", expiresAt)
	case !os.IsNotExist(err):
		slog.Warn("client certificate unreadable, creating a new one: trust it on every Siemens NCU again",
			"cert_der", paths.CertificateFile, "error", err)
	}
	if err := generateCertificates(dir, paths); err != nil {
		return Paths{}, err
	}
	return paths, nil
}

// existingCertificateExpiry returns when the certificate in paths ends; an error when it or its
// key is missing (os.IsNotExist) or cannot be read.
func existingCertificateExpiry(paths Paths) (time.Time, error) {
	if _, err := os.Stat(paths.PrivateKeyFile); err != nil {
		return time.Time{}, err
	}
	der, err := os.ReadFile(paths.CertificateFile)
	if err != nil {
		return time.Time{}, err
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse client certificate: %w", err)
	}
	return certificate.NotAfter, nil
}

// generateCertificates writes a new self-signed client certificate and key to paths, with PEM
// copies beside them.
func generateCertificates(dir string, paths Paths) error {
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

	uri, err := url.Parse(ApplicationURI)
	if err != nil {
		return fmt.Errorf("parse URI: %w", err)
	}

	tmpl := x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "GoOPCUAClient", Organization: []string{"7thSense"}},
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(certificateValidity),
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

	if err := os.WriteFile(paths.CertificateFile, der, 0644); err != nil {
		return fmt.Errorf("write client_cert.der: %w", err)
	}
	if err := os.WriteFile(paths.PrivateKeyFile, keyDER, 0600); err != nil {
		return fmt.Errorf("write client_key.der: %w", err)
	}

	certPEMPath := filepath.Join(dir, "client_cert.pem")
	if err := writePEM(certPEMPath, "CERTIFICATE", der, 0644); err != nil {
		return err
	}
	keyPEMPath := filepath.Join(dir, "client_key.pem")
	if err := writePEM(keyPEMPath, "RSA PRIVATE KEY", keyDER, 0600); err != nil {
		return err
	}

	slog.Info("client certificate created: trust client_cert.der on every Siemens NCU",
		"cert_der", paths.CertificateFile,
		"key_der", paths.PrivateKeyFile,
		"cert_pem", certPEMPath,
		"key_pem", keyPEMPath,
	)
	return nil
}

// writePEM writes bytes as one PEM block of blockType to path.
func writePEM(path, blockType string, bytes []byte, mode os.FileMode) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create %s: %w", filepath.Base(path), err)
	}
	if err := pem.Encode(file, &pem.Block{Type: blockType, Bytes: bytes}); err != nil {
		file.Close()
		return fmt.Errorf("encode %s: %w", filepath.Base(path), err)
	}
	return file.Close()
}
