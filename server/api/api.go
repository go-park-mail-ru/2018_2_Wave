package api

import (
	psql "Wave/server/database"
	"Wave/utiles/misc"
	"Wave/utiles/models"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type Handler struct {
	DB psql.DatabaseModel
}

func (h *Handler) SlashHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	rw.WriteHeader(http.StatusOK)

	return
}

func (h *Handler) RegisterHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

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

func (h *Handler) GetMeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	cookie := misc.GetSessionCookie(r)
	if isLogged, errLog := h.DB.IsLoggedIn(cookie); !isLogged || errLog != nil {
		fr := models.ForbiddenRequest{
			Reason: "You are not logged in.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(rw, string(payload))

		return
	}

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
	rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	cookie := misc.GetSessionCookie(r)
	if isLogged, errLog := h.DB.IsLoggedIn(cookie); !isLogged || errLog != nil {
		fr := models.ForbiddenRequest{
			Reason: "You are not logged in.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		return
	}

	editUser := models.UserEdit{
		NewUsername: r.FormValue("newUsername"),
		NewPassword: r.FormValue("newPassword"),
		NewAvatar:   r.FormValue("newAvatar"),
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

func (h *Handler) GetUserHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

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

func (h *Handler) GetLeadersHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	//vars := mux.Vars(r)
	//leaders, err := h.DB.GetTopUsers(strconv.ParseInt(vars["count"]), strconv.ParseInt(vars["page"])
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
}

func (h *Handler) LoginHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

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

func (h *Handler) LogoutHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	cookie := misc.GetSessionCookie(r)
	if isLogged, errLog := h.DB.IsLoggedIn(cookie); !isLogged || errLog != nil {
		fr := models.ForbiddenRequest{
			Reason: "You are not logged in.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(rw, string(payload))

		return
	}
	fmt.Println(cookie)

	if err := h.DB.LogOut(cookie); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.SetCookie(rw, misc.MakeSessionCookie(""))
	rw.WriteHeader(http.StatusOK)

	return
}
