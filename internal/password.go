package internal

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func ReadPassword() (string, error) {
	return ReadPasswordPrompt("Master password: ")
}

func ReadPasswordPrompt(prompt string) (string, error) {
	fmt.Print(prompt)
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // move to next line after hidden input
	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}
	if len(passwordBytes) == 0 {
		return "", fmt.Errorf("password cannot be empty")
	}
	return string(passwordBytes), nil
}
