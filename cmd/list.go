package cmd

import (
	"ArgonVault/internal"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List vaults, or secrets within a vault",
	Long: `Without --vault, list all vaults.
With --vault, list secret names (and timestamps) in that vault.

Examples:
  argonv list
  argonv list --vault personal`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultName, _ := cmd.Flags().GetString("vault")

		if vaultName == "" {
			vaults, err := internal.ListVaults()
			if err != nil {
				return fmt.Errorf("list vaults: %w", err)
			}
			if len(vaults) == 0 {
				fmt.Println("no vaults yet — create one with `argonv init --name <name>`")
				return nil
			}
			for _, v := range vaults {
				fmt.Println(v.Name)
			}
			return nil
		}

		v, err := internal.GetVault(vaultName)
		if err != nil {
			if errors.Is(err, internal.ErrNotFound) {
				return fmt.Errorf("vault %q does not exist", vaultName)
			}
			return fmt.Errorf("lookup vault: %w", err)
		}

		secrets, err := internal.ListSecrets(v.ID)
		if err != nil {
			return fmt.Errorf("list secrets: %w", err)
		}
		if len(secrets) == 0 {
			fmt.Printf("vault %q is empty\n", vaultName)
			return nil
		}
		for _, s := range secrets {
			fmt.Printf("%s\t(created %s, updated %s)\n",
				s.Name,
				s.CreatedAt.Format("2006-01-02 15:04"),
				s.UpdatedAt.Format("2006-01-02 15:04"),
			)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("vault", "V", "", "vault whose secrets to list (omit to list vaults)")
}