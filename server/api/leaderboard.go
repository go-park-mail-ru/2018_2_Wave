package api

import (
	"Wave/server"
	"Wave/server/types"

	"github.com/valyala/fasthttp"
)

// OnLeaderbordGET - public API
// walhalla: {
// 		URI: 		/users/:offset/:limit,
// 		Method: 	GET,
// 		Data: 		uri,
// 		Target:  	types.Pagination,
// 		Validation: true,
// 		Auth: 		true
// }
func OnLeaderbordGET(ctx *fasthttp.RequestCtx, sv *server.Server, p types.Pagination) {
	data, err := sv.DB.GetTopUsers(p.Offset, p.Limit).MarshalJSON()
	
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write(data)
	}
}
