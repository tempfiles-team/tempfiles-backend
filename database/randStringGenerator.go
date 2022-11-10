package database

import (
	"math/rand"
	"time"
)

func randInit() {
	rand.Seed(time.Now().UnixNano())
}

var charsets = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = charsets[rand.Intn(len(charsets))]
	}
	return string(b)
}
