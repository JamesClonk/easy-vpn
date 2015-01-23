package rng

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Rng_GenerateUsername(t *testing.T) {
	username := GenerateUsername()
	if assert.NotNil(t, username) {
		assert.Equal(t, 8, len(username))
	}
}

func Test_Rng_GeneratePassword(t *testing.T) {
	password := GeneratePassword()
	if assert.NotNil(t, password) {
		assert.Equal(t, 12, len(password))
	}
}
