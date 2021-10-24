package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"scripts-creator/cmd/api/scripts"

	"github.com/stretchr/testify/assert"
)

type FakeCreator struct {
	path      string
	err       error
	WasCalled bool
}

func (c *FakeCreator) Create(_ scripts.Definition) (string, error) {
	c.WasCalled = true
	return c.path, c.err
}

func TestHandler_createScript(t *testing.T) {
	type testCase struct {
		reqBodyStr       string
		creator          *FakeCreator
		wantsCallCreator bool
		wantedStatusCode int
	}

	cases := map[string]testCase{
		"success": {
			reqBodyStr: getExampleReqBody(),
			creator: &FakeCreator{
				path: "/var/tmp/myscript.sh",
			},
			wantsCallCreator: true,
			wantedStatusCode: http.StatusCreated,
		},
		"bad request": {
			reqBodyStr:       "saraza",
			wantsCallCreator: false,
			wantedStatusCode: http.StatusBadRequest,
		},
		"internal error": {
			reqBodyStr: getExampleReqBody(),
			creator: &FakeCreator{
				err: errors.New("erro creating script"),
			},
			wantsCallCreator: true,
			wantedStatusCode: http.StatusInternalServerError,
		},
	}

	for name, args := range cases {
		t.Run(name, func(t *testing.T) {
			h := Handler{
				creator: args.creator,
			}
			r := setupRouter(h)
			resp := httptest.NewRecorder()
			reqBody := strings.NewReader(args.reqBodyStr)
			req := httptest.NewRequest("POST", "/scripts", reqBody)
			r.ServeHTTP(resp, req)

			if args.wantsCallCreator {
				assert.True(t, args.creator.WasCalled, "The scripts creator must be called")
			}
			assert.Equal(t, args.wantedStatusCode, resp.Code)
		})
	}
}

func getExampleReqBody() string {
	return `{
    "filename": "myscript.sh",
    "tasks": [
        {
            "name": "rm",
            "command": "rm -f /tmp/test",
            "dependencies": [
                "cat"
            ]
        },
        {
            "name": "cat",
            "command": "cat /tmp/test",
            "dependencies": [
                "chown",
                "chmod"
            ]
        },
        {
            "name": "touch",
            "command": "touch /tmp/test"
        },
        {
            "name": "chmod",
            "command": "chmod 600 /tmp/test",
            "dependencies": [
                "touch"
            ]
        },
        {
            "name": "chown",
            "command": "chown root:root /tmp/test",
            "dependencies": [
                "touch"
            ]
        }
    ]
}`
}
