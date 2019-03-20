package database

import (
	"Wave/internal/models"
	"testing"
)

type RegisterUserCase struct {
	user    models.UserEdit
	isError bool
}

// database-only unit tests

func TestRegisterUser(t *testing.T) {
	cases := []RegisterUserCase{
		RegisterUserCase{models.UserEdit{Username: "Ksenia", Password: "pass123"}, false}, // correct
		RegisterUserCase{models.UserEdit{Username: "Artem", Password: "pswd123"}, false},  // correct
		RegisterUserCase{models.UserEdit{Username: "Ksenia", Password: "newpass"}, true},  // user already exists
		RegisterUserCase{models.UserEdit{Username: "Ks", Password: "randompass"}, true},   // username too short
		RegisterUserCase{models.UserEdit{Username: "Ks", Password: "xz"}, true},           // username and password too short
		RegisterUserCase{models.UserEdit{Username: "RandomUser", Password: "xz"}, true},   // password too short
	}

	db := New()

	for caseNum, item := range cases {
		cookie, err := db.Register(item.user)

		if item.isError == false && err != nil {
			t.Errorf("[%d] unexpected error: %+v", caseNum, err)
		}

		if item.isError == true && err == nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}

		if item.isError == true && err != nil {
			t.Logf("[%d] got error: %+v", caseNum, err)
		}

		if err == nil {
			u, _ := db.GetMyProfile(cookie)
			if u.Username != item.user.Username {
				t.Errorf("[%d] wrong results: got %s, expected %s",
					caseNum, u.Username, item.user.Username)
			}
		}
	}

	db.Database.Exec("TRUNCATE userinfo CASCADE;")
}
