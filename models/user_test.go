package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserPassword(t *testing.T) {
	user := &User{}

	password := "CatsRuleTheWorld"
	user.HashPassword(password)

	// Our password shouldn't be plain-text
	assert.NotEqual(t, user.Password, password)

	err := user.ComparePassword(password)

	// Hashes should match
	assert.Nil(t, err)

	err = user.ComparePassword("CatsAreNotRulingTheWorld")
	// Hashes should not match
	assert.NotNil(t, err)
}
