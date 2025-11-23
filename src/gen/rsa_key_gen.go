package gen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// GenerateRSAKeys generates a new RSA key pair and returns PEM-encoded private and public keys.
func GenerateRSAKeys() (privateKeyPEM string, publicKeyPEM string, err error) {
	// Generate a 2048-bit private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// Encode private key to PKCS1 PEM format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}
	privateKeyPEM = string(pem.EncodeToMemory(privBlock))

	// Extract public key
	publicKey := &privateKey.PublicKey

	// Encode public key to PKIX PEM format
	pubDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubDER,
	}
	publicKeyPEM = string(pem.EncodeToMemory(pubBlock))

	if privateKeyPEM == "" || publicKeyPEM == "" {
		return "", "", errors.New("failed to encode RSA keys")
	}

	return privateKeyPEM, publicKeyPEM, nil
}
