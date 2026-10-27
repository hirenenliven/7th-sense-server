package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"time"
)

func main() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	uri, _ := url.Parse("urn:go-opcua-client")

	tmpl := x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "GoOPCUAClient", Organization: []string{"MyOrg"}},
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
		panic(err)
	}
	keyDER := x509.MarshalPKCS1PrivateKey(key)

	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	must(os.WriteFile("client_cert.der", der, 0644))
	must(os.WriteFile("client_key.der", keyDER, 0600))

	certOut, err := os.Create("client_cert.pem")
	must(err)
	must(pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der}))
	certOut.Close()

	keyOut, err := os.OpenFile("client_key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	must(err)
	must(pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDER}))
	keyOut.Close()

	fmt.Println("4 files created: client_cert.der, client_cert.pem, client_key.der, client_key.pem")
}
