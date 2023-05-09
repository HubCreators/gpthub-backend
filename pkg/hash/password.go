package hash

import (
	"crypto/sha1"
	"encoding/base64"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

func (h *SHA1Hasher) Hash(password string) (string, error) {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(password))
	if err != nil {
		return "", err
	}
	sha := base64.URLEncoding.EncodeToString(hasher.Sum([]byte(h.salt)))
	return sha, nil
}
