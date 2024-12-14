package random

import (
	"math/rand"
	"time"
)

func NewRandomString(aliasLength int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRTSUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789",
	)

	b := make([]rune, aliasLength)

	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
