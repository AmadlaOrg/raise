package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AmadlaOrg/raise/plugin"
	"github.com/spf13/cobra"
)

const (
	raiseVersion = "1.0.0"
)

// RaiseInfo holds the aggregated metadata for raise itself.
type RaiseInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	EntityTypes []string `json:"entity_types"`
	Description string   `json:"description"`
}

var (
	infoPluginNew = plugin.New

	// InfoCmd outputs raise metadata including aggregated plugin entity types.
	InfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Output raise metadata as JSON",
		Long:  "Outputs raise metadata including entity types aggregated from all discovered raise-* plugins.",
		RunE:  runInfo,
	}
)

func runInfo(cmd *cobra.Command, args []string) error {
	svc := infoPluginNew()

	// raise itself handles Infrastructure.
	seen := map[string]bool{
		"amadla.org/entity/infrastructure@v1.0.0": true,
	}
	entityTypes := []string{"amadla.org/entity/infrastructure@v1.0.0"}

	// Aggregate entity types from all discovered plugins.
	plugins, err := svc.Discover()
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "warning: failed to discover plugins: %v\n", err)
	} else {
		for _, name := range plugins {
			info, err := svc.GetInfo(name)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: failed to get info from %s: %v\n", name, err)
				continue
			}
			for _, et := range info.Supports {
				if !seen[et] {
					seen[et] = true
					entityTypes = append(entityTypes, et)
				}
			}
		}
	}

	out := RaiseInfo{
		Name:        "raise",
		Version:     raiseVersion,
		EntityTypes: entityTypes,
		Description: "Infrastructure provisioning with raise-* plugins",
	}

	data, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("failed to marshal info: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
