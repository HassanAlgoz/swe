package utils

import (
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789_"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}
	return string(b)
}
