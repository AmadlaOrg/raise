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

func TestRunInfo_AggregatesPluginEntityTypes(t *testing.T) {
	mock := &mockPluginService{
		discoverResult: []string{"raise-libvirt", "raise-aws"},
		infoResults: map[string]*plugin.Info{
			"raise-libvirt": {
				Name:     "raise-libvirt",
				Version:  "1.0.0",
				Engine:   "libvirt",
				Supports: []string{"amadla.org/entity/infrastructure@v1.0.0", "amadla.org/entity/infrastructure/vm@v1.0.0"},
			},
			"raise-aws": {
				Name:     "raise-aws",
				Version:  "1.0.0",
				Engine:   "aws",
				Supports: []string{"amadla.org/entity/infrastructure@v1.0.0", "amadla.org/entity/infrastructure/cloud/compute@v1.0.0"},
			},
		},
	}

	origNew := infoPluginNew
	defer func() { infoPluginNew = origNew }()
	infoPluginNew = func() plugin.Service { return mock }

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
	assert.Contains(t, info.EntityTypes, "amadla.org/entity/infrastructure@v1.0.0")
	assert.Contains(t, info.EntityTypes, "amadla.org/entity/infrastructure/vm@v1.0.0")
	assert.Contains(t, info.EntityTypes, "amadla.org/entity/infrastructure/cloud/compute@v1.0.0")
	// No duplicates — infrastructure@v1.0.0 appears once despite both plugins reporting it.
	count := 0
	for _, et := range info.EntityTypes {
		if et == "amadla.org/entity/infrastructure@v1.0.0" {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func TestRunInfo_NoPlugins(t *testing.T) {
	mock := &mockPluginService{
		discoverResult: nil,
	}

	origNew := infoPluginNew
	defer func() { infoPluginNew = origNew }()
	infoPluginNew = func() plugin.Service { return mock }

	var stdout bytes.Buffer
	cmd := InfoCmd
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	var info RaiseInfo
	err = json.Unmarshal(stdout.Bytes(), &info)
	require.NoError(t, err)

	assert.Equal(t, []string{"amadla.org/entity/infrastructure@v1.0.0"}, info.EntityTypes)
}

func TestRunInfo_PluginInfoError(t *testing.T) {
	mock := &mockPluginService{
		discoverResult: []string{"raise-broken"},
		infoErrs: map[string]error{
			"raise-broken": fmt.Errorf("connection refused"),
		},
	}

	origNew := infoPluginNew
	defer func() { infoPluginNew = origNew }()
	infoPluginNew = func() plugin.Service { return mock }

	var stdout, stderr bytes.Buffer
	cmd := InfoCmd
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	var info RaiseInfo
	err = json.Unmarshal(stdout.Bytes(), &info)
	require.NoError(t, err)

	// Still returns raise's own entity type even if plugin fails.
	assert.Equal(t, []string{"amadla.org/entity/infrastructure@v1.0.0"}, info.EntityTypes)
}
