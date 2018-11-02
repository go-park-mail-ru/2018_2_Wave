package user

import (
	"Wave/app/generated/restapi/operations/user"
	"Wave/utiles/walhalla"

	"github.com/go-openapi/runtime/middleware"
)

// walhalla:gen
func Leaders(params user.LeadersParams, ctx *walhalla.Context, model *Model) middleware.Responder {
	leaders, err := model.GetTopUsers(int(*params.Body.Count), int(*params.Body.Page))
	
	if err != nil {
		return user.NewLeadersInternalServerError()
	}

	return user.NewLeadersOK().WithPayload(leaders)
}
