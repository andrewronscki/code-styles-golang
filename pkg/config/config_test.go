package config

import (
	"testing"

	afero "github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func setup_config_tests(t *testing.T) func() {
	invalid := []byte(`
	{
		"foo": "bar"
	}`,
	)
	env := []byte(`
	MTZ_DUMMY=DUMMY`,
	)
	appFs := afero.NewOsFs()
	appFs.Mkdir("./env", 0755)
	afero.WriteFile(appFs, "./env/dummy.env", env, 0644)
	afero.WriteFile(appFs, "./env/invalid.env", invalid, 0644)
	return func() {
		t.Cleanup(func() {
			appFs.Remove("./env/dummy.env")
			appFs.Remove("./env/invalid.env")
			appFs.Remove("./env")
		})
	}
}

func TestConfigInstance(t *testing.T) {
	cleanup := setup_config_tests(t)
	defer cleanup()

	t.Run("should get environment variable without prefix", func(t *testing.T) {
		// arrange
		start := StartConfig{
			Prefix:     "MTZ",
			ConfigPath: "env/dummy.env",
		}

		cfg := NewConfig(start)

		// act
		dummy := cfg.Standard.GetString("dummy")

		// assert
		assert.Equal(t, "DUMMY", dummy)
	})

	t.Run("should get environment variable with prefix", func(t *testing.T) {
		// arrange
		start := StartConfig{
			ConfigPath: "env/dummy.env",
		}

		cfg := NewConfig(start)
		// act

		dummy := cfg.Standard.GetString("mtz_dummy")

		// assert
		assert.Equal(t, "DUMMY", dummy)
	})

	t.Run("invalid env file, should panic", func(t *testing.T) {
		// arrange
		start := StartConfig{
			Prefix:     "MTZ",
			ConfigPath: "env/invalid.env",
		}

		// act
		cfg := func() {
			NewConfig(start)
		}

		// assert
		assert.Panics(t, cfg)
	})
}
