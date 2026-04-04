package cmd

import (
	"fmt"
	"os"

	"github.com/AmadlaOrg/raise/plugin"
	"github.com/spf13/cobra"
)

var (
	pluginsNew        = plugin.New
	pluginsOutputFlag string
	pluginsHeryFlag   bool

	// PluginsCmd lists all discovered raise plugins.
	PluginsCmd = &cobra.Command{
		Use:   "plugins",
		Short: "List discovered raise plugins",
		RunE:  runPlugins,
	}
)

func init() {
	PluginsCmd.Flags().StringVarP(&pluginsOutputFlag, "output", "o", "table", "Output format: table, json, yaml")
	PluginsCmd.Flags().BoolVar(&pluginsHeryFlag, "hery", false, "Wrap output in HERY envelope (_type, _body)")
}

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

	var rows []pluginRow
	for _, name := range plugins {
		info, err := svc.GetInfo(name)
		if err != nil {
			rows = append(rows, pluginRow{
				Plugin:      name,
				Engine:      "?",
				Version:     "?",
				Description: fmt.Sprintf("error: %v", err),
			})
			continue
		}
		rows = append(rows, pluginRow{
			Plugin:      name,
			Engine:      info.Engine,
			Version:     info.Version,
			Description: info.Description,
		})
	}

	f := parseFormat(pluginsOutputFlag)
	if pluginsHeryFlag {
		return writeHeryPluginsOutput(os.Stdout, f, rows)
	}
	return writePluginsTable(os.Stdout, f, rows)
}
