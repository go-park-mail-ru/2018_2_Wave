package game

import (
	"Wave/internal/logger"
	"Wave/internal/metrics"
	"Wave/internal/misc"
	"Wave/internal/services/auth/proto"

	"time"
	"net/http"

	"github.com/gorilla/websocket"

	"Wave/internal/application/manager"
	"Wave/internal/application/room"
	"Wave/internal/application/snake"
)

// TODO:: get the value from configuration files
const wsAppTickRate = 500 * time.Millisecond

type Handler struct {
	LG *logger.Logger
	Prof *metrics.Profiler
	AuthManager auth.AuthClient

	wsApp *manager.Manager
	upgrader websocket.Upgrader
}

func NewHandler(LG *logger.Logger, Prof *metrics.Profiler) *Handler{
	return &Handler {
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
		LG: LG,
		Prof: Prof,
	}
}

func (h *Handler) WSHandler(rw http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		panic(err)
	}
	if h.AuthManager == nil {
		panic("empty auth manager")
	}

	go func() {
		var (
			cookie        = misc.GetSessionCookie(r)
			userInfo, err = h.AuthManager.Info(r.Context(), &auth.Cookie{CookieValue: cookie})
			username room.UserID
		)
		if err != nil {
			username = h.wsApp.GetNextUserID()
		} else {
			username = room.UserID(userInfo.Username)
		}

		user := room.NewUser(username, ws)
		user.LG = h.LG
		user.AddToRoom(h.wsApp)
		user.Listen()
	}()
}
