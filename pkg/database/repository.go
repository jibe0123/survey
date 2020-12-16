package database

import "database/sql"

type Repository struct {
	Conn *sql.DB
}
