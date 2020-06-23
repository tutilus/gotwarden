package util

import (
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

// HashPassword provide a base64 password with 2 step hashing
func HashPassword(password, salt []byte, iteration int) string {

	// First step: get a key with the data provided
	key := pbkdf2.Key(password, salt, iteration, 32, sha256.New)

	// Second step: Get the hash (password as salt and only one iteration) from the key this time
	hash := pbkdf2.Key(key, password, 1, 32, sha256.New)

	return base64.StdEncoding.EncodeToString(hash)
}
