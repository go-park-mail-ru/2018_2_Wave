package game

import (
	"net/http"
	"time"

	"Wave/internal/application/snake"
	"Wave/internal/application/snake_manager"
	"Wave/internal/database"
	"Wave/internal/logger"
	"Wave/internal/metrics"

	"github.com/gorilla/websocket"
)

// TODO:: get the value from configuration files
const wsAppTickRate = 16 * time.Millisecond

type Handler struct {
	LG   *logger.Logger
	Prof *metrics.Profiler
	DB   *database.DatabaseModel

	cookieToRand map[string]int64
	randToCookie map[int64]string

	wsApp    *manager.Manager
	upgrader websocket.Upgrader
}

func NewHandler(LG *logger.Logger, Prof *metrics.Profiler, db *database.DatabaseModel) *Handler {
	return &Handler{
		wsApp: func() *manager.Manager {
			wsApp := manager.New("", wsAppTickRate, nil, Prof)
			wsApp.CreateLobby("snake", nil, snake.RoomType)
			go wsApp.Start()
			return wsApp
		}(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		cookieToRand: make(map[string]int64),
		randToCookie: make(map[int64]string),
		LG:           LG,
		DB:           db,
		Prof:         Prof,
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
		defer func() {
			if err := recover(); err != nil {
				h.LG.Sugar.Infof("Shit happens, sorry")
			}
		}()
		username := ""
		if err := ws.ReadJSON(&username); err != nil {
			h.LG.Sugar.Infof("WS get user name sheet %s", err)
		}

		u, err := h.wsApp.CreateUser(username, ws)
		if err != nil {
			panic(err)
		}
		u.Start()
	}()
}
