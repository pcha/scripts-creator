package scripts

import (
	"io"
	"os"
	"path/filepath"
)

type Creator interface {
	Create(d Definition) (string, error)
}

type onDiskCreator struct {
	Sorter      Sorter
	CreationDir string
	Open        Opener
}

func GetCreator(creationDir string) Creator {
	return onDiskCreator{
		Sorter:      TopologicalSorter{},
		CreationDir: creationDir,
		Open: func(filename string) (io.WriteCloser, error) {
			return os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0774)
		},
	}
}

type Opener func(filename string) (io.WriteCloser, error)

func (c onDiskCreator) Create(d Definition) (string, error) {
	d.Tasks = c.Sorter.Sort(d.Tasks)
	path2file := filepath.Join(c.CreationDir, d.Filename)
	script, err := c.Open(path2file)
	if err != nil {
		return "", err
	}
	for _, t := range d.Tasks {
		_, err = script.Write([]byte(t.Command + "\n"))
		if err != nil {
			return "", err
		}
	}
	return path2file, script.Close()
}
