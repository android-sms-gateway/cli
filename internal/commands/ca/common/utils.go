package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/android-sms-gateway/client-go/ca"
	"github.com/urfave/cli/v2"
)

func newServerCertificateRequestPEM(template x509.CertificateRequest, priv *ecdsa.PrivateKey) ([]byte, error) {
	template.SignatureAlgorithm = x509.ECDSAWithSHA256
	template.ExtraExtensions = []pkix.Extension{
		{
			Id:       []int{2, 5, 29, 15}, // keyUsage OID
			Critical: true,
			Value:    []byte{0x03, 0x02, 0x05, 0xa0}, // nonRepudiation, digitalSignature, keyEncipherment
		},
		{
			Id:       []int{2, 5, 29, 37}, // extendedKeyUsage OID
			Critical: false,
			Value:    []byte{0x30, 0x06, 0x06, 0x04, 0x2b, 0x06, 0x01, 0x05, 0x05, 0x07, 0x03, 0x01}, // serverAuth
		},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate request: %w", err)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}), nil
}

func requestCertificate(c *cli.Context, typ ca.CSRType, template x509.CertificateRequest) error {
	log.Println("Generating private key...")
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to generate private key: %s", err.Error()), codes.InternalError)
	}

	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to marshal private key: %s", err.Error()), codes.InternalError)
	}

	privPemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	})

	log.Println("Creating certificate request...")
	csrPemBytes, err := newServerCertificateRequestPEM(template, priv)
	if err != nil {
		return cli.Exit(err.Error(), codes.InternalError)
	}

	log.Println("Sending certificate request...")
	client := metadata.GetCAClient(c.App.Metadata)

	resp, err := client.PostCSR(c.Context, ca.PostCSRRequest{
		Type:    typ,
		Content: string(csrPemBytes),
	})
	if err != nil {
		return cli.Exit(err.Error(), codes.ClientError)
	}

	timeout := time.After(30 * time.Second)
	for resp.Certificate == "" {
		select {
		case <-c.Context.Done():
			return cli.Exit("Cancelled", codes.ClientError)
		case <-timeout:
			return cli.Exit("Timeout waiting for certificate", codes.ClientError)
		case <-time.After(1 * time.Second):
		}

		log.Println("Waiting for certificate response...")
		resp, err = client.GetCSRStatus(c.Context, resp.RequestID)
		if err != nil {
			return cli.Exit(err.Error(), codes.ClientError)
		}
	}

	log.Println("Saving certificate...")
	if err := os.WriteFile(c.String("out"), []byte(resp.Certificate), 0644); err != nil {
		return cli.Exit(err.Error(), codes.OutputError)
	}
	if err := os.WriteFile(c.String("keyout"), privPemBytes, 0400); err != nil {
		os.Remove(c.String("out"))
		return cli.Exit(err.Error(), codes.OutputError)
	}

	log.Printf("Certificate saved to %s\n", c.String("out"))
	log.Printf("Private key saved to %s\n", c.String("keyout"))

	return nil
}
