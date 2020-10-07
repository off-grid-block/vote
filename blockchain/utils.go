package blockchain

import (
	"time"
	"math/rand"
	"strings"
)

const (
	// length of generated salt
	length = 16
)

// https://yourbasic.org/golang/generate-random-string/
func GenerateRandomSalt() string {

	// choose random seed using current time in nanoseconds	
	rand.Seed(time.Now().UnixNano())

	// set of possible characters to be used in salt
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}
