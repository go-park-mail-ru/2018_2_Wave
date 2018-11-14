package database

import (
	"Wave/utiles/misc"
	"Wave/utiles/models"
	"log"
	"time"
)

const (
	SessionTable = "session"
	CookieCol    = "cookie"
)

func (model *DatabaseModel) Logtest() {
	model.LG.Sugar.Infow("time", time.Now())
}

func (model *DatabaseModel) LogIn(credentials models.UserCredentials) (cookie string, err error) {
	if isPresent, problem := model.present(UserInfoTable, UsernameCol, credentials.Username); isPresent && problem == nil {
		var psswd string
		row := model.Database.QueryRowx("SELECT password FROM userinfo WHERE username=$1", credentials.Username)
		err := row.Scan(&psswd)

		if err != nil {
			//log.Fatal(err)
			panic(err)
			return "", err
		}

		if misc.PasswordsMatched(psswd, credentials.Password) {
			cookie := misc.GenerateCookie()
			model.Database.MustExec("INSERT INTO session(uid, cookie) VALUES((SELECT uid FROM userinfo WHERE username=$1), $2);", credentials.Username, cookie)

			log.Println("login successful: cookie set")

			return cookie, nil
		} else {
			log.Println("login failed: wrong password")

			return "", nil
		}
	} else {
		//log.Fatal(problem)
		panic(problem)
	}

	return "", nil
}

func (model *DatabaseModel) LogOut(cookie string) error {
	log.Println(cookie)
	model.Database.QueryRowx("DELETE FROM session WHERE cookie=$1;", cookie)
	log.Println("logout successful")

	return nil
}
