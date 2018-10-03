package api

import (
	"Wave/server"
	"strconv"

	"github.com/valyala/fasthttp"
)

// OnLeaderbordGET - public API
// walhalla: {
// 		URI: 		/users/:start/:count,
// 		Method: 	GET,
// 		Auth: 		yes
// }
func OnLeaderbordGET(ctx *fasthttp.RequestCtx, sv *server.Server) {
	var (
		start, err1 = strconv.Atoi(ctx.UserValue("start").(string))
		count, err2 = strconv.Atoi(ctx.UserValue("count").(string))
	)
	if err1 != nil || err2 != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if data, err := sv.DB.GetTopUsers(start, count).MarshalJSON(); err == nil {
		ctx.Write(data)
	} else {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}
