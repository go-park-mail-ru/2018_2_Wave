package api

import (
	//lg "Wave/utiles/logger"
	"Wave/utiles/misc"
	"Wave/utiles/models"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func (h *Handler) SlashHandler(rw http.ResponseWriter, r *http.Request) {
	h.DB.Logtest()
	rw.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) RegisterPOSTHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.UserCredentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	log.Println(user.Username)
	log.Println(user.Password)

	cookie, err := h.DB.SignUp(user)
	log.Println(cookie)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}
	if cookie == "" {
		fr := models.ForbiddenRequest{
			Reason: "Username already in use.",
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

func (h *Handler) MeGETHandler(rw http.ResponseWriter, r *http.Request) {
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

func (h *Handler) EditMePUTHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)

	editUser := models.UserEdit{
		Username: r.FormValue("newUsername"),
		Password: r.FormValue("newPassword"),
		//Avatar:   r.FormValue("newAvatar"),
	}

	isUpdated, err := h.DB.UpdateProfile(editUser, cookie)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	if !isUpdated {
		fr := models.ForbiddenRequest{
			Reason: "Incorrect password.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		return
	}

	rw.WriteHeader(http.StatusOK)

	return
}

func (h *Handler) UserGETHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profile, err := h.DB.GetProfile(vars["name"])

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	if reflect.DeepEqual(models.UserExtended{}, profile) {
		rw.WriteHeader(http.StatusNotFound)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := profile.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	return
}

func (h *Handler) LeadersGETHandler(rw http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//leaders, err := h.DB.GetTopUsers(strconv.ParseInt(vars["count"]), strconv.ParseInt(vars["page"])
	/*
		pagination := models.Pagination{
			Page:  r.FormValue("page"),
			Count: r.FormValue("count"),
		}

		c, _ := strconv.Atoi(pagination.Count)
		p, _ := strconv.Atoi(pagination.Page)
		leaders, err := h.DB.GetTopUsers(c, p)

		if err != nil || reflect.DeepEqual(models.Leaders{}, leaders) {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		rw.WriteHeader(http.StatusOK)
		payload, _ := leaders.MarshalJSON()
		fmt.Fprintln(rw, string(payload))

		return
	*/
}

func (h *Handler) LoginPOSTHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.UserCredentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	log.Println(user.Username)
	log.Println(user.Password)

	cookie, err := h.DB.LogIn(user)
	log.Println(cookie)

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

func (h *Handler) LogoutDELETEHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)

	fmt.Println(cookie)

	if err := h.DB.LogOut(cookie); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.SetCookie(rw, misc.MakeSessionCookie(""))
	rw.WriteHeader(http.StatusOK)

	return
}

func (h *Handler) EditMeOPTHandler(rw http.ResponseWriter, r *http.Request) {
}

func (h *Handler) LogoutOPTHandler(rw http.ResponseWriter, r *http.Request) {
}
