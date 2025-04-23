package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// This struct contains config info.
type Config struct {
	Port int    // file server port
	Root string // root dir
	Qr   bool   // show qr
	Open bool   // open browser
}

// This function loads YAML config, specified by a cmd-line argument.
func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("specified config file %s does not exist", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
