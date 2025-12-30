package kap

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptPrefix = "$2y$"
)

type Secret []string

func (ss Secret) Has(secret string) bool {
	for _, s := range ss {
		if strings.HasPrefix(s, BcryptPrefix) {
			if bcrypt.CompareHashAndPassword([]byte(s), []byte(secret)) == nil {
				return true
			}
		} else {
			if s == secret {
				return true
			}
		}
	}

	return false
}
