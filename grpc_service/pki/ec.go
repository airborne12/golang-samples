package pki

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// LoadECKeyFromFile loads EC key from unencrypted PEM file.
func LoadECKeyFromFile(fileName string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read signing key file: %v", err)
	}

	privateKeyPEM, _ := pem.Decode(privateKeyBytes)
	if privateKeyPEM == nil {
		return nil, fmt.Errorf("failed to decode pem signing key file: %v", err)
	}

	privateKey, err := x509.ParseECPrivateKey(privateKeyPEM.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signing key file: %v", err)
	}

	return privateKey, nil
}
