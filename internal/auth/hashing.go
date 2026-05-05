package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (hash string, err error) {
	hash, err = argon2id.CreateHash(password, argon2id.DefaultParams)
	return
}

func CheckPasswordHash(password, hash string) (match bool, err error) {
	match, err = argon2id.ComparePasswordAndHash(password, hash)
	return
}

