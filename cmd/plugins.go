package cmd

import (
	"fmt"
	"os"

	"github.com/AmadlaOrg/raise/plugin"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	pluginsNew = plugin.New

	// PluginsCmd lists all discovered raise plugins.
	PluginsCmd = &cobra.Command{
		Use:   "plugins",
		Short: "List discovered raise plugins",
		RunE:  runPlugins,
	}
)

func runPlugins(cmd *cobra.Command, args []string) error {
	svc := pluginsNew()

	plugins, err := svc.Discover()
	if err != nil {
		return fmt.Errorf("failed to discover plugins: %w", err)
	}

	if len(plugins) == 0 {
		fmt.Fprintln(os.Stderr, "No raise plugins found in PATH.")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Plugin", "Engine", "Version", "Description")

	for _, name := range plugins {
		info, err := svc.GetInfo(name)
		if err != nil {
			table.Append(name, "?", "?", fmt.Sprintf("error: %v", err))
			continue
		}
		table.Append(name, info.Engine, info.Version, info.Description)
	}

	table.Render()
	return nil
}
