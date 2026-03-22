package cmd

import (
	"fmt"
	"os"

	"github.com/AmadlaOrg/raise/plugin"
	"github.com/spf13/cobra"
)

var (
	statusFrom string

	statusPluginNew = plugin.New

	// StatusCmd shows infrastructure status.
	StatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show infrastructure status",
		Long:  "Shows status of managed infrastructure. If --from is specified, queries a single provider; otherwise queries all discovered plugins.",
		RunE:  runStatus,
	}
)

func init() {
	StatusCmd.Flags().StringVar(&statusFrom, "from", "", "Provider plugin name (optional; queries all if omitted)")
}

func runStatus(cmd *cobra.Command, args []string) error {
	svc := statusPluginNew()

	if statusFrom != "" {
		pluginName := "raise-" + statusFrom
		code, err := svc.Exec(pluginName, []string{"status"}, os.Stdin, os.Stdout, os.Stderr)
		if err != nil {
			return fmt.Errorf("failed to get status from %s: %w", statusFrom, err)
		}
		if code != 0 {
			os.Exit(code)
		}
		return nil
	}

	// Query all discovered plugins for status.
	plugins, err := svc.Discover()
	if err != nil {
		return fmt.Errorf("failed to discover plugins: %w", err)
	}

	if len(plugins) == 0 {
		fmt.Fprintln(os.Stderr, "No raise plugins found in PATH.")
		return nil
	}

	for _, name := range plugins {
		fmt.Fprintf(os.Stdout, "--- %s ---\n", name)
		code, err := svc.Exec(name, []string{"status"}, os.Stdin, os.Stdout, os.Stderr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error querying %s: %v\n", name, err)
			continue
		}
		if code != 0 {
			fmt.Fprintf(os.Stderr, "%s exited with code %d\n", name, code)
		}
		fmt.Fprintln(os.Stdout)
	}

	return nil
}
