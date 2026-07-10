package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/AmadlaOrg/raise/entity"
	"github.com/AmadlaOrg/raise/plugin"
	"github.com/spf13/cobra"
)

var (
	upProvider string
	upFilePath string

	upPluginNew            = plugin.New
	upReadProvider         = entity.ReadProvider
	upReadProviderFromData = entity.ReadProviderFromData

	// UpCmd provisions infrastructure via a raise plugin.
	UpCmd = &cobra.Command{
		Use:   "up [name]",
		Short: "Provision infrastructure using a raise plugin",
		Long:  "Provisions infrastructure by delegating to the specified raise-* plugin (--provider flag).",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runUp,
	}
)

func init() {
	UpCmd.Flags().StringVar(&upProvider, "provider", "", "Provider plugin name (e.g. libvirt, virtualbox, aws); auto-detected from entity file if omitted")
	UpCmd.Flags().StringVarP(&upFilePath, "file", "f", "", "Infrastructure definition file (YAML or JSON)")
}

func runUp(cmd *cobra.Command, args []string) error {
	provider := upProvider
	stdin := io.Reader(os.Stdin)

	// With -f -, the definition arrives on stdin (e.g. piped from
	// `hery compose --dir`); buffer it so the provider can be sniffed before
	// the bytes are handed to the plugin.
	var stdinData []byte
	if upFilePath == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		stdinData = data
		stdin = bytes.NewReader(stdinData)
	}

	// Auto-detect provider from the entity input if --provider is not specified.
	if provider == "" {
		switch {
		case upFilePath == "":
			return fmt.Errorf("either --provider or -f must be specified")
		case upFilePath == "-":
			p, err := upReadProviderFromData(stdinData)
			if err != nil {
				return fmt.Errorf("failed to detect provider from stdin: %w", err)
			}
			provider = p
		default:
			p, err := upReadProvider(upFilePath)
			if err != nil {
				return fmt.Errorf("failed to detect provider from entity file: %w", err)
			}
			provider = p
		}
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
	code, err := svc.Exec(pluginName, pluginArgs, stdin, os.Stdout, os.Stderr)
	if err != nil {
		return fmt.Errorf("failed to run up: %w", err)
	}
	if code != 0 {
		os.Exit(code)
	}

	return nil
}
