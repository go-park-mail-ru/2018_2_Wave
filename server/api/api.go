package api

import (
	psql "Wave/server/database"
	"Wave/utiles/misc"
	"Wave/utiles/models"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type Handler struct {
	DB psql.DatabaseModel
}

func (h *Handler) RegisterHandler(rw http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetMeHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)
	profile, err := h.DB.GetMyProfile(cookie)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := profile.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	return
}

func (h *Handler) EditMeHandler(rw http.ResponseWriter, r *http.Request) {
}

func (h *Handler) GetUserHandler(rw http.ResponseWriter, r *http.Request) {
}

func (h *Handler) GetLeadersHandler(rw http.ResponseWriter, r *http.Request) {
}

func (h *Handler) LoginHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.UserCredentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	cookie, err := h.DB.LogIn(user)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	if cookie == "" {

		fr := models.ForbiddenRequest{
			Reason: "Incorrect password.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(rw, string(payload))

		return
	}

	sessionCookie := misc.MakeSessionCookie(cookie)
	http.SetCookie(rw, sessionCookie)
	rw.WriteHeader(http.StatusOK)

	return
}

func (h *Handler) LogoutHandler(rw http.ResponseWriter, r *http.Request) {
}
