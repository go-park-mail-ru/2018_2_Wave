package auth

import (
	psql "Wave/internal/database"
	lg "Wave/internal/logger"
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

func (authm *AuthManager) Create(ctx context.Context, cs *auth.Credentials) (*auth.Cookie, error) {

}

func (authm *AuthManager) Check(ctx context.Context, cs *auth.Credentials) (*auth.Bool, error) {

}

func (authm *AuthManager) Delete(ctx context.Context, ck *auth.Cookie) (*auth.Bool, error) {

}

/*
func (authm *AuthManager) Create(ctx context.Context, in *session.Session) (*session.SessionID, error) {
		if isPresent, problem := authm.DB.present(UserInfoTable, UsernameCol, in.Login); isPresent && problem == nil {

			authm.LG.Sugar.Infow(
				"signup failed, user already exists",
				"source", "database.go",
				"who", "SignUp",
			)

			return &session.SessionID{ID: ""}, nil
		} else if problem != nil {

			authm.LG.Sugar.Infow(
				"signup succeded",
				"source", "database.go",
				"who", "SignUp",
			)

			return &session.SessionID{}, problem
		} else if !isPresent {
			cookie := misc.GenerateCookie()
			hashedPsswd := misc.GeneratePasswordHash(in.Password)

			if in.Avatar != "" {
				authm.Database.MustExec(`
					INSERT INTO userinfo(username,password,avatar)
					VALUES($1, $2, $3)
				`, in.Login, hashedPsswd, in.Avatar)
			} else {
				authm.Database.MustExec(`
					INSERT INTO userinfo(username,password)
					VALUES($1, $2)
				`, in.Login, hashedPsswd)
			}

			authm.Database.MustExec(`
				INSERT INTO session(uid, cookie)
				VALUES(
					(SELECT uid FROM userinfo WHERE username=$1),
					$2
				)
			`, in.Login, cookie)

			authm.LG.Sugar.Infow(
				"signup succeded",
				"source", "database.go",
				"who", "SignUp",
			)

			return &session.SessionID{ID: cookie}, nil
		}

	return  &session.SessionID{ID: ""}, nil
}

func (authm *AuthManager) Check(ctx context.Context, in *session.SessionID) (*session.Nothing, error) {
	fmt.Println("call Check", in)
	if in.ID == "" {
		return &session.Nothing{Dummy: false}, nil
	}

	return &session.Nothing{Dummy: true}, nil
}

func (authm *AuthManager) Delete(ctx context.Context, in *session.SessionID) (*session.Nothing, error) {
	fmt.Println("call Delete", in)
	authm.Database.QueryRowx(`DELETE FROM session WHERE cookie=$1;`, in.ID)

	authm.LG.Sugar.Infow(
		"logout succeded",
		"source", "database.go",
		"who", "LogOut",
	)

	return &session.Nothing{Dummy: true}, nil
}*/
