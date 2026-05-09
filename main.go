/*
Copyright © 2026 ahmed idani <ahmed.idani@insat.ucar.tn>
*/
package main

import (
	"ArgonVault/cmd"
	"ArgonVault/internal"
	"errors"
	"fmt"
	"os"
)

var DATA_PATH string = "./data/storage.db"

func checkIfFileExists(filePath string) bool {
	info, error := os.Stat(filePath)
	fmt.Println("file infos:", info)
	return !errors.Is(error, os.ErrNotExist)
}
func main() {
	if !checkIfFileExists(DATA_PATH) {
		internal.InitStorage(DATA_PATH)
	}
	cmd.Execute()
}
