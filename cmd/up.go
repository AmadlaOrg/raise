package cmd

import (
	"fmt"
	"os"

	"github.com/AmadlaOrg/raise/plugin"
	"github.com/spf13/cobra"
)

var (
	upFrom     string
	upFilePath string

	upPluginNew = plugin.New

	// UpCmd provisions infrastructure via a raise plugin.
	UpCmd = &cobra.Command{
		Use:   "up [name]",
		Short: "Provision infrastructure using a raise plugin",
		Long:  "Provisions infrastructure by delegating to the specified raise-* plugin (--from flag).",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runUp,
	}
)

func init() {
	UpCmd.Flags().StringVar(&upFrom, "from", "", "Provider plugin name (e.g. libvirt, virtualbox, aws)")
	UpCmd.Flags().StringVarP(&upFilePath, "file", "f", "", "Infrastructure definition file (YAML or JSON)")
	_ = UpCmd.MarkFlagRequired("from")
}

func runUp(cmd *cobra.Command, args []string) error {
	pluginName := "raise-" + upFrom

	pluginArgs := []string{"up"}
	if len(args) > 0 {
		pluginArgs = append(pluginArgs, args[0])
	}
	if upFilePath != "" {
		pluginArgs = append(pluginArgs, "-f", upFilePath)
	}

	svc := upPluginNew()
	code, err := svc.Exec(pluginName, pluginArgs, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		return fmt.Errorf("failed to run up: %w", err)
	}
	if code != 0 {
		os.Exit(code)
	}

	return nil
}
