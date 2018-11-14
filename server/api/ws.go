package api

import (
	"Wave/server/room"
	"net/http"
)

// WSHandler - create ws connection
func (h *Handler) WSHandler(rw http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		UID := h.wsApp.GetNextUserID()
		user := room.NewUser(UID, ws)
		user.AddToRoom(h.wsApp)
		user.Listen()
	}()
}
