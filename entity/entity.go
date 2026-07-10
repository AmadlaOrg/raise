package entity

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Infrastructure represents the minimal fields needed from an Infrastructure entity.
type Infrastructure struct {
	Type string             `yaml:"_type"`
	Body InfrastructureBody `yaml:"_body"`
}

// InfrastructureBody holds the body of an Infrastructure entity.
type InfrastructureBody struct {
	Provider string `yaml:"provider"`
}

var osReadFile = os.ReadFile

// ErrNoProvider is returned when no document in the input carries
// _body.provider.
var ErrNoProvider = errors.New("no provider field")

// ReadProvider reads the provider field from an Infrastructure entity file.
func ReadProvider(path string) (string, error) {
	data, err := osReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read entity file %s: %w", path, err)
	}

	provider, err := ReadProviderFromData(data)
	switch {
	case errors.Is(err, ErrNoProvider):
		return "", fmt.Errorf("%w in entity file %s", ErrNoProvider, path)
	case err != nil:
		return "", fmt.Errorf("failed to parse entity file %s: %w", path, err)
	}
	return provider, nil
}

// ReadProviderFromData finds the provider field in entity data: a single
// Infrastructure entity or a multi-doc YAML stream as emitted by
// `hery compose --dir` (the first document carrying _body.provider wins).
func ReadProviderFromData(data []byte) (string, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	for {
		var entity Infrastructure
		if err := dec.Decode(&entity); err != nil {
			if errors.Is(err, io.EOF) {
				return "", ErrNoProvider
			}
			return "", err
		}
		if entity.Body.Provider != "" {
			return entity.Body.Provider, nil
		}
	}
}
