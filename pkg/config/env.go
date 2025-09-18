package config

import (
    "errors"
    "fmt"
    "os"
    "time"

    "github.com/joho/godotenv"
    v "github.com/spf13/viper"
)

// Defines the environment variables reader.
type Environment struct {
    viper *v.Viper
}

type EnvironmentScope struct {
    Env *Environment
}

// Instance to read environment variables.
var Env *Environment

var scopes map[string]*EnvironmentScope

func init() {
    scopes = map[string]*EnvironmentScope{}
}

// Creates a new environment variables scope and loads the .env file with the given path
// into the new scope.
func LoadScopedEnv(key, prefix, path string) (*EnvironmentScope, error) {
    _, found := scopes[key]

    if found {
        msg := fmt.Sprintf("scope with key %s was already created", key)
        return nil, errors.New(msg)
    }

    scope := &EnvironmentScope{
        Env: LoadEnv(prefix, path),
    }

    scopes[key] = scope

    return scope, nil
}

// Retrives the environment scope with the given key.
func Scope(key string) *EnvironmentScope {
    if scope, found := scopes[key]; found {
        return scope
    }

    return nil
}

// Loads the .env file with the given path.
// Automatically add environment variables (.env files included).
func LoadEnv(prefix, path string) *Environment {
    loadDotEnv(path)

    v := createViper(prefix)

    Env = &Environment{
        viper: v,
    }

    return Env
}

// Given the right key, it returns any value from an environment variable as an interface.
func (env *Environment) Get(key string) any {
    return env.viper.Get(key)
}

// Given the key, it returns the environment variable value as a string.
func (env *Environment) GetString(key string) string {
    return env.viper.GetString(key)
}

// Given the key, it returns the environment variable value as a string slice.
func (env *Environment) GetStringSlice(key string) []string {
    return env.viper.GetStringSlice(key)
}

// Given the right key, it returns the the environment variable value as a map, with string as key and interface as value.
func (env *Environment) GetStringMap(key string) map[string]any {
    return env.viper.GetStringMap(key)
}

// Given the key, it returns the environment variable value as bool.
func (env *Environment) GetBool(key string) bool {
    return env.viper.GetBool(key)
}

// Given the key, it returns the environment variable value as int.
func (env *Environment) GetInt(key string) int {
    return env.viper.GetInt(key)
}

// Given the key, it returns the environment variable value as float64.
func (env *Environment) GetFloat64(key string) float64 {
    return env.viper.GetFloat64(key)
}

// Given the key, it returns the environment variable value as time.Time.
func (env *Environment) GetTime(key string) time.Time {
    return env.viper.GetTime(key)
}

// Given the key, it returns the environment variable value as time.Duration.
func (env *Environment) GetDuration(key string) time.Duration {
    return env.viper.GetDuration(key)
}

// It returns the Viper instance to handle environment variables functionality
func (env *Environment) Viper() *v.Viper {
    return env.viper
}

func loadDotEnv(path string) {
    if path != "" {
        if _, err := os.Stat(path); err == nil {
            err := godotenv.Load(path)
            if err != nil {
                panic(err)
            }
        }
    }
}

func createViper(prefix string) *v.Viper {
    env := v.New()
    env.SetEnvPrefix(prefix)
    env.AutomaticEnv()
    return env
}