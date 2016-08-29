package helpers

import (
	"math/rand"
)

// GetRandomKey gets random key from provided map
func GetRandomKey(data map[string]bool) string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	n := rand.Intn(len(keys))
	return keys[n]
}
