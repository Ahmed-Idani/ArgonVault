package cmd

import (
	"ArgonVault/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the secret vault",
	Long: `You create a vault logical structure to seperate secret concerns and usage
treat it as a directory'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("vault name is required (use --name/-n)")
		}
		if err := internal.CreateVault(name); err != nil {
			return fmt.Errorf("create vault: %w", err)
		}
		fmt.Printf("vault %q created\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("name", "n", "", "provide a name to the vault")
}