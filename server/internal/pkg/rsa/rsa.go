package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"
)

// GenerateKeyPair generates a 2048-bit RSA key pair.
// Returns: privateKeyPEM (PEM format), publicKeyBase64 (raw base64, Java-compatible), err
func GenerateKeyPair() (privateKeyPEM, publicKeyBase64 string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// Private key: store as PEM for internal use
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})

	// Public key: return as raw base64 (Java-compatible, no PEM headers)
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubBase64 := base64.StdEncoding.EncodeToString(pubBytes)

	return string(privPEM), pubBase64, nil
}

// ExtractPublicKeyBase64 extracts raw base64 public key from PEM format.
func ExtractPublicKeyBase64(publicKeyPEM string) string {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return publicKeyPEM // Already base64, return as-is
	}
	return base64.StdEncoding.EncodeToString(block.Bytes)
}

// Decrypt decrypts a base64-encoded ciphertext using the PEM-encoded private key.
func Decrypt(privateKeyPEM, ciphertext string) (string, error) {
	log.Printf("[DEBUG RSA] Decrypt called, privateKeyPEM len: %d, ciphertext len: %d", len(privateKeyPEM), len(ciphertext))

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		log.Printf("[DEBUG RSA] Failed to decode PEM block")
		return "", errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Printf("[DEBUG RSA] Failed to parse private key: %v", err)
		return "", err
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		log.Printf("[DEBUG RSA] Failed to decode base64: %v", err)
		return "", err
	}

	log.Printf("[DEBUG RSA] Ciphertext bytes len: %d, key size: %d", len(ciphertextBytes), privateKey.N.BitLen())

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertextBytes)
	if err != nil {
		log.Printf("[DEBUG RSA] DecryptPKCS1v15 failed: %v", err)
		return "", err
	}

	log.Printf("[DEBUG RSA] Decryption successful, plaintext: %q", string(plaintext))
	return string(plaintext), nil
}
