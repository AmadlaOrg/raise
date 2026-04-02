package entity

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Infrastructure represents the minimal fields needed from an Infrastructure entity.
type Infrastructure struct {
	Type string              `yaml:"_type"`
	Body InfrastructureBody  `yaml:"_body"`
}

// InfrastructureBody holds the body of an Infrastructure entity.
type InfrastructureBody struct {
	Provider string `yaml:"provider"`
}

var osReadFile = os.ReadFile

// ReadProvider reads the provider field from an Infrastructure entity file.
func ReadProvider(path string) (string, error) {
	data, err := osReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read entity file %s: %w", path, err)
	}

	var entity Infrastructure
	if err := yaml.Unmarshal(data, &entity); err != nil {
		return "", fmt.Errorf("failed to parse entity file %s: %w", path, err)
	}

	if entity.Body.Provider == "" {
		return "", fmt.Errorf("no provider field in entity file %s", path)
	}

	return entity.Body.Provider, nil
}
