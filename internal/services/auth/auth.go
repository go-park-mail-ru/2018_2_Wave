package auth

import (
	psql "Wave/internal/database"
	lg "Wave/internal/logger"
	"Wave/internal/misc"
	auth "Wave/internal/services/auth/proto"

	"golang.org/x/net/context"
)

type AuthManager struct {
	DB psql.DatabaseModel
	LG *lg.Logger
}

func NewAuthManager(lg_ *lg.Logger, db_ *psql.DatabaseModel) *AuthManager {
	return &AuthManager{
		LG: lg_,
		DB: *db_,
	}
}

func (authm *AuthManager) Create(ctx context.Context, credentials *auth.Credentials) (*auth.Cookie, error) {
	if psql.ValidateUname(credentials.Username) && psql.ValidatePassword(credentials.Password) {
		if isPresent, problem := authm.DB.Present(psql.UserInfoTable, psql.UsernameCol, credentials.Username); isPresent && problem == nil {

			authm.LG.Sugar.Infow(
				"signup failed, user already exists",
				"source", "auth.go",
				"who", "Create",
			)

			return &auth.Cookie{CookieValue: ""}, nil
		} else if problem != nil {

			authm.LG.Sugar.Infow(
				"signup succeded",
				"source", "auth.go",
				"who", "Create",
			)

			return &auth.Cookie{CookieValue: ""}, problem
		} else if !isPresent {
			cookie := misc.GenerateCookie()
			hashedPsswd := misc.GeneratePasswordHash(credentials.Password)

			if credentials.Avatar != "" {
				authm.DB.Database.MustExec(`
					INSERT INTO userinfo(username,password,avatar)
					VALUES($1, $2, $3)
				`, credentials.Username, hashedPsswd, credentials.Avatar)
			} else {
				authm.DB.Database.MustExec(`
					INSERT INTO userinfo(username,password)
					VALUES($1, $2)
				`, credentials.Username, hashedPsswd)
			}

			authm.DB.Database.MustExec(`
				INSERT INTO session(uid, cookie)
				VALUES(
					(SELECT uid FROM userinfo WHERE username=$1),
					$2
				)
			`, credentials.Username, cookie)

			authm.DB.LG.Sugar.Infow(
				"signup succeded",
				"source", "auth.go",
				"who", "Create",
			)

			return &auth.Cookie{CookieValue: cookie}, nil
		}
	}

	return &auth.Cookie{CookieValue: ""}, nil
}

func (authm *AuthManager) Delete(ctx context.Context, cookie *auth.Cookie) (*auth.Bool, error) {
	authm.DB.Database.QueryRowx(`
		DELETE
		FROM session
		WHERE cookie=$1;
	`, cookie.CookieValue)

	authm.LG.Sugar.Infow(
		"logout succeded",
		"source", "auth.go",
		"who", "Delete",
	)

	return &auth.Bool{Resp: true}, nil
}

func (authm *AuthManager) Info(ctx context.Context, cookie *auth.Cookie) (*auth.UserInfo, error) {
	username := ""
	if err := authm.DB.Database.QueryRow(`
		SELECT username
		FROM (
			SELECT uid
			FROM session
			WHERE cookie=$1
		) u
		INNER JOIN userinfo USING(uid)
	`, cookie.CookieValue).Scan(&username); err != nil {
		return nil, err
	}
	return &auth.UserInfo{Username: username}, nil
}
