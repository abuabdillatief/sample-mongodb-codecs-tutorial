package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// Encrypt encrypts PII data using AES-GCM.
func Encrypt(plainText, key string) (string, error) {
	// Convert the key to a 32-byte slice
	keyBytes := []byte(key)
	switch len(keyBytes) {
	case 16, 24, 32:
	default:
		return "", errors.New("encryption key must be 16, 24, or 32 bytes long")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Generate a nonce (12 bytes for GCM)
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	cipherText := aesGCM.Seal(nil, nonce, []byte(plainText), nil)

	// Combine nonce and ciphertext for storage
	result := append(nonce, cipherText...)

	// Return base64 encoded result
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts the encrypted data back to plaintext.
func Decrypt(cipherTextBase64, key string) (string, error) {
	// Convert the key to a 32-byte slice
	keyBytes := []byte(key)
	switch len(keyBytes) {
	case 16, 24, 32:
	default:
		return "", errors.New("encryption key must be 16, 24, or 32 bytes long")
	}

	// Decode the base64-encoded ciphertext
	data, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Extract nonce and ciphertext
	if len(data) < 12 {
		return "", errors.New("invalid ciphertext")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM mode: %w", err)
	}
	
	nonce, cipherText := data[:12], data[12:]

	// Decrypt the ciphertext
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plainText), nil
}