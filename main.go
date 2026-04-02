package main

import (
	"fmt"
	"os"

	"github.com/AmadlaOrg/raise/cmd"
	"github.com/spf13/cobra"
)

const (
	appName = "raise"
	version = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:     appName,
	Short:   "Infrastructure provisioning CLI with raise-* plugins",
	Version: version,
}

func init() {
	rootCmd.AddCommand(cmd.InfoCmd)
	rootCmd.AddCommand(cmd.UpCmd)
	rootCmd.AddCommand(cmd.HaltCmd)
	rootCmd.AddCommand(cmd.DestroyCmd)
	rootCmd.AddCommand(cmd.SSHCmd)
	rootCmd.AddCommand(cmd.StatusCmd)
	rootCmd.AddCommand(cmd.PluginsCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
