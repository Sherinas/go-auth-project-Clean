package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidPassword(t *testing.T) {

	tests := []struct {
		password string
		Expected bool
	}{

		{"Password@123", true},
		{"password@123", false},
		{"PASSWORD@123", false},
		{"Password@aaa", false},
		{"Password123", false},
		{"P@1a", false},
	}

	for _, test := range tests {

		result := isValidPassword(test.password)
		assert.Equal(t, test.Expected, result, "Password: %s", test.password)
	}
}
