package cmd

import (
	"github.com/AmadlaOrg/raise/plugin"
	"github.com/spf13/cobra"
)

const (
	raiseVersion = "1.0.0"
)

// RaiseInfo holds the metadata for raise itself.
type RaiseInfo struct {
	Name        string   `json:"name" yaml:"name"`
	Version     string   `json:"version" yaml:"version"`
	Supports    []string `json:"supports" yaml:"supports"`
	Description string   `json:"description" yaml:"description"`
}

var (
	infoPluginNew  = plugin.New
	infoOutputFlag string
	infoHeryFlag   bool

	// InfoCmd outputs raise metadata.
	InfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Output raise metadata",
		Long:  "Outputs raise metadata including name, version, supported entity types, and description.",
		RunE:  runInfo,
	}
)

func init() {
	InfoCmd.Flags().StringVarP(&infoOutputFlag, "output", "o", "table", "Output format: table, json, yaml")
	InfoCmd.Flags().BoolVar(&infoHeryFlag, "hery", false, "Wrap output in HERY envelope (_type, _body)")
}

func runInfo(cmd *cobra.Command, args []string) error {
	out := RaiseInfo{
		Name:    "raise",
		Version: raiseVersion,
		Supports: []string{
			"amadla.org/entity/infrastructure@^v1.0.0",
		},
		Description: "Infrastructure provisioning with raise-* plugins",
	}

	svc := infoPluginNew()
	format := parseFormat(infoOutputFlag)
	if infoHeryFlag {
		return writeHeryInfoOutput(cmd.OutOrStdout(), format, out)
	}
	return writeOutput(cmd.OutOrStdout(), format, out, svc)
}
