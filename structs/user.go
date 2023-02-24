package structs

import "database/sql"

type User struct {
	ID int
	Email sql.NullString
	Password sql.NullString
}