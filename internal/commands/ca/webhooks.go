package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"time"

	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/android-sms-gateway/client-go/ca"
	"github.com/urfave/cli/v2"
)

var webhooks = &cli.Command{
	Name:      "webhooks",
	Aliases:   []string{"wh"},
	Usage:     "Issue a new certificate for receiving webhooks to local IP address",
	Args:      true,
	ArgsUsage: "IP address",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "out",
			Usage:    "Certificate output file",
			Required: false,
			Value:    "server.crt",
		},
		&cli.StringFlag{
			Name:     "keyout",
			Usage:    "Private key output file",
			Required: false,
			Value:    "server.key",
		},
	},
	Action: func(c *cli.Context) error {
		ip := c.Args().Get(0)
		if ip == "" {
			return cli.Exit("IP address is empty", codes.ParamsError)
		}

		netipAddr, err := netip.ParseAddr(ip)
		if err != nil {
			return cli.Exit(err.Error(), codes.ParamsError)
		}

		if !netipAddr.IsPrivate() {
			return cli.Exit("IP address is not private", codes.ParamsError)
		}

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

		subject := pkix.Name{
			CommonName: netipAddr.String(),
		}

		csrTemplate := x509.CertificateRequest{
			Subject:            subject,
			SignatureAlgorithm: x509.ECDSAWithSHA256,
			// Key Usage: nonRepudiation, digitalSignature, keyEncipherment
			ExtraExtensions: []pkix.Extension{
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
			},
			IPAddresses: []net.IP{netipAddr.AsSlice()},
		}

		log.Println("Creating certificate request...")
		csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, priv)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Failed to create certificate request: %s", err.Error()), codes.InternalError)
		}

		csrPemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})

		client := metadata.GetCAClient(c.App.Metadata)

		log.Println("Sending certificate request...")
		resp, err := client.PostCSR(c.Context, ca.PostCSRRequest{
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
	},
}
