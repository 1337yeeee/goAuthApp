package coockiesController

import (
	"net/http"
	"strconv"
	"fmt"

	"movies_crud/data"
	"movies_crud/structs"
)

type User = structs.User
var user User

func SetUserCookie(w http.ResponseWriter, userID int) http.ResponseWriter {
	cookie := http.Cookie{
		Name:	"user_id",
		Value:	strconv.Itoa(userID),
		Path:	"/",
		// MaxAge:   3600,
	}

	http.SetCookie(w, &cookie)

	return w
}

func GetUserCookie(r *http.Request) *User {
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

func DelUserCookie(w http.ResponseWriter) http.ResponseWriter {
	cookie := http.Cookie{
		Name:	"user_id",
		Value:	"",
		Path:	"/",
		MaxAge:   -1,
	}

	http.SetCookie(w, &cookie)

	return w
}