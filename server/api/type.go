package api

import (
	psql "Wave/server/database"
	"Wave/server/room/app"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// TODO:: configs:: room tick rate
const wsAppTickRate = 16 * time.Millisecond

type Handler struct {
	DB       psql.DatabaseModel
	wsApp    *app.App
	upgrader websocket.Upgrader
	//LG     *lg.Logger
}

func New(model *psql.DatabaseModel) *Handler {
	return &Handler{
		wsApp: func() *app.App {
			wsApp := app.New("app", wsAppTickRate)
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
