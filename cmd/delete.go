package cmd

import (
	"ArgonVault/internal"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a secret, or an entire vault",
	Long: `Delete a single secret from a vault, or delete the whole vault when --name is omitted.

Examples:
  argonv delete --vault personal --name aws_key
  argonv delete --vault personal`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultName, _ := cmd.Flags().GetString("vault")
		name, _ := cmd.Flags().GetString("name")

		if vaultName == "" {
			return fmt.Errorf("vault name is required (use --vault/-V)")
		}

		if name == "" {
			if err := internal.DeleteVault(vaultName); err != nil {
				if errors.Is(err, internal.ErrNotFound) {
					return fmt.Errorf("vault %q does not exist", vaultName)
				}
				_ = internal.LogAction("DELETE", vaultName, "", "FAILURE")
				return fmt.Errorf("delete vault: %w", err)
			}
			_ = internal.LogAction("DELETE", vaultName, "", "SUCCESS")
			fmt.Printf("vault %q deleted\n", vaultName)
			return nil
		}

		v, err := internal.GetVault(vaultName)
		if err != nil {
			if errors.Is(err, internal.ErrNotFound) {
				return fmt.Errorf("vault %q does not exist", vaultName)
			}
			return fmt.Errorf("lookup vault: %w", err)
		}

		if err := internal.DeleteSecret(v.ID, name); err != nil {
			if errors.Is(err, internal.ErrNotFound) {
				return fmt.Errorf("secret %q not found in vault %q", name, vaultName)
			}
			_ = internal.LogAction("DELETE", vaultName, name, "FAILURE")
			return fmt.Errorf("delete secret: %w", err)
		}

		_ = internal.LogAction("DELETE", vaultName, name, "SUCCESS")
		fmt.Printf("secret %q deleted from vault %q\n", name, vaultName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("vault", "V", "", "vault to operate on")
	deleteCmd.Flags().StringP("name", "n", "", "name of the secret to delete (omit to delete the whole vault)")
}