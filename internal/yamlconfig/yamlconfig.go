package yamlconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(fPath string) error {
	f, err := os.ReadFile(fPath)
	if err != nil {
		return fmt.Errorf("failed to read: %w", err)
	}

	var cfg map[string]string

	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	if len(cfg) == 0 {
		return fmt.Errorf("config file %q is empty", fPath)
	}

	for k, v := range cfg {
		os.Setenv(k, v)
	}

	return nil
}
