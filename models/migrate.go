package models

import (
	"database/sql"
)

//Migrate autoMigrate
func Migrate(db *sql.DB) {
	stmt, err := db.Prepare(`CREATE Table IF NOT EXISTS users(id int PRIMARY KEY, external_id int, name varchar(50), username varchar(30), language varchar(10), avatar varchar(100), created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP);`)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	stmt, err = db.Prepare(`CREATE Table IF NOT EXISTS bots(id int, name varchar(50), username varchar(30));`)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
}
