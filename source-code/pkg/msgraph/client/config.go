package client

import (
	"errors"
	"os"
)

// Config holds the configuration details required for the client
type Config struct {
	TenantID     string
	ClientID     string
	ClientSecret string
}

// NewConfigFromEnvironmentVariables creates a config from environment variables
func NewConfigFromEnvironmentVariables() (*Config, error) {

	config := &Config{
		TenantID:     os.Getenv("AZURE_TENANT_ID"),
		ClientID:     os.Getenv("AZURE_CLIENT_ID"),
		ClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
	}

	if config.TenantID == "" || config.ClientID == "" || config.ClientSecret == "" {
		err := errors.New("Error creating config from environment variables: the AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET environment variables must be set")
		return nil, err
	}

	return config, nil
}
