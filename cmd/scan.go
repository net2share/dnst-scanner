package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan resolvers (placeholder)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ips, err := loadResolvers()
		if err != nil {
			return err
		}
		fmt.Printf("Loaded %d resolvers\n", len(ips))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().String("tunnel-domain", "", "Tunnel domain to test")
	scanCmd.MarkFlagRequired("tunnel-domain")
}
