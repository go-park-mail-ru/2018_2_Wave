package misc

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(psswd string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(psswd), bcrypt.DefaultCost)

	return hash
}

func PasswordsMatched(stored []byte, provided string) bool {
	if err := bcrypt.CompareHashAndPassword(stored, []byte(provided)); err != nil {
		return false
	}

	return true
}
