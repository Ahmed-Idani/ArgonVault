/*
Copyright © 2026 ahmed idani <ahmed.idani@insat.ucar.tn>
*/
package main

import (
	"ArgonVault/cmd"
	"ArgonVault/internal"
	"errors"
	"os"
)

func checkIfFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}
func main() {
	if !checkIfFileExists(internal.DataPath) {
		if err := internal.InitStorage(); err != nil {
			panic(err)
		}
	}
	cmd.Execute()
}
