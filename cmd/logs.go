package cmd

import (
	"ArgonVault/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show audit log entries",
	Long: `Print audit log entries (newest first). With --vault, only entries for that vault.

Examples:
  argonv logs
  argonv logs --vault personal`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultName, _ := cmd.Flags().GetString("vault")

		var (
			logs []internal.AuditLog
			err  error
		)
		if vaultName == "" {
			logs, err = internal.GetLogs()
		} else {
			logs, err = internal.GetLogsByVault(vaultName)
		}
		if err != nil {
			return fmt.Errorf("read audit log: %w", err)
		}
		if len(logs) == 0 {
			fmt.Println("no audit entries")
			return nil
		}

		for _, l := range logs {
			vault := "-"
			if l.VaultName.Valid {
				vault = l.VaultName.String
			}
			secret := "-"
			if l.SecretName.Valid {
				secret = l.SecretName.String
			}
			fmt.Printf("%s\t%s\t%s\tvault=%s\tsecret=%s\n",
				l.Timestamp.Format("2006-01-02 15:04:05"),
				l.Action,
				l.Status,
				vault,
				secret,
			)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringP("vault", "V", "", "filter to a single vault")
}