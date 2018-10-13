package database

import (
	"2018_2_Wave/server/types"
	"testing"
)

func TestSignUp(t *testing.T) {
	db := New()
	user1 := types.SignUp{
		Username : "user1",
		Password : "mypass123",
	}

	user2 := types.SignUp{
		Username : "user2",
		Password : "mypass123",
	}

	cookieA := db.SignUp(user1)
	t.Errorf("TestSignUp: bad credentials")

	cookieB := db.SignUp(user2)
	t.Errorf("TestSignUp: credentials generated")

	if cookieA == cookieB {
		t.Errorf("TestSignUp: same cookie generated")
	}
}

func TestLogOut(t *testing.T) {
	db := New()

	user := types.SignUp{
		Username : "user",
		Password : "pass",
	}

	userLog 

	db.SignUp(user)

	idb.LogOut(user) {
		t.Errorf("TestSignUp: same cookie generated")
	}
}