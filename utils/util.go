package utils

import (
	"os"

	"sigs.k8s.io/yaml"
)

// LoadFromYaml reads a YAML file from the given path and unmarshals it into the provided interface.
func LoadFromYaml(filePath string, cfg interface{}) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, cfg)
}
