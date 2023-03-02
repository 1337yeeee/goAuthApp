package singingController

import (
	"database/sql"
	"net/http"
	"log"
	"fmt"

	"movies_crud/structs"
	cookie "movies_crud/coockiesController"
	h "movies_crud/helper"
	"movies_crud/data"
)

type User = structs.User
var user User

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.Templating(w, "signup", "sign_layout")
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		user = User{
			Email: sql.NullString{String: r.Form.Get("email"), Valid: true},
			Password: sql.NullString{String: r.Form.Get("password"), Valid: true},
		}

		id, err := data.AddUser(user)
		fmt.Println(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "./", http.StatusFound)
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.Templating(w, "signin", "sign_layout")
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		userId, err := data.GetUserLogin(r.Form.Get("email"), r.Form.Get("password"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if userId == 0 {
			fmt.Println(userId)
		} else {
			user, err := data.GetUser(userId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w = cookie.SetUserCookie(w, user.ID)

			http.Redirect(w, r, "./", http.StatusFound)
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w = cookie.DelUserCookie(w)
	http.Redirect(w, r, "./", http.StatusFound)
}