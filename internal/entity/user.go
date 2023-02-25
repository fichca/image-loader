package entity

import "database/sql"

type User struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	Login       string         `db:"login"`
	Password    string         `db:"password"`
	Description sql.NullString `db:"description"`
}
