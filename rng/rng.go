package rng

import (
	"math/rand"
	"time"
)

var chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateUsername() string {
	firstChar := chars[rand.Intn(len(chars)-10)]
	return string(firstChar) + getRandomCharacters(7)
}

func GeneratePassword() string {
	return getRandomCharacters(12)
}

func getRandomCharacters(num int) string {
	out := make([]rune, num)
	for idx := range out {
		out[idx] = chars[rand.Intn(len(chars))]
	}
	return string(out)
}
