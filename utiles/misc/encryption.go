package misc

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(psswd string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(psswd), bcrypt.MinCost)

	if err != nil {
		//log.Fatal(err)
		panic(err)
		return []byte{}
	}

	return hash
}

func PasswordsMatched(stored string, provided string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(provided)); err != nil {
		log.Println(err)
		return false
	}

	log.Println("encryption: passwords matched")
	return true
}
