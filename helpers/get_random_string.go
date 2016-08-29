package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

// GetRandomString generates random string of length n
func GetRandomString(n int) string {
	s := make([]byte, n)
	rand.Read(s)
	return hex.EncodeToString(s)[:n]
}
