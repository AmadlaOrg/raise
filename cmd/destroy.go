package cmd

import (
	"fmt"
	"os"

	"github.com/AmadlaOrg/raise/plugin"
	"github.com/spf13/cobra"
)

var (
	destroyProvider string

	destroyPluginNew = plugin.New

	// DestroyCmd destroys infrastructure via a raise plugin.
	DestroyCmd = &cobra.Command{
		Use:   "destroy [name]",
		Short: "Destroy infrastructure using a raise plugin",
		Long:  "Destroys infrastructure by delegating to the specified raise-* plugin (--provider flag).",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDestroy,
	}
)

func init() {
	DestroyCmd.Flags().StringVar(&destroyProvider, "provider", "", "Provider plugin name (e.g. libvirt, virtualbox, aws)")
	_ = DestroyCmd.MarkFlagRequired("provider")
}

func runDestroy(cmd *cobra.Command, args []string) error {
	pluginName := "raise-" + destroyProvider

	pluginArgs := []string{"destroy"}
	if len(args) > 0 {
		pluginArgs = append(pluginArgs, args[0])
	}

	svc := destroyPluginNew()
	code, err := svc.Exec(pluginName, pluginArgs, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		return fmt.Errorf("failed to run destroy: %w", err)
	}
	if code != 0 {
		os.Exit(code)
	}

	return nil
}
