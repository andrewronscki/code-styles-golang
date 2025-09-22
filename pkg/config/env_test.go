package config

import (
	"fmt"
	"testing"
	"time"

	afero "github.com/spf13/afero"
	v "github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setup_env(t *testing.T) func() {
	env := []byte(`
	MTZ_DUMMY=DUMMY
	MTZ_INTERFACE=20
	MTZ_SLICE=val1 val2
	MTZ_MAP={"key1":"val1", "key2":"val2"}
	MTZ_BOOL=true
	MTZ_INT=1
	MTZ_FLOAT=4.66
	MTZ_TIME="2020-04-03T22:50:45Z"
	MTZ_DURATION=4h`,
	)

	scope := []byte(`
	SC1_FOO=FOO
	SC2_BAR=BAR
	`)

	appFs := afero.NewOsFs()
	appFs.Mkdir("./env", 0755)
	afero.WriteFile(appFs, "./env/dummy.env", env, 0644)
	afero.WriteFile(appFs, "./env/scope.env", scope, 0644)

	return func() {
		t.Cleanup(func() {
			appFs.Remove("./env/dummy.env")
			appFs.Remove("./env/scope.env")
			appFs.Remove("./env")
			Env = nil
			scopes = map[string]*EnvironmentScope{}
		})
	}
}

func TestLoadEnv_WhenUseCorrectPrefixAndPath_ShouldInstanciateModuleInstance(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")

	// act
	check := func() {
		fmt.Printf("%t", Env != nil)
	}

	// assert
	assert.NotPanics(t, check)
}

func TestEnv_WhenAccessingNotLoadedModuleInstance_ShouldPanic(t *testing.T) {
	// arrange
	Env = nil

	// act
	check := func() {
		fmt.Printf("%s", Env.GetString("dummy"))
	}

	// assert
	assert.Panics(t, check)
}

func TestGet_WhenGivenTheRightKey_ShouldReturnInterface(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")

	// act
	i := Env.Get("interface")

	// assert
	assert.Equal(t, "20", i)
}

func TestGetString_WhenGivenTheRightKey_ShouldReturnString(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")

	// act
	dummy := Env.Get("dummy")

	// assert
	assert.IsType(t, "DUMMY", dummy)
}

func TestGetStringSlice_WhenGivenTheRightKey_ShouldReturnStringSlice(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")

	// act
	slice := Env.GetStringSlice("slice")

	// assert
	assert.Equal(t, []string{"val1", "val2"}, slice)
}

func TestGetStringMap_WhenGivenTheRightKey_ShouldReturnStringMap(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")

	expected := map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}

	// act
	m := Env.GetStringMap("map")

	// assert
	assert.Equal(t, expected, m)
}

func TestGetBool_WhenGivenTheRightKey_ShouldReturnBool(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")

	// act
	b := Env.GetBool("bool")

	// assert
	assert.True(t, b)
}

func TestGetInt_WhenGivenTheRightKey_ShouldReturnInt(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")

	// act
	integer := Env.GetInt("int")

	// assert
	assert.Equal(t, 1, integer)
}

func TestGetFloat64_WhenGivenTheRightKey_ShouldReturnFloat64(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")
	expected := 4.66

	// act
	integers := Env.GetFloat64("float")

	// assert
	assert.Equal(t, expected, integers)
}

func TestGetTime_WhenGivenTheRightKey_ShouldReturnTime(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")
	expected, _ := time.Parse(time.RFC3339, "2020-04-03T22:50:45Z")

	// act
	date := Env.GetTime("time")

	// assert
	assert.Equal(t, expected, date)
}

func TestGetDuration_WhenGivenTheRightKey_ShouldReturnDuration(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")
	expected, _ := time.ParseDuration("4h")

	// act
	hours := Env.GetDuration("duration")

	// assert
	assert.Equal(t, expected, hours)
}

func TestGetViper_WhenCallingViperFunction_ShouldReturnInstanceFromViper(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadEnv("MTZ", "env/dummy.env")
	expected := v.New()

	// act
	viper := Env.Viper()

	// assert
	assert.IsType(t, expected, viper)
}

func TestScopedEnvironment_WhenLoaded_ShoudHaveOnlyScopedVariables(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadScopedEnv("key", "MTZ", "env/dummy.env")
	LoadEnv("", "")
	scope := Scope("key")
	// act
	exists := scope.Env.GetString("dummy")
	notExists := Env.GetString("dummy")

	// assert
	assert.NotEmpty(t, exists)
	assert.Empty(t, notExists)
}

func TestScopedEnvironment_WhenKeyIsDuplicated_ShouldReturnError(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()
	LoadScopedEnv("key", "MTZ", "env/dummy.env")
	// act

	_, err := LoadScopedEnv("key", "MTZ", "env/dummy.env")

	// assert
	assert.Error(t, err)
}

func TestScopedEnvironment_WhenKeyNotExists_ShouldReturnNil(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()

	// act
	scope := Scope("key")

	// assert
	assert.Nil(t, scope)
}

func TestScopedEnvironment_WhenCoexistentScopes_ShouldReturnScopedVariables(t *testing.T) {
	// arrange
	cleanup := setup_env(t)
	defer cleanup()

	scope1, _ := LoadScopedEnv("scope1", "SC1", "env/scope.env")
	scope2, _ := LoadScopedEnv("scope2", "SC2", "env/scope.env")

	// act
	f1 := scope1.Env.GetString("foo")
	f2 := scope2.Env.GetString("foo")
	b1 := scope2.Env.GetString("bar")
	b2 := scope1.Env.GetString("bar")

	// assert
	assert.Equal(t, "FOO", f1)
	assert.Equal(t, "", f2)
	assert.Equal(t, "BAR", b1)
	assert.Equal(t, "", b2)
}
