package utils

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Generate random number (min, max)
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

// Generate random string (n)
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		char := alphabet[rand.Intn(k)]
		sb.WriteByte(char)
	}

	return sb.String()
}

// Generate random owner
func RandomOwner() string {
	return RandomString(8)
}

// Generate random money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// Generate random currency
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}