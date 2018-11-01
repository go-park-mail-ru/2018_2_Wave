package session

import (
	"Wave/utiles/walhalla"
	"Wave/app/misc"
	"Wave/app/generated/models"

	"strconv"
	"log"
	"github.com/jmoiron/sqlx"
)

type Model struct {
	db *sqlx.DB
}

func NewModel(ctx *walhalla.Context) *Model {
	return &Model{
		db: ctx.DB,
	}
}

const (
	UserInfoTable = "userinfo"
	UsernameCol   = "username"

	SessionTable = "session"
	CookieCol   = "cookie"
)

func (model *Model) present(tableName string, colName string, target string) (bool, error) {
	var exists string
	model.db.Select(&exists, "SELECT EXISTS (SELECT true FROM " + tableName + " WHERE " + colName + "='" + target + "')")

	fl, errParse := strconv.ParseBool(exists)
	if errParse != nil {
		return false, errParse
	}

	return fl, nil
}

func (model *Model) LogIn(credentials &models.UserCredentials) (cookie string, err error) {
	if isPresent, problem := model.present(UserInfoTable, UsernameCol, *credentials.Username); isPresent && problem == nil {
		var psswd string
		row := model.db.QueryRowx("SELECT password FROM userinfo WHERE username=$1", *credentials.Username)
		err := row.Scan(&psswd)

		if err != nil {
			return "", err
		}

		if psswd == *credentials.Password {
			cookie := misc.GenerateCookie()
			model.db.MustExec("INSERT INTO session(uid, cookie) VALUES((SELECT uid FROM userinfo WHERE username=$1), $2);", *credentials.Username, cookie)

			log.Println("login successful: cookie set")

			return cookie, nil
		} else {
			log.Println("login failed: wrong password")

			return "", nil
		}
	}

	return "", nil
}

func (model *Model) LogOut(cookie string) error {
	model.db.MustExec("DELETE FROM session WHERE cookie=$1", cookie)
	log.Println("logout successful")

	return nil
}