package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	lettersNumber = []rune("0123456789")
	lettersAlpha  = []rune("abcdefghijklmnopqrstuvwxyz")
	letters       = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
)

func GenNumberSecret(n int) string {
	return genSecret(lettersNumber, n)
}

func GenSecret(n int) string {
	return genSecret(letters, n)
}

func genSecret(letters []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
