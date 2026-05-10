package cmd

import (
	"ArgonVault/internal"
	"ArgonVault/internal/ui"
	"fmt"

	"github.com/spf13/cobra"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate-key",
	Short: "Rotate the master password",
	Long: `Change the master password for the entire store. Every vault is given a
fresh salt and every secret is re-encrypted under the new key, atomically.

You will be prompted for the current master password and the new one.

Example:
  argonv rotate-key`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.RotateMasterPassword(); err != nil {
			_ = internal.LogAction("ROTATE", "", "", "FAILURE")
			return fmt.Errorf("rotate master password: %w", err)
		}
		_ = internal.LogAction("ROTATE", "", "", "SUCCESS")
		ui.Success("master password rotated")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rotateCmd)
}
