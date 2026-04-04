package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/AmadlaOrg/raise/plugin"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

type format int

const (
	formatTable format = iota
	formatJSON
	formatYAML
)

func parseFormat(value string) format {
	switch strings.ToLower(value) {
	case "json":
		return formatJSON
	case "yaml":
		return formatYAML
	default:
		return formatTable
	}
}

func writeOutput(w io.Writer, f format, data any, _ plugin.Service) error {
	switch f {
	case formatJSON:
		return writeJSON(w, data)
	case formatYAML:
		return writeYAML(w, data)
	default:
		return writeInfoTable(w, data)
	}
}

func writeJSON(w io.Writer, data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, err = fmt.Fprintln(w, string(bytes))
	return err
}

func writeYAML(w io.Writer, data any) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	_, err = fmt.Fprint(w, string(bytes))
	return err
}

func writeInfoTable(w io.Writer, data any) error {
	switch v := data.(type) {
	case RaiseInfo:
		table := tablewriter.NewWriter(w)
		table.Header("Field", "Value")
		table.Append("Name", v.Name)
		table.Append("Version", v.Version)
		table.Append("Supports", strings.Join(v.Supports, ", "))
		table.Append("Description", v.Description)
		table.Render()
	default:
		return writeJSON(w, data)
	}
	return nil
}

func writePluginsTable(w io.Writer, f format, plugins []pluginRow) error {
	switch f {
	case formatJSON:
		return writeJSON(w, plugins)
	case formatYAML:
		return writeYAML(w, plugins)
	default:
		table := tablewriter.NewWriter(w)
		table.Header("Plugin", "Engine", "Version", "Description")
		for _, p := range plugins {
			table.Append(p.Plugin, p.Engine, p.Version, p.Description)
		}
		table.Render()
		return nil
	}
}

type pluginRow struct {
	Plugin      string `json:"plugin" yaml:"plugin"`
	Engine      string `json:"engine" yaml:"engine"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

// heryEnvelope ensures _type is serialized before _body in both JSON and YAML.
type heryEnvelope struct {
	Type string `json:"_type" yaml:"_type"`
	Body any    `json:"_body" yaml:"_body"`
}

func writeHeryInfoOutput(w io.Writer, f format, data any) error {
	return writeHeryEnvelope(w, f, "amadla.org/entity/tools/info@v1.0.0", data)
}

func writeHeryPluginsOutput(w io.Writer, f format, data any) error {
	return writeHeryEnvelope(w, f, "amadla.org/entity/tools/plugins@v1.0.0", data)
}

func writeHeryEnvelope(w io.Writer, f format, entityType string, data any) error {
	envelope := heryEnvelope{
		Type: entityType,
		Body: data,
	}

	switch f {
	case formatJSON:
		return writeJSON(w, envelope)
	case formatTable:
		// Render _type as header, then table for _body
		fmt.Fprintf(w, "_type: %s\n\n", entityType)
		switch v := data.(type) {
		case RaiseInfo:
			return writeInfoTable(w, v)
		case []pluginRow:
			return writePluginsTable(w, formatTable, v)
		default:
			return writeYAML(w, data)
		}
	default:
		return writeYAML(w, envelope)
	}
}
