package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildConfig(t *testing.T) {
	genDirVal := "/var/generated"
	resetEnv, err := overrideEnv(genDirVarName, genDirVal)
	if err != nil {
		t.Skipf("Error overriding envvar %q", err.Error())
		return
	}
	defer func() {
		if err := resetEnv(); err != nil {
			t.Fatal(err)
		}
	}()

	cfg := buildConfig()
	assert.Equal(t, Config{genDirVal}, cfg)
}

func overrideEnv(name, value string) (func() error, error) {
	originalVal := os.Getenv(name)
	err := os.Setenv(name, value)
	return func() error {
		return os.Setenv(name, originalVal)
	}, err
}
