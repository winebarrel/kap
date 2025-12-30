package kap

import (
	"errors"
	"slices"
	"strings"

	"github.com/alecthomas/kong"
	"golang.org/x/crypto/bcrypt"
)

var (
	bcryptPrefixs = []string{
		"$2$",
		"$2a$",
		"$2x$",
		"$2y$",
		"$2b$",
	}
)

type Secrets []Secret

func (ss Secrets) Has(secret string) bool {
	for _, s := range ss {
		if s.equal(secret) {
			return true
		}
	}

	return false
}

type Secret struct {
	equal func(string) bool
}

func (s *Secret) Decode(ctx *kong.DecodeContext) error {
	var secret string

	if err := ctx.Scan.PopValueInto("secret", &secret); err != nil {
		return err
	}

	if secret == "" {
		return errors.New("cannot set secret value to empty string")
	}

	if slices.ContainsFunc(bcryptPrefixs, func(prefix string) bool { return strings.HasPrefix(secret, prefix) }) {
		s.equal = func(input string) bool {
			return bcrypt.CompareHashAndPassword([]byte(secret), []byte(input)) == nil
		}
	} else {
		s.equal = func(input string) bool {
			return input == secret
		}
	}

	return nil
}
