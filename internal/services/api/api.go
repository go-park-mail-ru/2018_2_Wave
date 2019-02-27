package api

import (
	psql "Wave/internal/database"
	lg "Wave/internal/logger"
	mc "Wave/internal/metrics"
	"Wave/internal/misc"
	"Wave/internal/models"

	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type Handler struct {
	DB   psql.DatabaseModel
	LG   *lg.Logger
	Prof *mc.Profiler
	//AuthManager auth.AuthClient

}

func NewHandler(DB *psql.DatabaseModel, LG *lg.Logger, Prof *mc.Profiler) *Handler {
	h := &Handler{
		DB:   *DB,
		LG:   LG,
		Prof: Prof,
	}
	return h
}

func (h *Handler) uploadHandler(r *http.Request) (created bool, path string) {
	file, _, err := r.FormFile("avatar")

	if err != nil || file == nil {

		h.LG.Sugar.Infow("upload failed, unable to read from FormFile or avatar not provided, default avatar set",
			"source", "api.go",
			"who", "uploadHandler")

		return true, "img/avatars/default"
	}

	defer file.Close()

	prefix := "img/avatars/"
	hash := ksuid.New()
	fileName := hash.String()

	createPath := "./" + prefix + fileName
	path = prefix + fileName

	out, err := os.Create(createPath)
	defer out.Close()

	if err != nil {

		h.LG.Sugar.Infow("upload failed, file couldn't be created",
			"source", "api.go",
			"who", "uploadHandler")

		return false, ""
	}

	_, err = io.Copy(out, file)
	if err != nil {

		h.LG.Sugar.Infow("upload failed, couldn't copy data",
			"source", "api.go",
			"who", "uploadHandler")

		return false, ""
	}

	h.LG.Sugar.Infow("upload succeeded",
		"source", "api.go",
		"who", "uploadHandler")

	return true, path
}

func (h *Handler) SlashHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

/******************** Register POST ********************/

func (h *Handler) RegisterPOSTHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.UserEdit{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	var created bool
	created, user.Avatar = h.uploadHandler(r)

	if !created {
		fr := models.ForbiddenRequest{
			Reason: "Bad avatar",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(rw, string(payload))
		return
	}

	cookie, err := h.DB.Register(user)

	if err != nil {
		if err.Error() == "non-valid" {
			fr := models.ForbiddenRequest{
				Reason: "Bad username or/and password",
			}

			payload, _ := fr.MarshalJSON()
			rw.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(rw, string(payload))

			h.LG.Sugar.Infow("/users failed, bad username or/and password.",
				"source", "api.go",
				"who", "RegisterPOSTHandler")

			h.Prof.HitsStats.
				WithLabelValues("403", "FORBIDDEN").
				Add(1)

			return
		}

		if err.Error() == "exists" {
			fr := models.ForbiddenRequest{
				Reason: "Username already in use.",
			}

			payload, _ := fr.MarshalJSON()
			rw.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(rw, string(payload))

			h.LG.Sugar.Infow("/users failed, username already in use.",
				"source", "api.go",
				"who", "RegisterPOSTHandler")

			h.Prof.HitsStats.
				WithLabelValues("403", "FORBIDDEN").
				Add(1)

			return
		}
	}

	sessionCookie := misc.MakeSessionCookie(cookie)
	http.SetCookie(rw, sessionCookie)
	rw.WriteHeader(http.StatusCreated)

	h.LG.Sugar.Infow("/users succeeded",
		"source", "api.go",
		"who", "RegisterPOSTHandler")

	h.Prof.HitsStats.
		WithLabelValues("201", "CREATED").
		Add(1)

	return
}

/******************** Me GET ********************/

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

	h.LG.Sugar.Infow("/users/me succeeded",
		"source", "api.go",
		"who", "MeGETHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

//
/******************** Edit PUT ********************/

func (h *Handler) EditMePUTHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)

	user := models.UserEdit{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	isCreated, avatarPath := h.uploadHandler(r)

	if isCreated && avatarPath != "/img/avatars/default" {
		user.Avatar = avatarPath
	} else if !isCreated {
		fr := models.ForbiddenRequest{
			Reason: "Bad avatar.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users/me failed, bad avatar.",
			"source", "api.go",
			"who", "EditMePUTHandler")

		h.Prof.HitsStats.
			WithLabelValues("403", "FORBIDDEN").
			Add(1)

		return
	}

	log.Println(user.Avatar)
	err := h.DB.UpdateProfile(user, cookie)

	if err != nil {
		fr := models.ForbiddenRequest{
			Reason: "Bad new username or/and password.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users/me failed",
			"source", "api.go",
			"who", "EditMePUTHandler")

		h.Prof.HitsStats.
			WithLabelValues("403", "FORBIDDEN").
			Add(1)

		return
	}

	h.LG.Sugar.Infow("/users/me succeeded, user profile updated",
		"source", "api.go",
		"who", "EditMePUTHandler")

	rw.WriteHeader(http.StatusOK)

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

/******************** User {name} GET ********************/

func (h *Handler) UserGETHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profile, err := h.DB.GetProfile(vars["name"])

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/users/{name} failed",
			"source", "api.go",
			"who", "UserGETHandler")

		h.Prof.HitsStats.
			WithLabelValues("500", "INTERNAL SERVER ERROR").
			Add(1)

		return
	}

	if reflect.DeepEqual(models.UserExtended{}, profile) {
		rw.WriteHeader(http.StatusNotFound)

		h.LG.Sugar.Infow("/users/{name} failed",
			"source", "api.go",
			"who", "UserGETHandler")

		h.Prof.HitsStats.
			WithLabelValues("404", "NOT FOUND").
			Add(1)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := profile.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/users/{name} succeeded",
		"source", "api.go",
		"who", "UserGETHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

/******************** Login POST ********************/

func (h *Handler) LoginPOSTHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.UserCredentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	cookie, err := h.DB.Login(user)

	if err != nil {
		rw.WriteHeader(http.StatusForbidden)
		// fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/session failed",
			"source", "api.go",
			"who", "LoginPOSTHandler")

		h.Prof.HitsStats.
			WithLabelValues("500", "INTERNAL SERVER ERROR").
			Add(1)

		return
	}

	if cookie == "" {
		fr := models.ForbiddenRequest{
			Reason: "Incorrect password or/and username.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/session failed, incorrect password or/and username",
			"source", "api.go",
			"who", "LoginPOSTHandler")

		h.Prof.HitsStats.
			WithLabelValues("401", "UNAUTHORIZED").
			Add(1)

		return
	}

	sessionCookie := misc.MakeSessionCookie(cookie)
	http.SetCookie(rw, sessionCookie)
	rw.WriteHeader(http.StatusOK)

	h.LG.Sugar.Infow("/session succeeded",
		"source", "api.go",
		"who", "LoginPOSTHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

/******************** Logout DELETE ********************/

func (h *Handler) LogoutDELETEHandler(rw http.ResponseWriter, r *http.Request) {
	/*
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
				"who", "LogoutDELETEHandler")

			h.Prof.HitsStats.
				WithLabelValues("500", "INTERNAL SERVER ERROR").
				Add(1)

			return
		}

		http.SetCookie(rw, misc.MakeSessionCookie(""))
		rw.WriteHeader(http.StatusOK)

		h.LG.Sugar.Infow("/session succeeded",
			"source", "api.go",
			"who", "LogoutDELETEHandler")

		h.Prof.HitsStats.
			WithLabelValues("200", "OK").
			Add(1)

		return
	*/

	cookie := misc.GetSessionCookie(r)
	success := h.DB.Logout(cookie)

	if success != true {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/session failed",
			"source", "api.go",
			"who", "LogoutDELETEHandler")

		h.Prof.HitsStats.
			WithLabelValues("500", "INTERNAL SERVER ERROR").
			Add(1)

		return
	}

	http.SetCookie(rw, misc.MakeSessionCookie(""))
	rw.WriteHeader(http.StatusOK)

	h.LG.Sugar.Infow("/session succeeded",
		"source", "api.go",
		"who", "LogoutDELETEHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

/******************** OPTIONS ********************/

func (h *Handler) EditMeOPTHandler(rw http.ResponseWriter, r *http.Request) {

	h.LG.Sugar.Infow("/users/me succeeded",
		"source", "api.go",
		"who", "EditMeOPTHandler")

}

func (h *Handler) LogoutOPTHandler(rw http.ResponseWriter, r *http.Request) {

	h.LG.Sugar.Infow("/session succeeded",
		"source", "api.go",
		"who", "LogoutOPTHandler")

}

/******************** Applications ********************/

func (h *Handler) ShowAppsGETHandler(rw http.ResponseWriter, r *http.Request) {
	var apps models.Applications
	apps = h.DB.GetApps()

	payload, _ := apps.MarshalJSON()
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/apps succeeded",
		"source", "api.go",
		"who", "ShowAppsGETHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

func (h *Handler) ShowAppsPopularGETHandler(rw http.ResponseWriter, r *http.Request) {
	var apps models.Applications
	apps = h.DB.GetPopularApps()

	payload, _ := apps.MarshalJSON()
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/apps succeeded",
		"source", "api.go",
		"who", "ShowAppsGETHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

func (h *Handler) AppGETHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := h.DB.GetApp(vars["name"])

	if reflect.DeepEqual(models.Application{}, app) {
		rw.WriteHeader(http.StatusNotFound)

		h.LG.Sugar.Infow("/apps/{name} failed",
			"source", "api.go",
			"who", "AppGETHandler")

		h.Prof.HitsStats.
			WithLabelValues("404", "NOT FOUND").
			Add(1)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := app.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/apps/{name} succeeded",
		"source", "api.go",
		"who", "AppGETHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

func (h *Handler) AppPersonalGETHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)
	vars := mux.Vars(r)
	app := h.DB.GetAppPersonal(cookie, vars["name"])

	if reflect.DeepEqual(models.UserApplicationInstalled{}, app) {
		rw.WriteHeader(http.StatusNotFound)

		h.LG.Sugar.Infow("/me/apps/{name} failed",
			"source", "api.go",
			"who", "AppGETHandler")

		h.Prof.HitsStats.
			WithLabelValues("404", "NOT FOUND").
			Add(1)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := app.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/me/apps/{name} succeeded",
		"source", "api.go",
		"who", "AppPersonalGETHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

func (h *Handler) AddAppPOSTHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)
	appname := r.FormValue("name")

	h.DB.AddApp(cookie, appname)
	rw.WriteHeader(http.StatusOK)

	h.LG.Sugar.Infow("/apps succeeded",
		"source", "api.go",
		"who", "AddAppPOSTHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

func (h *Handler) MeShowAppsGetHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)

	var apps models.Applications
	apps = h.DB.GetMyApps(cookie)

	payload, _ := apps.MarshalJSON()
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/me/apps succeeded",
		"source", "api.go",
		"who", "MeShowAppsGetHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}

func (h *Handler) CategoryGETHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var apps models.Applications
	apps = h.DB.GetAppsByCattegory(vars["category"])

	payload, _ := apps.MarshalJSON()
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/apps succeeded",
		"source", "api.go",
		"who", "ShowAppsGETHandler")

	h.Prof.HitsStats.
		WithLabelValues("200", "OK").
		Add(1)

	return
}
