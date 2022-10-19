package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(input string) string {
	hash := sha256.New()

	var buf []byte
	hash.Write([]byte(input))
	buf = hash.Sum(nil)

	return hex.EncodeToString(buf)
}
