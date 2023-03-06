package main

import (
	"net/http"

	"golang_auth/data"
	"golang_auth/structs"
	singing "golang_auth/singingController"
	cookie "golang_auth/coockiesController"
	h "golang_auth/helper"
)

type User = structs.User
var user User

func indexHandler(w http.ResponseWriter, r *http.Request) {
	user := cookie.GetUserCookie(r)
	if user != nil {
		h.Templating(w, "index", "base", *user)
	} else {
		h.Templating(w, "index", "base")
	}
}

func main() {
	data.Init("test.db")
	fs := http.FileServer(http.Dir("assets"))
	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/signup", singing.SignUpHandler)
	mux.HandleFunc("/login", singing.SignInHandler)
	mux.HandleFunc("/logout", singing.LogoutHandler)
	http.ListenAndServe(":8000", mux)
}
