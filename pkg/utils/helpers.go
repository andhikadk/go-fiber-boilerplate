package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomString(byteLength int) string {
	if byteLength <= 0 {
		byteLength = 32
	}
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
