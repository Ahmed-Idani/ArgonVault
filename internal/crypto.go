package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

type EncryptedData struct {
	Ciphertext []byte
	IV         []byte
}

const (
	argonTime    = 3
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
)

func deriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
}

func Encrypt(plaintext []byte, password string, salt []byte) (*EncryptedData, error) {
	if len(salt) == 0 {
		return nil, errors.New("salt cannot be empty")
	}
	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, gcm.NonceSize()) // 12 bytes
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, iv, plaintext, nil)

	return &EncryptedData{Ciphertext: ciphertext, IV: iv}, nil
}

func Decrypt(ed *EncryptedData, password string, salt []byte) ([]byte, error) {
	if ed == nil {
		return nil, errors.New("encrypted data is nil")
	}
	if len(salt) == 0 {
		return nil, errors.New("salt cannot be empty")
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, ed.IV, ed.Ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed: wrong password or tampered data")
	}

	return plaintext, nil
}
