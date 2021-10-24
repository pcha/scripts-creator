package scripts

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeSorter struct {
	WasCalled bool
}

func (s *FakeSorter) Sort(ts []Task) []Task {
	s.WasCalled = true
	return ts
}

type FakeFile struct {
	path      string
	closed    bool
	writeable bool
	closeable bool
	content   string
}

func (f *FakeFile) Write(p []byte) (n int, err error) {
	if !f.writeable {
		return 0, errors.New("error on write")
	}
	f.content += string(p)
	return len(p), nil
}

func (f *FakeFile) Close() error {
	if !f.closeable {
		return errors.New("error on close")
	}
	f.closed = true
	return nil
}

func Test_onDiskCreator_Create(t *testing.T) {
	type TestCase struct {
		fakeFile    *FakeFile
		wontContent string
		wontPath    string
		wontErr     bool
	}
	cases := map[string]TestCase{
		"success - path with trailing slash": {
			fakeFile: &FakeFile{
				path:      "/var/gen/",
				closed:    false,
				writeable: true,
				closeable: true,
			},
			wontContent: "touch hellofile\necho 'hello world' > hellofile\nrm hellofile\n",
			wontPath:    "/var/gen/testscript.sh",
			wontErr:     false,
		},
		"success - path without trailing slash": {
			fakeFile: &FakeFile{
				path:      "/var/gen",
				closed:    false,
				writeable: true,
				closeable: true,
			},
			wontContent: "touch hellofile\necho 'hello world' > hellofile\nrm hellofile\n",
			wontPath:    "/var/gen/testscript.sh",
			wontErr:     false,
		},
		"error on open": {
			fakeFile: nil,
			wontErr:  true,
		},
		"error on write": {
			fakeFile: &FakeFile{
				path:      "/var/gen",
				closed:    false,
				writeable: false,
				closeable: true,
			},
			wontErr: true,
		},
		"error on close": {
			fakeFile: &FakeFile{
				path:      "/var/gen",
				closed:    false,
				writeable: true,
				closeable: false,
			},
			wontErr: true,
		},
	}

	for name, args := range cases {
		t.Run(name, func(t *testing.T) {
			sorter := new(FakeSorter)
			c := &onDiskCreator{
				sorter,
				"/var/gen",
				func(path string) (io.WriteCloser, error) {
					if args.fakeFile == nil {
						return nil, errors.New("can't open")
					}
					return args.fakeFile, nil
				},
			}
			def := Definition{
				Filename: "testscript.sh",
				Tasks: []Task{
					{
						Name:    "touch",
						Command: "touch hellofile",
					},
					{
						Name:    "echo",
						Command: "echo 'hello world' > hellofile",
						Dependencies: []string{
							"touch",
						},
					},
					{
						Name:    "rm",
						Command: "rm hellofile",
						Dependencies: []string{
							"echo",
						},
					},
				},
			}
			path, err := c.Create(def)

			assert.True(t, sorter.WasCalled)
			if args.wontErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, args.wontContent, args.fakeFile.content)
				assert.True(t, args.fakeFile.closed)
				assert.Equal(t, args.wontPath, path)
				assert.Nil(t, err)
			}
		})
	}
}
