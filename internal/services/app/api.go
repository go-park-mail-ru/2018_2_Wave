package api

import (
	psql "Wave/internal/database"
	lg "Wave/internal/logger"
	mc "Wave/internal/metrics"
	"Wave/internal/services/auth/proto"
	"Wave/internal/models"
	"Wave/internal/misc"

	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"os"
	"io"
	//"log"

	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	"golang.org/x/net/context"
)

type Handler struct {
	DB psql.DatabaseModel
	LG *lg.Logger
	Prof *mc.Profiler
	AuthManager auth.AuthClient
}

func (h *Handler) uploadHandler(r *http.Request) (created bool, path string) {
	file, _, err := r.FormFile("avatar")
	defer file.Close()

	if err != nil {

		h.LG.Sugar.Infow("upload failed, unable to read from FormFile, default avatar set",
		"source", "api.go",
		"who", "uploadHandler",)

        return true, "/img/avatars/default"
	}

	prefix := "/img/avatars/"
	hash := ksuid.New()
	fileName := hash.String()

	createPath := ".." + prefix + fileName
	path = prefix + fileName

	out, err := os.Create(createPath)
	defer out.Close()

    if err != nil {

		h.LG.Sugar.Infow("upload failed, file couldn't be created",
		"source", "api.go",
		"who", "uploadHandler",)

        return false, ""
    }

    _, err = io.Copy(out, file)
    if err != nil {

        h.LG.Sugar.Infow("upload failed, couldn't copy data",
		"source", "api.go",
		"who", "uploadHandler",)

		return false, ""
    }

	h.LG.Sugar.Infow("upload succeded",
		"source", "api.go",
		"who", "uploadHandler",)

	return true, path
}

func (h *Handler) SlashHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)

	h.Prof.HitsStats.
	WithLabelValues("200", "OK").
	Add(1)

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
		"who", "RegisterPOSTHandler",)

		h.Prof.HitsStats.
		WithLabelValues("400", "BAD REQUEST").
		Add(1)

		return
	}

	cookie, err := h.AuthManager.Create(
			context.Background(),
			&auth.Credentials{
			Username: user.Username,
			Password: user.Username,
			Avatar: user.Avatar,
		})

	if err != nil {

		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/users failed",
		"source", "api.go",
		"who", "RegisterPOSTHandler",)

		h.Prof.HitsStats.
		WithLabelValues("500", "INTERNAL SERVER ERROR").
		Add(1)

		return
	}
	// or validation failed
	if cookie.CookieValue == "" {
		fr := models.ForbiddenRequest{
			Reason: "Username already in use.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users failed, username already in use.",
		"source", "api.go",
		"who", "RegisterPOSTHandler",)

		h.Prof.HitsStats.
		WithLabelValues("403", "FORBIDDEN").
		Add(1)

		return
	}

	sessionCookie := misc.MakeSessionCookie(cookie.CookieValue)
	http.SetCookie(rw, sessionCookie)
	rw.WriteHeader(http.StatusCreated)

	h.LG.Sugar.Infow("/users succeded",
		"source", "api.go",
		"who", "RegisterPOSTHandler",)

	h.Prof.HitsStats.
	WithLabelValues("201", "CREATED").
	Add(1)

	return
}

func (h *Handler) MeGETHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)

	profile, err := h.DB.GetMyProfile(cookie)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		h.Prof.HitsStats.
		WithLabelValues("500", "INTERNAL SERVER ERROR").
		Add(1)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := profile.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/users/me succeded",
		"source", "api.go",
		"who", "MeGETHandler",)

	h.Prof.HitsStats.
	WithLabelValues("200", "OK").
	Add(1)

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
			Reason: "Update didn't happend, shitty avatar.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users/me failed, bad avatar.",
		"source", "api.go",
		"who", "EditMePUTHandler",)

		h.Prof.HitsStats.
		WithLabelValues("403", "FORBIDDEN").
		Add(1)

		return
	}

	_, err := h.DB.UpdateProfile(user, cookie)

	if err != nil {
		fr := models.ForbiddenRequest{
			Reason: "Update didn't happend, shitty username and/or password.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users/me failed",
		"source", "api.go",
		"who", "EditMePUTHandler",)

		h.Prof.HitsStats.
		WithLabelValues("403", "FORBIDDEN").
		Add(1)

		return
	}

	h.LG.Sugar.Infow("/users/me succeded, user profile updated",
	"source", "api.go",
	"who", "EditMePUTHandler",)

	rw.WriteHeader(http.StatusOK)

	h.Prof.HitsStats.
	WithLabelValues("200", "OK").
	Add(1)

	return
}

func (h *Handler) UserGETHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profile, err := h.DB.GetProfile(vars["name"])

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/users/{name} failed",
		"source", "api.go",
		"who", "UserGETHandler",)

		h.Prof.HitsStats.
		WithLabelValues("500", "INTERNAL SERVER ERROR").
		Add(1)

		return
	}

	if reflect.DeepEqual(models.UserExtended{}, profile) {
		rw.WriteHeader(http.StatusNotFound)

		h.LG.Sugar.Infow("/users/{name} failed",
		"source", "api.go",
		"who", "UserGETHandler",)

		h.Prof.HitsStats.
		WithLabelValues("404", "NOT FOUND").
		Add(1)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := profile.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/users/{name} succeded",
		"source", "api.go",
		"who", "UserGETHandler",)

	h.Prof.HitsStats.
	WithLabelValues("200", "OK").
	Add(1)

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
		"who", "LeadersGETHandler",)

		h.Prof.HitsStats.
		WithLabelValues("500", "INTERNAL SERVER ERROR").
		Add(1)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := leaders.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/users/leaders succeded",
	"source", "api.go",
	"who", "LeadersGETHandler",)

	h.Prof.HitsStats.
	WithLabelValues("200", "OK").
	Add(1)

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
		"who", "LoginPOSTHandler",)

		h.Prof.HitsStats.
		WithLabelValues("500", "INTERNAL SERVER ERROR").
		Add(1)

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
		"who", "LoginPOSTHandler",)

		h.Prof.HitsStats.
		WithLabelValues("401", "UNAUTHORIZED").
		Add(1)

		return
	}

	sessionCookie := misc.MakeSessionCookie(cookie)
	http.SetCookie(rw, sessionCookie)
	rw.WriteHeader(http.StatusOK)

	h.LG.Sugar.Infow("/session succeded",
		"source", "api.go",
		"who", "LoginPOSTHandler",)

	h.Prof.HitsStats.
	WithLabelValues("200", "OK").
	Add(1)

	return
}

func (h *Handler) LogoutDELETEHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)
	success, _ := h.AuthManager.Delete(
		context.Background(),
		&auth.Cookie{
			CookieValue: cookie,
		})

	if success.Resp != true {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/session failed",
		"source", "api.go",
		"who", "LogoutDELETEHandler",)

		h.Prof.HitsStats.
		WithLabelValues("500", "INTERNAL SERVER ERROR").
		Add(1)

		return
	}

	http.SetCookie(rw, misc.MakeSessionCookie(""))
	rw.WriteHeader(http.StatusOK)

	h.LG.Sugar.Infow("/session succeded",
		"source", "api.go",
		"who", "LogoutDELETEHandler",)

	h.Prof.HitsStats.
	WithLabelValues("200", "OK").
	Add(1)

	return
}

func (h *Handler) EditMeOPTHandler(rw http.ResponseWriter, r *http.Request) {

	h.LG.Sugar.Infow("/users/me succeded",
		"source", "api.go",
		"who", "EditMeOPTHandler",)

}

func (h *Handler) LogoutOPTHandler(rw http.ResponseWriter, r *http.Request) {

	h.LG.Sugar.Infow("/session succeded",
		"source", "api.go",
		"who", "LogoutOPTHandler",)

}