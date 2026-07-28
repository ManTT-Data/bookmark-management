package utils

import (
	"bytes"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// KeyGenerator represents the interface to generate random string keys
//
//go:generate mockery --name KeyGenerator --filename string.go
type KeyGenerator interface {
	GenerateKey(length int) string
}

type randomStringGenerator struct {
	rng *rand.Rand
}

// NewKeyGenerator returns a KeyGenerator
func NewKeyGenerator() KeyGenerator {
	return &randomStringGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateKey generates a random string with no special characters
func (r *randomStringGenerator) GenerateKey(length int) string {
	return randomString(r.rng, length)
}

func GenerateRandomString(length int) string {
	// Seed the random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	return randomString(rng, length)
}

func randomString(rng *rand.Rand, length int) string {
	var strBuilder bytes.Buffer
	for i := 0; i < length; i++ {
		strBuilder.WriteByte(charset[rng.Intn(len(charset))])
	}

	return strBuilder.String()
}
