// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql"
)

type User struct {
	ID        string
	Firstname sql.NullString
	Lastname  sql.NullString
	Email     string
	Password  string
}
