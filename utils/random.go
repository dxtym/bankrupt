package utils

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// generate random number (min, max)
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

// generate random string (n)
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		char := alphabet[rand.Intn(k)]
		sb.WriteByte(char)
	}

	return sb.String()
}

// generate random owner
func RandomOwner() string {
	return RandomString(8)
}

// generate random money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// generate random currency
func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}