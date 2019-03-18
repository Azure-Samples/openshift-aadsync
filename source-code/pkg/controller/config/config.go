package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// ControllerConfig holds the controller configuration details
type ControllerConfig struct {
	Namespace string   `yaml:"namespace"`
	Groups    []string `yaml:"groups"`
}

// NewControllerConfigFromFile creates a config from a file
func NewControllerConfigFromFile() (*ControllerConfig, error) {

	filepath := os.Getenv("AADSYNC_CONTROLLER_CONFIGFILE")
	if filepath == "" {
		filepath = "/etc/aadsynccontroller/config.yaml"
	}

	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var controllerConfig *ControllerConfig
	err = yaml.Unmarshal(configFile, &controllerConfig)
	if err != nil {
		return nil, err
	}

	return controllerConfig, nil
}
