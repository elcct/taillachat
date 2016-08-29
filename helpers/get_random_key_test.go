package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRandomKey(t *testing.T) {
	keys := map[string]bool{
		"one": true,
		"two": true,
	}

	freq := map[string]int{
		"one": 0,
		"two": 0,
	}

	for i := 1; i <= 10000; i++ {
		key := GetRandomKey(keys)
		if _, exists := keys[key]; !exists {
			// That's a bit lame, but if exists is not true, then something is
			// wrong with our function
			assert.Equal(t, false, exists)
		} else {
			freq[key]++
		}
	}
	delta := freq["one"] - freq["two"]
	if delta < 0 {
		delta = -delta
	}
	// Let's assume we accept generous 10% in favour of any key
	assert.Equal(t, true, delta < 1000)
}
