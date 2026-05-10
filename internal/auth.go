package internal

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const VaultVersion = "1"

const (
	ConfigKeyPasswordHash = "password_hash"
	ConfigKeyVaultVersion = "vault_version"
)

const passwordSaltLen = 16

func HashPassword(password string) (string, error) {
	salt := make([]byte, passwordSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return fmt.Sprintf(
		"argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func VerifyPassword(password, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 5 || parts[0] != "argon2id" {
		return false
	}
	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[2], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	got := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(expected)))
	return subtle.ConstantTimeCompare(expected, got) == 1
}

// EnsureMasterPassword prompts for the master password, performs first-time
// setup if no hash is recorded, otherwise verifies against the stored hash.
// Returns the password so callers can use it for vault key derivation.
func EnsureMasterPassword() (string, error) {
	stored, err := GetConfig(ConfigKeyPasswordHash)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return "", err
	}
	if stored == "" {
		pass, err := ReadPasswordPrompt("Set master password: ")
		if err != nil {
			return "", err
		}
		confirm, err := ReadPasswordPrompt("Confirm master password: ")
		if err != nil {
			return "", err
		}
		if pass != confirm {
			return "", fmt.Errorf("passwords do not match")
		}
		hash, err := HashPassword(pass)
		if err != nil {
			return "", err
		}
		if err := SetConfig(ConfigKeyPasswordHash, hash); err != nil {
			return "", err
		}
		return pass, nil
	}
	pass, err := ReadPassword()
	if err != nil {
		return "", err
	}
	if !VerifyPassword(pass, stored) {
		_ = LogAction("FAILED_AUTH", "", "", "FAILURE")
		return "", fmt.Errorf("wrong master password")
	}
	return pass, nil
}
