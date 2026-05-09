package internal

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var sqlCreateTables = []string{
	`CREATE TABLE IF NOT EXISTS vaults (
	id		INTEGER PRIMARY KEY AUTOINCREMENT,
	name	TEXT UNIQUE NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS secrets (
		id					INTEGER PRIMARY KEY AUTOINCREMENT,
		name 				TEXT NOT NULL,
		ciphertext			BLOB NOT NULL,
		iv					BLOB NOT NULL,
		vault_id			INTEGER NOT NULL,
		created_at 			DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at 			DATETIME DEFAULT CURRENT_TIMESTAMP,

		CONSTRAINT fk_vault FOREIGN KEY (vault_id) references vaults(id) ON DELETE CASCADE,
    	UNIQUE (name, vault_id)
	);`,
	`CREATE TABLE IF NOT EXISTS configs (
		key			TEXT PRIMARY KEY,
		value		TEXT NOT NULL 
	);`,
	`CREATE TABLE IF NOT EXISTS audit_log (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    action      TEXT NOT NULL CHECK(action IN ('INIT', 'STORE', 'GET', 'DELETE', 'ROTATE', 'FAILED_AUTH')),
    secret_name TEXT,
    vault_name  TEXT,
    status      TEXT NOT NULL CHECK(status IN ('SUCCESS', 'FAILURE')),
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP);
    `,
}

func InitStorage(dataPath string) {
	var err error
	db, err := sql.Open("sqlite3", dataPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	if err = db.Ping(); err != nil {
		panic(err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		panic(err)
	}
	for _, query := range sqlCreateTables {
		_, err := db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
}
