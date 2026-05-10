package internal

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
)

const vaultSaltLen = 16

// RotateMasterPassword changes the master password for the entire store.
// Every vault gets a fresh random salt, every secret is re-encrypted under
// the new key, and configs.password_hash is updated — all in one transaction
// so a failure mid-rotation rolls back cleanly.
func RotateMasterPassword() error {
	storedHash, err := GetConfig(ConfigKeyPasswordHash)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return errors.New("master password is not set; nothing to rotate")
		}
		return err
	}

	oldPass, err := ReadPasswordPrompt("Current master password: ")
	if err != nil {
		return err
	}
	if !VerifyPassword(oldPass, storedHash) {
		_ = LogAction("FAILED_AUTH", "", "", "FAILURE")
		return errors.New("wrong master password")
	}

	newPass, err := ReadPasswordPrompt("New master password: ")
	if err != nil {
		return err
	}
	confirm, err := ReadPasswordPrompt("Confirm new master password: ")
	if err != nil {
		return err
	}
	if newPass != confirm {
		return errors.New("passwords do not match")
	}
	if newPass == oldPass {
		return errors.New("new password must differ from the current one")
	}

	newHash, err := HashPassword(newPass)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT id, name, salt FROM vaults`)
	if err != nil {
		return fmt.Errorf("read vaults: %w", err)
	}
	var vaults []Vault
	for rows.Next() {
		var v Vault
		if err := rows.Scan(&v.ID, &v.Name, &v.Salt); err != nil {
			rows.Close()
			return fmt.Errorf("scan vault: %w", err)
		}
		vaults = append(vaults, v)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	for _, v := range vaults {
		if err := rotateVault(tx, v, oldPass, newPass); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(
		`INSERT INTO configs(key, value) VALUES(?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		ConfigKeyPasswordHash, newHash,
	); err != nil {
		return fmt.Errorf("update password hash: %w", err)
	}

	return tx.Commit()
}

func rotateVault(tx *sql.Tx, v Vault, oldPass, newPass string) error {
	type item struct {
		id        int
		name      string
		plaintext []byte
	}

	rows, err := tx.Query(
		`SELECT id, name, ciphertext, iv FROM secrets WHERE vault_id = ?`, v.ID,
	)
	if err != nil {
		return fmt.Errorf("read secrets in vault %q: %w", v.Name, err)
	}
	var items []item
	for rows.Next() {
		var s Secret
		if err := rows.Scan(&s.ID, &s.Name, &s.Ciphertext, &s.IV); err != nil {
			rows.Close()
			return fmt.Errorf("scan secret: %w", err)
		}
		pt, err := Decrypt(&EncryptedData{Ciphertext: s.Ciphertext, IV: s.IV}, oldPass, v.Salt)
		if err != nil {
			rows.Close()
			return fmt.Errorf("decrypt %q in vault %q: %w", s.Name, v.Name, err)
		}
		items = append(items, item{id: s.ID, name: s.Name, plaintext: pt})
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	newSalt := make([]byte, vaultSaltLen)
	if _, err := rand.Read(newSalt); err != nil {
		return fmt.Errorf("generate salt: %w", err)
	}

	if _, err := tx.Exec(`UPDATE vaults SET salt = ? WHERE id = ?`, newSalt, v.ID); err != nil {
		return fmt.Errorf("update salt for vault %q: %w", v.Name, err)
	}

	for _, it := range items {
		ed, encErr := Encrypt(it.plaintext, newPass, newSalt)
		zero(it.plaintext)
		if encErr != nil {
			return fmt.Errorf("re-encrypt %q in vault %q: %w", it.name, v.Name, encErr)
		}
		if _, err := tx.Exec(
			`UPDATE secrets SET ciphertext = ?, iv = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
			ed.Ciphertext, ed.IV, it.id,
		); err != nil {
			return fmt.Errorf("update secret %q in vault %q: %w", it.name, v.Name, err)
		}
	}
	return nil
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
