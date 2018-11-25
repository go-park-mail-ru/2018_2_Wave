package api

import (
	psql "Wave/server/database"
	lg "Wave/utiles/logger"
	"net/http"
	"time"
	"Wave/server/chat"
	"Wave/server/app"
	"Wave/server/room"
	"github.com/gorilla/websocket"
)

// TODO:: get the value from configuration files
const wsAppTickRate = 16 * time.Millisecond

type Handler struct {
	DB       psql.DatabaseModel
	wsApp    *app.App
	upgrader websocket.Upgrader
	LG       *lg.Logger
}

func New(model *psql.DatabaseModel) *Handler {
	return &Handler{
		wsApp: func() *app.App {
			wsApp := app.New("manager", wsAppTickRate, model)
			_, err := wsApp.CreateLobby(chat.RoomType, "global")
			if err != nil {
				panic(err)
			}
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

func (h *Handler) ChatHandler(rw http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		UID := h.wsApp.GetNextUserID()
		user := room.NewUser(UID, ws)
		user.LG = h.LG
		user.ID = "test id"
		user.AddToRoom(h.wsApp)
		user.Listen()
	}()
	return
}
