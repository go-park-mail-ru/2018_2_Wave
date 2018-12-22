package game

import (
	"net/http"
	"time"

	"Wave/internal/application/manager"
	"Wave/internal/application/room"
	"Wave/internal/application/snake"
	"Wave/internal/database"
	"Wave/internal/logger"
	"Wave/internal/metrics"
	"Wave/internal/misc"

	"github.com/gorilla/websocket"
)

// TODO:: get the value from configuration files
const wsAppTickRate = 16 * time.Millisecond

type Handler struct {
	LG   *logger.Logger
	Prof *metrics.Profiler
	DB   *database.DatabaseModel

	wsApp    *manager.Manager
	upgrader websocket.Upgrader
}

func NewHandler(LG *logger.Logger, Prof *metrics.Profiler, db *database.DatabaseModel) *Handler {
	return &Handler{
		wsApp: func() *manager.Manager {
			wsApp := manager.New("", wsAppTickRate, nil, Prof)
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
		LG:   LG,
		DB:   db,
		Prof: Prof,
	}
}

func (h *Handler) WSHandler(rw http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		panic(err)
	}
	if h.DB == nil {
		panic("no database")
	}

	go func() {
		var (
			cookie      = misc.GetSessionCookie(r)
			username, _ = h.DB.Info(cookie)
		)
		user := room.NewUser(h.wsApp.GetNextUserID(), ws)
		user.Name = username
		user.LG = h.LG
		user.AddToRoom(h.wsApp)
		user.Listen()
	}()
}
