package cmd

import (
	"fmt"
	"os"

	"github.com/AmadlaOrg/raise/entity"
	"github.com/AmadlaOrg/raise/plugin"
	"github.com/spf13/cobra"
)

var (
	upFrom     string
	upFilePath string

	upPluginNew   = plugin.New
	upReadProvider = entity.ReadProvider

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
	UpCmd.Flags().StringVar(&upFrom, "from", "", "Provider plugin name (e.g. libvirt, virtualbox, aws); auto-detected from entity file if omitted")
	UpCmd.Flags().StringVarP(&upFilePath, "file", "f", "", "Infrastructure definition file (YAML or JSON)")
}

func runUp(cmd *cobra.Command, args []string) error {
	provider := upFrom

	// Auto-detect provider from entity file if --from is not specified.
	if provider == "" {
		if upFilePath == "" {
			return fmt.Errorf("either --from or -f must be specified")
		}
		p, err := upReadProvider(upFilePath)
		if err != nil {
			return fmt.Errorf("failed to detect provider from entity file: %w", err)
		}
		provider = p
	}

	pluginName := "raise-" + provider

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
