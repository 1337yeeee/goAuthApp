package data

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	s "movies_crud/structs"
	"fmt"
)

type User = s.User

var DbName string

func Init(dbName string) error {
	DbName = dbName
	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL,
		password TEXT NOT NULL);
	`)
	if err != nil{
		return err
	}
	return nil
}

func AddUser(user User) (int, error) {
	var result_ sql.Result

	db, err := sql.Open("sqlite3", DbName)

	if err != nil {
		return 0, err
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)
	result_ = result

	if err != nil{
		fmt.Println(err)
		fmt.Println("\n\n\n\n")
		id, _ := result.LastInsertId()
		return  int(id), err
	}

	id, error := result_.LastInsertId()

	return int(id), error
}

func GetUsers() ([]User, error) {
	users := []User{}

	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		return users, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID, &u.Email)
		if err != nil{
			fmt.Println(err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

func GetUser(user_id int) (User, error) {
	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		return User{}, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, email FROM users WHERE id = ?", user_id)

	user := User{}
	err = row.Scan(&user.ID, &user.Email)

	return user, err
}

func IsUser(user_id int) (bool, error) {
	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		return false, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", user_id)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetUserLogin(email string, password string) (int, error) {
	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM users WHERE email = ? AND password = ?", email, password)

	var id int
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func DeleteUser(user User) error {
	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = ?", user.ID)

	return err
}

// func UpdateUser(user_id int, name string) error {
// 	db, err := sql.Open("sqlite3", DbName)
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// 	_, err = db.Exec("UPDATE users SET username = ? WHERE id = ?", name, user_id)
	
// 	return err
// }
