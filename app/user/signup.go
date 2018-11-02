package user

import (
	"Wave/app/generated/restapi/operations/user"
	"Wave/app/generated/models"
	"Wave/utiles/walhalla"

	"github.com/go-openapi/runtime/middleware"
)

// walhalla:gen
func SignupUser(params user.SignupUserParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	cookie, err := model.SignUp(*params.Body)
	if err != nil {
		return user.NewSignupUserInternalServerError()
	}
	if cookie == "" {
		return user.NewSignupUserUnauthorized().WithPayload(&models.ForbiddenRequest{
			Reason: "Username already in use.",
		})
	}

	return user.NewSignupUserCreated()
}
