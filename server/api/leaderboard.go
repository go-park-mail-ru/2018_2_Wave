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
	board, err := sv.DB.GetTopUsers(p.Limit, p.Offset)
	
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	} else {
		data, err := board.MarshalJSON()
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Write(data)
	}
}
