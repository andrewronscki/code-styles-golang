package config

import (
    v "github.com/spf13/viper"
)

// Defines all the necessary input to load environment variables.
type StartConfig struct {
    // Defines a prefix that environment variables will use.
    // If your prefix is mtz, your environment variables should start with MTZ_
    Prefix string
    // The path of the .env file to be loaded.
    // It is relative to the application entrypoint, e.g: main.go.
    ConfigPath string
}

// Defines the instance type to be used to retrieve environment variables.
type Config struct {
    // The viper instance to handle environment variables.
    Standard *v.Viper
}

// Creates a new `Config` instance, given a `StartConfig`
// Loads the .env file with the given path.
// Automatically add environment variables (.env files included).
func NewConfig(cfg StartConfig) *Config {
    loadDotEnv(cfg.ConfigPath)

    config := createViper(cfg.Prefix)

    return &Config{Standard: config}
}