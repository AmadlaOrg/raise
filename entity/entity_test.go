package entity

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadProvider(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "infrastructure.hery")

	content := `_type: amadla.org/entity/infrastructure@v1.0.0
_body:
  provider: libvirt
  ssh:
    user: root
    port: 22
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	provider, err := ReadProvider(path)
	require.NoError(t, err)
	assert.Equal(t, "libvirt", provider)
}

func TestReadProvider_AWS(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "infrastructure.hery")

	content := `_type: amadla.org/entity/infrastructure@v1.0.0
_body:
  provider: aws
  region: us-east-1
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	provider, err := ReadProvider(path)
	require.NoError(t, err)
	assert.Equal(t, "aws", provider)
}

func TestReadProvider_MissingProvider(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "infrastructure.hery")

	content := `_type: amadla.org/entity/infrastructure@v1.0.0
_body:
  ssh:
    user: root
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	_, err := ReadProvider(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no provider field")
}

func TestReadProvider_FileNotFound(t *testing.T) {
	_, err := ReadProvider("/nonexistent/path.hery")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read entity file")
}

func TestReadProvider_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.hery")

	require.NoError(t, os.WriteFile(path, []byte("just a string"), 0644))

	_, err := ReadProvider(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse entity file")
}

func TestReadProviderFromData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		provider string
		wantErr  error
	}{
		{
			name: "single infrastructure doc",
			input: `_type: amadla.org/entity/infrastructure@v1.0.0
_body:
  provider: libvirt
`,
			provider: "libvirt",
		},
		{
			name: "multi-doc graph with provider in first doc",
			input: `---
_type: amadla.org/entity/infrastructure@v1.0.0
_body:
  provider: libvirt
  ssh:
    user: root
---
_type: amadla.org/entity/infrastructure/vm@v1.0.0
_body:
  image: /iso/rocky.iso
`,
			provider: "libvirt",
		},
		{
			name: "multi-doc graph with provider in later doc",
			input: `---
_type: amadla.org/entity/system@v1.0.0
_body:
  hostname: demo
---
_type: amadla.org/entity/infrastructure@v1.0.0
_body:
  provider: aws
`,
			provider: "aws",
		},
		{
			name: "no provider anywhere",
			input: `---
_type: amadla.org/entity/system@v1.0.0
_body:
  hostname: demo
`,
			wantErr: ErrNoProvider,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: ErrNoProvider,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := ReadProviderFromData([]byte(tt.input))
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.provider, provider)
		})
	}
}
