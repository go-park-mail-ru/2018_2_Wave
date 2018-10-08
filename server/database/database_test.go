package database

import (
	"Wave/server/types"
	"testing"
)

func TestSignUp(t *testing.T) {
	db := New()

	user1A := types.User{
		Username: "A",
		Password: "kek",
	}
	user1B := types.User{
		Username: "A",
	}

	if cookie := db.SignUp(user1A); cookie == "" {
		t.Errorf("Empty cookie")
	}
	if !db.IsSignedUp(user1B) {
		t.Errorf("Registration failed")
	}
}

func TestLogIn(t *testing.T) {
	db := New()

	user1A := types.User{
		Username: "A",
		Password: "kek",
	}
	user1B := types.User{
		Username: "A",
		Password: "lol",
	}

	cookie1 := db.SignUp(user1A)
	cookie2 := db.LogIn(user1A)
	cookie3 := db.LogIn(user1B)

	if cookie2 == "" {
		t.Errorf("Empty cookie from valid user")
	}

	if cookie3 != "" {
		t.Errorf("User with incorrect password mustn't be logged in")
	}

	if cookie1 == cookie2 {
		t.Errorf("Cookies for each session must be unique")
	}
}
