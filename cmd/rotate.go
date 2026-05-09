package cmd

import (
	"ArgonVault/internal"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate-key",
	Short: "Rotate the encryption key of a vault",
	Long: `Re-encrypt every secret in the vault under a freshly derived key.
You will be prompted for the current master password and the new one.

Example:
  argonv rotate-key --vault personal`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultName, _ := cmd.Flags().GetString("vault")
		if vaultName == "" {
			return fmt.Errorf("vault name is required (use --vault/-V)")
		}

		v, err := internal.GetVault(vaultName)
		if err != nil {
			if errors.Is(err, internal.ErrNotFound) {
				return fmt.Errorf("vault %q does not exist", vaultName)
			}
			return fmt.Errorf("lookup vault: %w", err)
		}

		if err := internal.RotateVaultKey(v); err != nil {
			_ = internal.LogAction("ROTATE", vaultName, "", "FAILURE")
			return fmt.Errorf("rotate key: %w", err)
		}

		_ = internal.LogAction("ROTATE", vaultName, "", "SUCCESS")
		fmt.Printf("encryption key rotated for vault %q\n", vaultName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rotateCmd)
	rotateCmd.Flags().StringP("vault", "V", "", "vault whose key to rotate")
}