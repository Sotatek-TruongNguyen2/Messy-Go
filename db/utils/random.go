package db_utils

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnop"

func RandomInteger(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInteger(50, 200)
}

func RandomCurrency() string {
	allowedCurrency := []string {"USD", "EUR", "VND", "RUP"}
	return allowedCurrency[rand.Intn(len(allowedCurrency))]
}