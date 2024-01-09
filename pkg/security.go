package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func CalculateHMAC(method, path string, body []byte, timeStamp string) string {
	key := []byte(Config.SecretKey)

	// Create an HMAC hasher using SHA-256 hash function and the secret key
	hasher := hmac.New(sha256.New, key)

	// Write the data to the hasher (method, path, and body)
	hasher.Write([]byte(method))
	hasher.Write([]byte(path))
	hasher.Write([]byte(timeStamp))
	hasher.Write(body)

	// Calculate the HMAC
	hmacResult := hasher.Sum(nil)

	// Convert the HMAC to a hexadecimal string representation
	return hex.EncodeToString(hmacResult)
}
