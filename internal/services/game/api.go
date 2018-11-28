package game

import (
	"Wave/internal/logger"
	"Wave/internal/metrics"
	"Wave/internal/services/auth/proto"

	"time"
	"net/http"

	"github.com/gorilla/websocket"

	"Wave/application/manager"
	"Wave/application/room"
	"Wave/application/snake"
)

// TODO:: get the value from configuration files
const wsAppTickRate = 500 * time.Millisecond

type Handler struct {
	LG *logger.Logger
	Prof *metrics.Profiler
	AuthManager auth.AuthClient

	wsApp *app.App
	upgrader websocket.Upgrader
}

func NewHandler(LG *logger.Logger, Prof *metrics.Profiler, AuthManager auth.AuthClient) *Handler{
	return &Handler {
		wsApp: func() *app.App {
			wsApp := app.New("app", wsAppTickRate, nil, Prof)
			wsApp.CreateLobby(snake.RoomType, "snake")
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
		LG: LG,
		Prof: Prof,
		AuthManager: AuthManager,
	}
}

func (h *Handler) WSHandler(rw http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		user := room.NewUser(h.wsApp.GetNextUserID(), ws)
		user.LG = h.LG
		user.AddToRoom(h.wsApp)
		user.Listen()
	}()
}
