package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"bytes"
	"log"
	"fmt"

	"movies_crud/data"
	"movies_crud/structs"
)

type User = structs.User
var user User

const tplPath = "templates/"

func templating(w http.ResponseWriter, tmplName string, layout string, args ...any) {
	buf := &bytes.Buffer{}

	tmpl, err := template.New(layout).ParseFiles(tplPath+tmplName+".html", tplPath+layout+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(args) == 0 {
		err = tmpl.Execute(buf, nil)
	} else {
		err = tmpl.Execute(buf, args[0])
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

func setUserCookie(w http.ResponseWriter, userID int) http.ResponseWriter {
	cookie := http.Cookie{
		Name:	"user_id",
		Value:	strconv.Itoa(userID),
		Path:	"/",
		// MaxAge:   3600,
	}

	http.SetCookie(w, &cookie)

	return w
}

func getUserCookie(r *http.Request) *User {
	cookie, err := r.Cookie("user_id")
	if err == nil {
		if cookie.Value != "" {
			id, _ := strconv.Atoi(cookie.Value)
			user, err := data.GetUser(id)

			if err != nil {
				fmt.Println(err)
			}

			return &user
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func delUserCookie(w http.ResponseWriter) http.ResponseWriter {
	cookie := http.Cookie{
		Name:	"user_id",
		Value:	"",
		Path:	"/",
		MaxAge:   -1,
	}

	http.SetCookie(w, &cookie)

	return w
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserCookie(r)
	if user != nil {
		templating(w, "index", "base", *user)
	} else {
		templating(w, "index", "base")
	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templating(w, "signup", "sign_layout")
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

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templating(w, "signin", "sign_layout")
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

			w = setUserCookie(w, user.ID)

			http.Redirect(w, r, "./", http.StatusFound)
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w = delUserCookie(w)
	http.Redirect(w, r, "./", http.StatusFound)
}

func main() {
	data.Init("test.db")
	fs := http.FileServer(http.Dir("assets"))
	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/signup", signUpHandler)
	mux.HandleFunc("/login", signInHandler)
	mux.HandleFunc("/logout", logoutHandler)
	http.ListenAndServe(":8080", mux)
}
