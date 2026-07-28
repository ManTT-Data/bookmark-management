package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomStringGenerator_GenerateKey(t *testing.T) {
	testKeyGen := NewKeyGenerator()
	res := testKeyGen.GenerateKey(7)
	assert.Len(t, res, 7)
	for _, c := range res {
		assert.Contains(t, charset, string(c))
	}
}

func TestGenerateRandomString(t *testing.T) {
	res := GenerateRandomString(7)
	assert.Len(t, res, 7)
	for _, c := range res {
		assert.Contains(t, charset, string(c))
	}
}
