package utils

import (
	"math/rand"
	"time"
)

var charsets = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString() string {

	r := rand.NewSource(time.Now().UnixNano())
	b := make([]rune, 5)
	for i := range b {
		b[i] = charsets[r.Int63()%int64(len(charsets))]
	}
	return string(b)
}
