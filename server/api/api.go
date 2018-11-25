package api

import (
	psql "Wave/server/database"
	lg "Wave/utiles/logger"
	"net/http"

	//"github.com/gorilla/websocket"

	_ "github.com/lib/pq"
)

type Handler struct {
	DB psql.DatabaseModel
	LG *lg.Logger
}

func (h *Handler) SlashHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	return
}
