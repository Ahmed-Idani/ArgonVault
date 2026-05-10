package cmd

import (
	"ArgonVault/internal"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve and decrypt a secret from a vault",
	Long: `Fetch a named secret from a vault and print the decrypted value.

Example:
  argonv get --vault personal --name aws_key`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultName, _ := cmd.Flags().GetString("vault")
		name, _ := cmd.Flags().GetString("name")

		if vaultName == "" {
			return fmt.Errorf("vault name is required (use --vault/-V)")
		}
		if name == "" {
			return fmt.Errorf("secret name is required (use --name/-n)")
		}

		v, err := internal.GetVault(vaultName)
		if err != nil {
			if errors.Is(err, internal.ErrNotFound) {
				return fmt.Errorf("vault %q does not exist", vaultName)
			}
			return fmt.Errorf("lookup vault: %w", err)
		}

		s, err := internal.GetSecret(v.ID, name)
		if err != nil {
			if errors.Is(err, internal.ErrNotFound) {
				_ = internal.LogAction("GET", vaultName, name, "FAILURE")
				return fmt.Errorf("secret %q not found in vault %q", name, vaultName)
			}
			_ = internal.LogAction("GET", vaultName, name, "FAILURE")
			return fmt.Errorf("read secret: %w", err)
		}

		pass, err := internal.EnsureMasterPassword()
		if err != nil {
			_ = internal.LogAction("GET", vaultName, name, "FAILURE")
			return err
		}

		plaintext, err := internal.Decrypt(&internal.EncryptedData{Ciphertext: s.Ciphertext, IV: s.IV}, pass, v.Salt)
		if err != nil {
			_ = internal.LogAction("GET", vaultName, name, "FAILURE")
			return fmt.Errorf("decrypt secret: %w", err)
		}

		_ = internal.LogAction("GET", vaultName, name, "SUCCESS")
		fmt.Println(string(plaintext))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("vault", "V", "", "vault holding the secret")
	getCmd.Flags().StringP("name", "n", "", "name of the secret")
}
