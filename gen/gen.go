package gen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/blake2b"
)

// Fingerprint creates a BLAKE2b hash composed of appName and the current timestamp
func Fingerprint(appName string) (string, error) {
	// Get the current timestamp
	timestamp := time.Now().UTC().UnixNano()
	// Concatenate appName and timestamp to create the input for hashing
	input := fmt.Sprintf("%s:%d", appName, timestamp)

	// Generate a random key for the BLAKE2b hash
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}

	// Create the BLAKE2b hash
	hash, err := blake2b.New256(key)
	if err != nil {
		return "", fmt.Errorf("failed to create BLAKE2b hash: %w", err)
	}

	// Write the input to the hash
	hash.Write([]byte(input))

	// Get the resulting hash
	result := hash.Sum(nil)

	// Encode the hash to a hexadecimal string
	secret := hex.EncodeToString(result)

	return secret, nil
}
