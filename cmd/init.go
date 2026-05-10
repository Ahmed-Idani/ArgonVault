package cmd

import (
	"ArgonVault/internal"
	"ArgonVault/internal/ui"
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
		if _, err := internal.EnsureMasterPassword(); err != nil {
			return err
		}
		if err := internal.CreateVault(name); err != nil {
			return fmt.Errorf("create vault: %w", err)
		}
		ui.Success("vault %q created", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("name", "n", "", "provide a name to the vault")
}