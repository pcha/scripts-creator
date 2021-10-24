package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testGenDirVarName = "gendir_test"

func TestIntegration(t *testing.T) {
	// overriding env for test
	if testGenDir := os.Getenv(testGenDirVarName); testGenDir != "" {
		resetEnv, err := overrideEnv(genDirVarName, testGenDir)
		if err != nil {
			t.Skipf("Error overriding env var: %v", err.Error())
			return
		}
		defer func() {
			err := resetEnv()
			if err != nil {
				t.Fatal(err.Error())
			}
		}()
	}
	go main()
	reqBody := strings.NewReader(getExampleReqBody())
	resp, err := http.Post("http://localhost:8080/scripts", "application/json", reqBody)
	if err != nil {
		t.Error(err)
		return
	}
	requestOk := assert.Equal(t, http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if !requestOk {
		t.Error(string(respBodyBytes))
		return
	}
	var unparsedBody = &struct {
		Location string `json:"location"`
	}{}
	err = json.Unmarshal(respBodyBytes, unparsedBody)
	if err != nil {
		t.Error(err)
		return
	}
	scriptpath := unparsedBody.Location
	info, err := os.Stat(scriptpath)
	permStr := info.Mode().String()
	groupExec := permStr[6:7]
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "x", groupExec, "The group must have execution rights")
}
