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

func checkIfFileExists(filePath string) bool {
	info, error := os.Stat(filePath)
	fmt.Println("file infos:", info)
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
