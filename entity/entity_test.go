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
