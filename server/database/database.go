package database

import (
	lg "Wave/utiles/logger"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseModel struct {
	Database *sqlx.DB
	LG       *lg.Logger
}

func New(lg_ *lg.Logger) *DatabaseModel {
	postgr := &DatabaseModel{
		LG: lg_,
	}

	var err error

	dbuser := os.Getenv("WAVE_DB_USER")
	dbpassword := os.Getenv("WAVE_DB_PASSWORD")
	dbname := os.Getenv("WAVE_DB_NAME")

	postgr.Database, err = sqlx.Connect("postgres", "user="+dbuser+" password="+dbpassword+" dbname='"+dbname+"' "+"sslmode=disable")

	if err != nil {
		postgr.LG.Sugar.Panicw(
			"PostgreSQL connection establishment failed",
			"source", "database.go",
			"who", "New",
		)
		panic(err)
	}

	postgr.LG.Sugar.Infow(
		"PostgreSQL connection establishment succeded",
		"source", "database.go",
		"who", "New",
	)

	return postgr
}

const (
	UserInfoTable = "userinfo"
	UsernameCol   = "username"

	SessionTable = "session"
	CookieCol    = "cookie"
)

func (model *DatabaseModel) present(tableName string, colName string, target string) (fl bool, err error) {
	var exists string
	row := model.Database.QueryRowx("SELECT EXISTS (SELECT true FROM " + tableName + " WHERE " + colName + "='" + target + "');")
	err = row.Scan(&exists)

	if err != nil {

		model.LG.Sugar.Infow(
			"Scan failed",
			"source", "database.go",
			"who", "present",
		)

		return false, err
	}

	fl, err = strconv.ParseBool(exists)

	if err != nil {

		model.LG.Sugar.Infow(
			"strconv.ParseBool failed",
			"source", "database.go",
			"who", "present",
		)

		return false, err
	}

	return fl, nil
}

func (model *DatabaseModel) AddMsg(roomid, msg string) error {
	model.Database.MustExec("INSERT INTO usermessage(roomid, msg) VALUES($1, $2)", roomid, msg)
	return nil
}

func (model *DatabaseModel) GetSenderId(cookie string) cookie {
	//`SELECT DISTINCT username FROM userinfo JOIN session WHERE username.uid=session.uid`
	return ""
}
