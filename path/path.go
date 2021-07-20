package path

import (
	"github.com/DawnBreather/go-commons/logger"
	"os"
	"path/filepath"
	"strings"
)

type Path struct {
	path string
}

const (
	FILE = iota
	DIRECTORY
)

var _logger = logger.New()

func (p *Path) SetPath(path string) *Path {
	p.path = filepath.FromSlash(path)
	return p
}
func (p *Path) GetPath() string {
	return p.path
}
func (p *Path) SetCompositePath(subPaths ...string) *Path {
	path := strings.Join(subPaths, "/")
	p.path = filepath.FromSlash(path)
	return p
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

func (p *Path) MkdirAll(mode os.FileMode){
	if err := os.MkdirAll(p.path, mode); err != nil {
		_logger.Errorf("Unable to `mkdir -p` on path { %s }: %v", p.path, err)
	}
}

func (p *Path) IsFileOrDir() int {
	fi, _ := os.Stat(p.path)

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return DIRECTORY
	case mode.IsRegular():
		return FILE
	}

	return -1
}

func (p *Path) IsFile() bool {
	fi, _ := os.Stat(p.path)
	return fi.Mode().IsRegular()
}

func (p *Path) IsDirectory() bool {
	fi, _ := os.Stat(p.path)
	return fi.Mode().IsDir()
}

//func (p *Path) GetFileObject() file.File {
//	var f = file.File{}
//	f.SetPath(p.path)
//	return f
//}