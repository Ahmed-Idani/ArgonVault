package cmd

import (
	"ArgonVault/internal"
	"ArgonVault/internal/ui"
	"fmt"

	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store a secret in a vault",
	Long: `Store a named secret inside an existing vault.

Example:
  argonv store --vault personal --name aws_key --data AKIAIOSFODNN7EXAMPLE`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultName, _ := cmd.Flags().GetString("vault")
		name, _ := cmd.Flags().GetString("name")
		data, _ := cmd.Flags().GetString("data")

		if vaultName == "" {
			return fmt.Errorf("vault name is required (use --vault/-V)")
		}
		if name == "" {
			return fmt.Errorf("secret name is required (use --name/-n)")
		}
		if data == "" {
			return fmt.Errorf("secret data is required (use --data/-d)")
		}

		v, err := internal.GetVault(vaultName)
		if err != nil {
			if err == internal.ErrNotFound {
				return fmt.Errorf("vault %q does not exist; run `argonv init --name %s` first", vaultName, vaultName)
			}
			return fmt.Errorf("lookup vault: %w", err)
		}
		pass, err := internal.EnsureMasterPassword()
		if err != nil {
			_ = internal.LogAction("STORE", vaultName, name, "FAILURE")
			return err
		}

		ed, err := internal.Encrypt([]byte(data), pass, v.Salt)
		if err != nil {
			_ = internal.LogAction("STORE", vaultName, name, "FAILURE")
			return fmt.Errorf("encrypt secret: %w", err)
		}

		if err := internal.SaveSecret(v.ID, name, ed.Ciphertext, ed.IV); err != nil {
			_ = internal.LogAction("STORE", vaultName, name, "FAILURE")
			return fmt.Errorf("save secret: %w", err)
		}

		_ = internal.LogAction("STORE", vaultName, name, "SUCCESS")
		ui.Success("secret %q stored in vault %q", name, vaultName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
	storeCmd.Flags().StringP("vault", "V", "", "vault to store the secret in")
	storeCmd.Flags().StringP("name", "n", "", "name of the secret")
	storeCmd.Flags().StringP("data", "d", "", "secret value to store")
}
