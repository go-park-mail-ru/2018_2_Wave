package main

import (
	"Wave/handles"
	"Wave/server"
)

func main() {
	println(" -- -- -- -- -- -- -- -- -- --")

	srv := server.New()

	srv.Get("/index.html", srv.StaticServer)
	srv.Get("/app.bundle.js", srv.StaticServer)

	srv.Get("/img/avatars/:uid", handles.OnAvatarGET)

	srv.Post("/signup", handles.OnSignUpPOST)
	srv.Post("/login", handles.OnLogInPOST)

	srv.Get("/profile", handles.OnProfileGET)
	srv.Post("/profile", handles.OnProfilePOST)

	srv.Start(":8080")
}
