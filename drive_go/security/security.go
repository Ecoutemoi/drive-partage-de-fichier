package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// Hash hash un mot de passe.
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword compare un hash et un mot de passe en clair.
func VerifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Encrypt chiffre des données (AES-GCM). Retourne une base64 string (en bytes).
func Encrypt(plain []byte, key []byte) ([]byte, error) {
	k := sha256.Sum256(key) // clé 32 bytes
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// On préfixe le nonce
	ciphertext := gcm.Seal(nonce, nonce, plain, nil)

	out := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(out, ciphertext)
	return out, nil
}

// Decrypt déchiffre une string base64 produite par Encrypt.
func Decrypt(enc string, key []byte) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}

	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ns := gcm.NonceSize()
	if len(raw) < ns {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := raw[:ns], raw[ns:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
