package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dua",
	Short: "Fast, friendly disk-usage CLI tool",
	Long:  `dua scans path and reports disk usage.`,
}

func initCfg() {
	viper.SetConfigName("dua")               // looks for dua.{yaml|toml|json}
	viper.AddConfigPath("$HOME/.config/dua") // path to look for the config file in

	// defaults
	viper.SetDefault("top", "10")

	// env: DUA_WORKERS, DUA_EXCLUDE, etc.
	viper.SetEnvPrefix("dua")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	_ = viper.ReadInConfig() // okay if not found
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize(initCfg)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
