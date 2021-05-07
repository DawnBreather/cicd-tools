package path

import (
	"fmt"
	"os"
)

type Path struct {
	path string
}

const (
	FILE = iota
	DIRECTORY
)

func (p *Path) SetPath(path string) *Path {
	p.path = path
	return p
}
func (p *Path) GetPath() string {
	return p.path
}

// Exists returns whether the given file or directory exists
func (p *Path) Exists() bool {
	_, err := os.Stat(p.path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (p *Path) IsFileOrDir() (fileOrDir int, ok bool) {
	fi, err := os.Stat(p.path)
	if err != nil {
		fmt.Println(err)
		return -1, false
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return DIRECTORY, true
	case mode.IsRegular():
		return FILE, true
	}

	return -1, false
}
