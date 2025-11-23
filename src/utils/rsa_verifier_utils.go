package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

func VerifyKeyPair(privateB64, publicB64 string) error {
	privBytes, err := base64.StdEncoding.DecodeString(privateB64)
	if err != nil {
		return err
	}
	pubBytes, err := base64.StdEncoding.DecodeString(publicB64)
	if err != nil {
		return err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privBytes)
	if err != nil {
		return err
	}
	publicKey, err := x509.ParsePKCS1PublicKey(pubBytes)
	if err != nil {
		return err
	}

	// Generate random message
	message := []byte("test-message")
	hash := sha256.Sum256(message)

	// Sign with private key
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, hash[:])
	if err != nil {
		return err
	}

	// Verify with public key
	err = rsa.VerifyPKCS1v15(publicKey, 0, hash[:], signature)
	if err != nil {
		return errors.New("public/private key mismatch")
	}

	return nil
}
