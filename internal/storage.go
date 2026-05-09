package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var sqlCreateTables = []string{
	`CREATE TABLE vaults (
	id		INTEGER PRIMARY KEY AUTOINCREMENT,
	name	TEXT UNIQUE NOT NLL
	)`,
	`CREATE TABLE secrets (
		id		INTEGER PRIMARY KEY AUTOINCREMENT,
		name 	TEXT UNIQUE NOT NULL,
		ciphertext BLOB NOT NULL,
		iv BLOB NOT NULL,
		vault_id INTEGER NOT NULL,

		CONSTRAINT fk_vault FOREIGN KEY (vault_id) references vaults(id) ON DELETE CASCADE
	)`,
}

func init() {
	var err error
	db, err := sql.Open("sqlite3", "./argonvault.db")
	defer db.Close()
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	for _, query := range sqlCreateTables {
		_, err := db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Database initialized and connected.")
}
