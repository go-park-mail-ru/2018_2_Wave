package misc

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(psswd string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(psswd), bcrypt.MinCost)

	if err != nil {
		//panic(err)
		return []byte{}
	}

	return hash
}

func PasswordsMatched(stored string, provided string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(provided)); err != nil {
		return false
	}

	return true
}
