package cfgyaml

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(fPath string) error {
	f, err := os.ReadFile(fPath)
	if err != nil {
		return err
	}

	var cfg map[string]string

	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return err
	}

	if len(cfg) == 0 {
		return fmt.Errorf("Config file %q is empty.", fPath)
	}

	for k, v := range cfg {
		os.Setenv(k, v)
	}

	return nil
}
