package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/AmadlaOrg/raise/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPluginService struct {
	discoverResult []string
	discoverErr    error
	infoResults    map[string]*plugin.Info
	infoErrs       map[string]error
}

func (m *mockPluginService) Discover() ([]string, error) {
	return m.discoverResult, m.discoverErr
}

func (m *mockPluginService) GetInfo(name string) (*plugin.Info, error) {
	if m.infoErrs != nil {
		if err, ok := m.infoErrs[name]; ok {
			return nil, err
		}
	}
	if m.infoResults != nil {
		if info, ok := m.infoResults[name]; ok {
			return info, nil
		}
	}
	return nil, fmt.Errorf("no info for %s", name)
}

func (m *mockPluginService) Exec(pluginName string, args []string, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	return 0, nil
}

func TestRunInfo_OutputsRaiseMetadata(t *testing.T) {
	mock := &mockPluginService{}

	origNew := infoPluginNew
	defer func() { infoPluginNew = origNew }()
	infoPluginNew = func() plugin.Service { return mock }

	origFlag := infoOutputFlag
	defer func() { infoOutputFlag = origFlag }()
	infoOutputFlag = "json"

	var stdout bytes.Buffer
	cmd := InfoCmd
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	var info RaiseInfo
	err = json.Unmarshal(stdout.Bytes(), &info)
	require.NoError(t, err)

	assert.Equal(t, "raise", info.Name)
	assert.Equal(t, raiseVersion, info.Version)
	assert.Equal(t, "Infrastructure provisioning with raise-* plugins", info.Description)
	assert.Equal(t, []string{"amadla.org/entity/infrastructure@^v1.0.0"}, info.Supports)
}

func TestRunInfo_DoesNotAggregatePluginTypes(t *testing.T) {
	mock := &mockPluginService{
		discoverResult: []string{"raise-libvirt", "raise-aws"},
		infoResults: map[string]*plugin.Info{
			"raise-libvirt": {
				Name:     "raise-libvirt",
				Version:  "1.0.0",
				Engine:   "libvirt",
				Supports: []string{"amadla.org/entity/infrastructure@^v1.0.0", "amadla.org/entity/infrastructure/vm@^v1.0.0"},
			},
			"raise-aws": {
				Name:     "raise-aws",
				Version:  "1.0.0",
				Engine:   "aws",
				Supports: []string{"amadla.org/entity/infrastructure@^v1.0.0", "amadla.org/entity/infrastructure/cloud@^v1.0.0"},
			},
		},
	}

	origNew := infoPluginNew
	defer func() { infoPluginNew = origNew }()
	infoPluginNew = func() plugin.Service { return mock }

	origFlag := infoOutputFlag
	defer func() { infoOutputFlag = origFlag }()
	infoOutputFlag = "json"

	var stdout bytes.Buffer
	cmd := InfoCmd
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	var info RaiseInfo
	err = json.Unmarshal(stdout.Bytes(), &info)
	require.NoError(t, err)

	// raise core only declares its own entity type, not plugin sub-types
	assert.Equal(t, []string{"amadla.org/entity/infrastructure@^v1.0.0"}, info.Supports)
	assert.NotContains(t, info.Supports, "amadla.org/entity/infrastructure/vm@^v1.0.0")
	assert.NotContains(t, info.Supports, "amadla.org/entity/infrastructure/cloud@^v1.0.0")
}

func TestRunInfo_DefaultTableFormat(t *testing.T) {
	mock := &mockPluginService{}

	origNew := infoPluginNew
	defer func() { infoPluginNew = origNew }()
	infoPluginNew = func() plugin.Service { return mock }

	origFlag := infoOutputFlag
	defer func() { infoOutputFlag = origFlag }()
	infoOutputFlag = "table"

	var stdout bytes.Buffer
	cmd := InfoCmd
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	output := stdout.String()
	assert.Contains(t, output, "raise")
	assert.Contains(t, output, raiseVersion)
}

func TestRunInfo_YAMLFormat(t *testing.T) {
	mock := &mockPluginService{}

	origNew := infoPluginNew
	defer func() { infoPluginNew = origNew }()
	infoPluginNew = func() plugin.Service { return mock }

	origFlag := infoOutputFlag
	defer func() { infoOutputFlag = origFlag }()
	infoOutputFlag = "yaml"

	var stdout bytes.Buffer
	cmd := InfoCmd
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	output := stdout.String()
	assert.Contains(t, output, "name: raise")
	assert.Contains(t, output, "version:")
	assert.Contains(t, output, "supports:")
}
