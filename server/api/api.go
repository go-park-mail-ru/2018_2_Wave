package api

import (
	psql "Wave/server/database"
	"Wave/application/room"
	"Wave/application/manager"
	lg "Wave/utiles/logger"
	"Wave/utiles/misc"
	"Wave/utiles/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
)

// TODO:: get the value from configuration files
const wsAppTickRate = 16 * time.Millisecond

type Handler struct {
	DB       *psql.Model
	wsApp    *app.App
	upgrader websocket.Upgrader
	LG       *lg.Logger
}

func New(model *psql.Model) *Handler {
	return &Handler{
		wsApp: func() *app.App {
			wsApp := app.New("app", wsAppTickRate, model)
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
		DB: model,
	}
}

func (h *Handler) uploadHandler(r *http.Request) (created bool, path string) {
	file, _, err := r.FormFile("avatar")
	defer file.Close()

	if err != nil {

		h.LG.Sugar.Infow("upload failed, not able to read from FormFile, default avatar set",
			"source", "api.go",
			"who", "uploadHandler")

		return true, "" // setting default avatar
	}

	prefix := "/img/avatars/"
	hash := ksuid.New()
	fileName := hash.String()

	createPath := "." + prefix + fileName
	log.Println(fileName)
	path = prefix + fileName

	out, err := os.Create(createPath)
	defer out.Close()

	if err != nil {

		h.LG.Sugar.Infow("upload failed, file couldn't be created",
			"source", "api.go",
			"who", "uploadHandler")

		//file.Close()
		//out.Close()

		return false, ""
	}

	_, err = io.Copy(out, file)
	if err != nil {

		h.LG.Sugar.Infow("upload failed, couldn't copy data",
			"source", "api.go",
			"who", "uploadHandler")

		//file.Close()
		//out.Close()

		return false, ""
	}

	h.LG.Sugar.Infow("upload succeded",
		"source", "api.go",
		"who", "uploadHandler")

	//file.Close()
	//out.Close()

	return true, path
}

func (h *Handler) SlashHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) RegisterPOSTHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.UserEdit{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	isCreated, avatarPath := h.uploadHandler(r)

	if isCreated && avatarPath != "" {
		user.Avatar = avatarPath
	} else if !isCreated {
		fr := models.ForbiddenRequest{
			Reason: "Bad avatar.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users failed, bad avatar.",
			"source", "api.go",
			"who", "RegisterPOSTHandler")

		return
	}

	cookie, err := h.DB.SignUp(user)

	if err != nil {

		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/users failed",
			"source", "api.go",
			"who", "RegisterPOSTHandler")

		return
	}

	if cookie == "" {
		fr := models.ForbiddenRequest{
			Reason: "Username already in use.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users failed, username already in use.",
			"source", "api.go",
			"who", "RegisterPOSTHandler")

		return
	}

	sessionCookie := misc.MakeSessionCookie(cookie)
	http.SetCookie(rw, sessionCookie)
	rw.WriteHeader(http.StatusCreated)

	h.LG.Sugar.Infow("/users succeded",
		"source", "api.go",
		"who", "RegisterPOSTHandler")

	return
}

func (h *Handler) MeGETHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)

	profile, err := h.DB.GetMyProfile(cookie)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/users/me failed",
			"source", "api.go",
			"who", "MeGETHandler")

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := profile.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/users/me succeded",
		"source", "api.go",
		"who", "MeGETHandler")

	return
}

func (h *Handler) EditMePUTHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)

	user := models.UserEdit{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	isCreated, avatarPath := h.uploadHandler(r)

	if isCreated && avatarPath != "" {
		user.Avatar = avatarPath
	} else if !isCreated {
		fr := models.ForbiddenRequest{
			Reason: "Update didn't happend, bad avatar.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users/me failed, bad avatar.",
			"source", "api.go",
			"who", "EditMePUTHandler")

		return
	}

	isUpdated, err := h.DB.UpdateProfile(user, cookie)

	if err != nil {
		fr := models.ForbiddenRequest{
			Reason: "Update didn't happend, something went wrong.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users/me failed",
			"source", "api.go",
			"who", "EditMePUTHandler")

		return
	}

	if !isUpdated {
		fr := models.ForbiddenRequest{
			Reason: "Nothing happened actually.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users/me succeded, nothing changed",
			"source", "api.go",
			"who", "EditMePUTHandler")

		return
	}

	h.LG.Sugar.Infow("/users/me succeded, user profile updated",
		"source", "api.go",
		"who", "EditMePUTHandler")

	rw.WriteHeader(http.StatusOK)

	return
}

func (h *Handler) UserGETHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profile, err := h.DB.GetProfile(vars["name"])

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/users/{name} failed",
			"source", "api.go",
			"who", "UserGETHandler")

		return
	}

	if reflect.DeepEqual(models.UserExtended{}, profile) {
		rw.WriteHeader(http.StatusNotFound)

		h.LG.Sugar.Infow("/users/{name} failed",
			"source", "api.go",
			"who", "UserGETHandler")

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := profile.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/users/{name} succeded",
		"source", "api.go",
		"who", "UserGETHandler")

	return
}

func (h *Handler) LeadersGETHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, _ := strconv.Atoi(vars["count"])
	p, _ := strconv.Atoi(vars["page"])
	leaders, err := h.DB.GetTopUsers(c, p)

	if err != nil || reflect.DeepEqual(models.Leaders{}, leaders) {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/users/leaders failed",
			"source", "api.go",
			"who", "LeadersGETHandler")

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := leaders.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/users/leaders succeded",
		"source", "api.go",
		"who", "LeadersGETHandler")

	return
}

func (h *Handler) LoginPOSTHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.UserCredentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	cookie, err := h.DB.LogIn(user)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/session failed",
			"source", "api.go",
			"who", "LoginPOSTHandler")

		return
	}

	if cookie == "" {

		fr := models.ForbiddenRequest{
			Reason: "Incorrect password/username.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/session failed, incorrect password/username",
			"source", "api.go",
			"who", "LoginPOSTHandler")

		return
	}

	sessionCookie := misc.MakeSessionCookie(cookie)
	http.SetCookie(rw, sessionCookie)
	rw.WriteHeader(http.StatusOK)

	h.LG.Sugar.Infow("/session succeded",
		"source", "api.go",
		"who", "LoginPOSTHandler")

	return
}

func (h *Handler) LogoutDELETEHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)

	if err := h.DB.LogOut(cookie); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/session failed",
			"source", "api.go",
			"who", "LogoutDELETEHandler")

		return
	}

	http.SetCookie(rw, misc.MakeSessionCookie(""))
	rw.WriteHeader(http.StatusOK)

	h.LG.Sugar.Infow("/session succeded",
		"source", "api.go",
		"who", "LogoutDELETEHandler")

	return
}

func (h *Handler) EditMeOPTHandler(rw http.ResponseWriter, r *http.Request) {

	h.LG.Sugar.Infow("/users/me succeded",
		"source", "api.go",
		"who", "EditMeOPTHandler")

}

func (h *Handler) LogoutOPTHandler(rw http.ResponseWriter, r *http.Request) {

	h.LG.Sugar.Infow("/session succeded",
		"source", "api.go",
		"who", "LogoutOPTHandler")

}

// WSHandler - create ws connection
func (h *Handler) WSHandler(rw http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		UID := h.wsApp.GetNextUserID()
		user := room.NewUser(UID, ws)
		user.LG = h.LG
		user.AddToRoom(h.wsApp)
		user.Listen()
	}()
}
