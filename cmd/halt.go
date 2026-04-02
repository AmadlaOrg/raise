package cmd

import (
	"fmt"
	"os"

	"github.com/AmadlaOrg/raise/plugin"
	"github.com/spf13/cobra"
)

var (
	haltProvider string

	haltPluginNew = plugin.New

	// HaltCmd stops infrastructure via a raise plugin.
	HaltCmd = &cobra.Command{
		Use:   "halt [name]",
		Short: "Stop infrastructure using a raise plugin",
		Long:  "Stops (halts) infrastructure by delegating to the specified raise-* plugin (--provider flag).",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runHalt,
	}
)

func init() {
	HaltCmd.Flags().StringVar(&haltProvider, "provider", "", "Provider plugin name (e.g. libvirt, virtualbox, aws)")
	_ = HaltCmd.MarkFlagRequired("provider")
}

func runHalt(cmd *cobra.Command, args []string) error {
	pluginName := "raise-" + haltProvider

	pluginArgs := []string{"halt"}
	if len(args) > 0 {
		pluginArgs = append(pluginArgs, args[0])
	}

	svc := haltPluginNew()
	code, err := svc.Exec(pluginName, pluginArgs, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		return fmt.Errorf("failed to run halt: %w", err)
	}
	if code != 0 {
		os.Exit(code)
	}

	return nil
}
