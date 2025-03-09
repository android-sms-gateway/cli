package ca

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"net/netip"

	"github.com/android-sms-gateway/cli/internal/core/codes"
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

		csrTemplate := x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName: netipAddr.String(),
			},
			IPAddresses: []net.IP{netipAddr.AsSlice()},
		}

		return requestCertificate(c, csrTemplate)
	},
}
