package internal

import (
	"database/sql"
	"fmt"
	"time"
)

type Vault struct {
	ID   int
	Name string
	Salt []byte
}

type Secret struct {
	ID         int
	Name       string
	Ciphertext []byte
	IV         []byte
	VaultID    int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type SecretMeta struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuditLog struct {
	ID         int
	Action     string
	SecretName sql.NullString
	VaultName  sql.NullString
	Status     string
	Timestamp  time.Time
}

var sqlCreateTables = []string{
	`CREATE TABLE IF NOT EXISTS vaults (
	id		INTEGER PRIMARY KEY AUTOINCREMENT,
	name	TEXT UNIQUE NOT NULL,
	salt    BLOB NOT NULL
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

func InitStorage() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	for _, query := range sqlCreateTables {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("create table: %w", err)
		}
	}
	if _, err := db.Exec(
		`INSERT INTO configs(key, value) VALUES(?, ?) ON CONFLICT(key) DO NOTHING`,
		ConfigKeyVaultVersion, VaultVersion,
	); err != nil {
		return fmt.Errorf("set vault version: %w", err)
	}
	return nil
}
