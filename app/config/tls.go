// Copyright (c) 2023 Fajar Laksono All Rights Reserved.

package config

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

const (
	// DefaultSSLCertPath default certificate path for linux alpine
	DefaultSSLCertPath = "/etc/ssl/certs/ca-certificates.crt"
)

var (
	ErrPrivateKeyInPKCS8Wrapping = errors.New("found unknown private key type in PKCS8 wrapping")
	ErrNoCertOrPrivateKeyFound   = errors.New("no certificate or private key found")
)

// NewTLS creates TLS configuration based on given certificate path.
func NewTLS(certPath string) (*tls.Config, error) {
	certInByte, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	var cert tls.Certificate
	for {
		certBlock, certRest := pem.Decode(certInByte)
		if certBlock == nil {
			break
		}

		if certBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certBlock.Bytes)
		} else {
			cert.PrivateKey, err = parsePrivateKey(certBlock.Bytes)
			if err != nil {
				return nil, err
			}
		}

		certInByte = certRest
	}

	if len(cert.Certificate) < 1 && cert.PrivateKey == nil {
		return nil, ErrNoCertOrPrivateKeyFound
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConf, nil
}

func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}

	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, ErrPrivateKeyInPKCS8Wrapping
		}
	}

	key, err := x509.ParseECPrivateKey(der)
	if err != nil {
		return nil, err
	}

	return key, nil
}
