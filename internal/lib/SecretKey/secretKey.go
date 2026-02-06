package secretkey

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("file is not found")
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("invalid private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)

	if !ok {
		return nil, err
	}

	return rsaKey, nil
}
