package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRandomString(t *testing.T) {
	l := 32

	keys := map[string]bool{}

	for i := 1; i <= 100; i++ {
		key := GetRandomString(l)
		assert.Equal(t, l, len(key))
		if _, exists := keys[key]; exists {
			// That's a bit lame, but if exists is true, our random string
			// is definitely not random ;)
			assert.Equal(t, true, exists)
		} else {
			keys[key] = true
		}
	}
}
