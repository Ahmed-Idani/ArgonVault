package cmd

import (
	"ArgonVault/internal"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type exportedSecret struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type exportedVault struct {
	Vault   string           `json:"vault"`
	Secrets []exportedSecret `json:"secrets"`
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the decrypted contents of a vault as JSON",
	Long: `Decrypt every secret in a vault and write a JSON document to stdout (or --output).

Example:
  argonv export --vault personal --output personal.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultName, _ := cmd.Flags().GetString("vault")
		output, _ := cmd.Flags().GetString("output")

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

		metas, err := internal.ListSecrets(v.ID)
		if err != nil {
			return fmt.Errorf("list secrets: %w", err)
		}

		pass, err := internal.EnsureMasterPassword()
		if err != nil {
			return err
		}

		bundle := exportedVault{Vault: v.Name, Secrets: make([]exportedSecret, 0, len(metas))}
		for _, m := range metas {
			s, err := internal.GetSecret(v.ID, m.Name)
			if err != nil {
				return fmt.Errorf("read secret %q: %w", m.Name, err)
			}
			plaintext, err := internal.Decrypt(&internal.EncryptedData{Ciphertext: s.Ciphertext, IV: s.IV}, pass, v.Salt)
			if err != nil {
				return fmt.Errorf("decrypt secret %q: %w", m.Name, err)
			}
			bundle.Secrets = append(bundle.Secrets, exportedSecret{
				Name:      m.Name,
				Value:     string(plaintext),
				CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			})
		}

		buf, err := json.MarshalIndent(bundle, "", "  ")
		if err != nil {
			return fmt.Errorf("encode json: %w", err)
		}

		if output == "" {
			fmt.Println(string(buf))
			return nil
		}
		if err := os.WriteFile(output, buf, 0o600); err != nil {
			return fmt.Errorf("write %s: %w", output, err)
		}
		fmt.Printf("exported %d secret(s) to %s\n", len(bundle.Secrets), output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("vault", "V", "", "vault to export")
	exportCmd.Flags().StringP("output", "o", "", "file to write to (default: stdout)")
}