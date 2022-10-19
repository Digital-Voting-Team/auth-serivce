package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(input string) string {
	hash := sha256.New()

	var buf []byte
	hash.Write(decodeHex(input))
	buf = hash.Sum(nil)

	return hex.EncodeToString(buf)
}

func decodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return b
}
