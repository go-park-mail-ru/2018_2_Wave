package api

import (
	psql "Wave/server/database"
	"Wave/server/room/application"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// TODO:: configs:: room tick rate
const wsAppTickRate = 16 * time.Millisecond

type Handler struct {
	DB       psql.DatabaseModel
	wsApp    *application.Application
	upgrader websocket.Upgrader
	//LG     *lg.Logger
}

func New(model *psql.DatabaseModel) *Handler {
	return &Handler{
		wsApp: func() *application.Application {
			wsApp := application.New("app", wsAppTickRate)
			go wsApp.Run()
			return wsApp
		}(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		DB: *model,
	}
}
