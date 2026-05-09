package internal

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("not found")

// Vaults

func CreateVault(name string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO vaults(name) VALUES(?)`, name)
	return err
}

func GetVault(name string) (*Vault, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	v := &Vault{}
	err = db.QueryRow(`SELECT id, name, salt FROM vaults WHERE name = ?`, name).
		Scan(&v.ID, &v.Name, &v.Salt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return v, nil
}

func ListVaults() ([]Vault, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, name, salt FROM vaults ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vaults []Vault
	for rows.Next() {
		var v Vault
		if err := rows.Scan(&v.ID, &v.Name, &v.Salt); err != nil {
			return nil, err
		}
		vaults = append(vaults, v)
	}
	return vaults, rows.Err()
}

func DeleteVault(name string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec(`DELETE FROM vaults WHERE name = ?`, name)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// Secrets

func SaveSecret(vaultID int, name string, ciphertext []byte, iv []byte) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO secrets(name, ciphertext, iv, vault_id) VALUES(?, ?, ?, ?)`,
		name, ciphertext, iv, vaultID,
	)
	return err
}

func GetSecret(vaultID int, name string) (*Secret, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	s := &Secret{}
	err = db.QueryRow(
		`SELECT id, name, ciphertext, iv, vault_id, created_at, updated_at
		 FROM secrets WHERE name = ? AND vault_id = ?`,
		name, vaultID,
	).Scan(&s.ID, &s.Name, &s.Ciphertext, &s.IV, &s.VaultID, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func ListSecrets(vaultID int) ([]SecretMeta, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		`SELECT name, created_at, updated_at FROM secrets WHERE vault_id = ? ORDER BY name`,
		vaultID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metas []SecretMeta
	for rows.Next() {
		var m SecretMeta
		if err := rows.Scan(&m.Name, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		metas = append(metas, m)
	}
	return metas, rows.Err()
}

func DeleteSecret(vaultID int, name string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec(`DELETE FROM secrets WHERE name = ? AND vault_id = ?`, name, vaultID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func UpdateSecret(vaultID int, name string, ciphertext []byte, iv []byte) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec(
		`UPDATE secrets SET ciphertext = ?, iv = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE name = ? AND vault_id = ?`,
		ciphertext, iv, name, vaultID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// Config

func SetConfig(key string, value string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO configs(key, value) VALUES(?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	return err
}

func GetConfig(key string) (string, error) {
	db, err := OpenDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var value string
	err = db.QueryRow(`SELECT value FROM configs WHERE key = ?`, key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

// Audit

func LogAction(action string, vaultName string, secretName string, status string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO audit_log(action, vault_name, secret_name, status) VALUES(?, ?, ?, ?)`,
		action, nullableString(vaultName), nullableString(secretName), status,
	)
	return err
}

func GetLogs() ([]AuditLog, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		`SELECT id, action, secret_name, vault_name, status, timestamp
		 FROM audit_log ORDER BY timestamp DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanLogs(rows)
}

func GetLogsByVault(vaultName string) ([]AuditLog, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		`SELECT id, action, secret_name, vault_name, status, timestamp
		 FROM audit_log WHERE vault_name = ? ORDER BY timestamp DESC`,
		vaultName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanLogs(rows)
}

func scanLogs(rows *sql.Rows) ([]AuditLog, error) {
	var logs []AuditLog
	for rows.Next() {
		var l AuditLog
		if err := rows.Scan(&l.ID, &l.Action, &l.SecretName, &l.VaultName, &l.Status, &l.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}