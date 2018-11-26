package implementation

import (
	"fmt"
	"strconv"
	"Wave/session"

	"golang.org/x/net/context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"Wave/utiles/misc"
	//"Wave/utiles/models"
	lg "Wave/utiles/logger"
)

const sessKeyLen = 10

const (
	UserInfoTable = "userinfo"
	UsernameCol   = "username"

	SessionTable = "session"
	CookieCol    = "cookie"
)

type SessionManager struct {
	LG *lg.Logger
	Database *sqlx.DB
}

func NewSessionManager(lg_ *lg.Logger) *SessionManager {
	postgr, _ := sqlx.Connect("postgres", "user=waveapp password='surf' dbname='wave' sslmode=disable")
	return &SessionManager{
		Database: postgr,
		LG: lg_,
	}
}

func (sm *SessionManager) present(tableName string, colName string, target string) (fl bool, err error) {
	var exists string
	row := sm.Database.QueryRowx("SELECT EXISTS (SELECT true FROM " + tableName + " WHERE " + colName + "='" + target + "');")
	err = row.Scan(&exists)

	if err != nil {

		sm.LG.Sugar.Infow(
			"Scan failed",
			"source", "database.go",
			"who", "present",
		)

		return false, err
	}

	fl, err = strconv.ParseBool(exists)

	if err != nil {

		sm.LG.Sugar.Infow(
			"strconv.ParseBool failed",
			"source", "database.go",
			"who", "present",
		)

		return false, err
	}

	return fl, nil
}

func (sm *SessionManager) Create(ctx context.Context, in *session.Session) (*session.SessionID, error) {
		if isPresent, problem := sm.present(UserInfoTable, UsernameCol, in.Login); isPresent && problem == nil {
			
			sm.LG.Sugar.Infow(
				"signup failed, user already exists",
				"source", "database.go",
				"who", "SignUp",
			)

			return &session.SessionID{ID: ""}, nil
		} else if problem != nil {

			sm.LG.Sugar.Infow(
				"signup succeded",
				"source", "database.go",
				"who", "SignUp",
			)

			return &session.SessionID{}, problem
		} else if !isPresent {
			cookie := misc.GenerateCookie()
			hashedPsswd := misc.GeneratePasswordHash(in.Password)

			if in.Avatar != "" {
				sm.Database.MustExec(`
					INSERT INTO userinfo(username,password,avatar)
					VALUES($1, $2, $3)
				`, in.Login, hashedPsswd, in.Avatar)
			} else {
				sm.Database.MustExec(`
					INSERT INTO userinfo(username,password)
					VALUES($1, $2)
				`, in.Login, hashedPsswd)
			}

			sm.Database.MustExec(`
				INSERT INTO session(uid, cookie)
				VALUES(
					(SELECT uid FROM userinfo WHERE username=$1),
					$2
				)
			`, in.Login, cookie)

			sm.LG.Sugar.Infow(
				"signup succeded",
				"source", "database.go",
				"who", "SignUp",
			)

			return &session.SessionID{ID: cookie}, nil
		}

	return  &session.SessionID{ID: ""}, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *session.SessionID) (*session.Nothing, error) {
	fmt.Println("call Check", in)
	if in.ID == "" {
		return &session.Nothing{Dummy: false}, nil
	}

	return &session.Nothing{Dummy: true}, nil
}

func (sm *SessionManager) Delete(ctx context.Context, in *session.SessionID) (*session.Nothing, error) {
	fmt.Println("call Delete", in)
	sm.Database.QueryRowx(`DELETE FROM session WHERE cookie=$1;`, in.ID)

	sm.LG.Sugar.Infow(
		"logout succeded",
		"source", "database.go",
		"who", "LogOut",
	)

	return &session.Nothing{Dummy: true}, nil
}