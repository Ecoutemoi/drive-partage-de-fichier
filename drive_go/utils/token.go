package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func NewShareToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// URL-safe, pas de + / =
	return base64.RawURLEncoding.EncodeToString(b), nil
}
