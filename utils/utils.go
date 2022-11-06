package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hint(s string, size int) string {
	if len([]rune(s)) < size {
		return s
	}
	return string([]rune(s)[0:size])
}

func HashString(input string) string {
	hash := sha256.New()

	var buf []byte
	hash.Write([]byte(input))
	buf = hash.Sum(nil)

	return hex.EncodeToString(buf)
}
