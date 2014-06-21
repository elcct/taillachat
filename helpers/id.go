package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

func GetRandomString(n int) string {
	s := make([]byte, n) 
	rand.Read(s)	
	return hex.EncodeToString(s)
}
