package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/marcosartorato/dua/internal/format"
	"github.com/marcosartorato/dua/internal/scan"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan a path and report disk usage",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		opts := scan.Options{
			TopN: viper.GetInt("top"),
		}

		res, warnings, err := scan.Run(context.Background(), path, opts)
		for _, w := range warnings {
			fmt.Fprintln(os.Stderr, "warn:", w)
		}
		if err != nil {
			return err
		}
		return format.PrintTable(os.Stdout, res)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().IntP("top", "n", 10, "Show only the top N largest entries")
	_ = viper.BindPFlag("top", scanCmd.Flags().Lookup("top"))
}
