package api

import (
	psql "Wave/server/database"
	//lg "Wave/utiles/logger"
	"Wave/utiles/misc"
	"Wave/utiles/models"
	"fmt"
	"log"
	"net/http"
	"reflect"
	//"strconv"
	"time"
	//"encoding/json"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	_ "github.com/lib/pq"
)

type Handler struct {
	DB psql.DatabaseModel
	//LG *lg.Logger
}

func (h *Handler) SlashHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.UserCredentials{
		Username: "ebana",
		Password: "pizdec",
	}
	cookie, _ := h.DB.LogIn(user)
	log.Println(cookie)
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

/************************* websocket block ************************************/

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	waitTime = 15 * time.Second
)

//action_uid uniqely generated on the front
//action_id : 
// 1 - add user to the room
// 2 - remove user from the room
// 3 - start
// 4 - rollback

type lobbyReq struct {
	actionID 	string `json:"action_id"`
	actionUID 	string `json:"action_uid"`
	username 	string `json:"username"`
}

type lobbyRespGenereic struct {
	actionUID string `json:"action_id"`
	status 	 string `json:"status"`
}

func contains(sl []string, str string) bool {
    for _, cur := range sl {
        if str == cur {
            return true
        }
    }
    return false
}

func (h *Handler) LobbyHandler(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Fatal(err)
		}

	lobby := []string{}

	go func(client *websocket.Conn, lb []string){
		ticker := time.NewTicker(waitTime)
		defer func() {
			ticker.Stop()
			client.Close()
		}()
			for {
					in := lobbyReq{}
				
					err := client.ReadJSON(&in)
					if err != nil {
						break
					}
			
					fmt.Printf("Got message: %#v\n", in)

					out := lobbyRespGenereic{}

					switch in.actionID {
						case "1": 
							if in.username == "" {
								break
							}
							out.actionUID = in.actionUID
							lb = append(lb, in.username)
							out.status = "success" 

							if err = client.WriteJSON(out); err != nil {
								break
							}
		
						case "2":
							if in.username == "" {
								break
							}
							out.actionUID = in.actionUID
							if contains(lb, in.username) {
								for _, cur := range lb {
									if cur == in.username {
										cur = ""
									}
								}
								out.status = "success"
								if err = client.WriteJSON(out); err != nil {
									break
								}
							} else {
								out.status = "failure"
								if err = client.WriteJSON(out); err != nil {
									break
								}			
							}
						case "3":
							out.actionUID = in.actionUID
							out.status = "success"
							if err = client.WriteJSON(out); err != nil {
								break
							}		
							
						case "4":
							out.actionUID = in.actionUID
							out.status = "success"
							if err = client.WriteJSON(out); err != nil {
								break
							}		
					}

					<-ticker.C
		}
	}(ws, lobby)
	return
}