package cicd_envsubst

import (
	"github.com/DawnBreather/go-commons/file"
	"github.com/DawnBreather/go-commons/logger"
	pth "github.com/DawnBreather/go-commons/path"
	"os"
	"strings"
)

var _logger = logger.New()

//func main(){
//	Execute()
//}

func Execute(){
	_logger.Infof("Replacing environment variables in the following locations: { %v }", strings.Join(os.Args[1:], ", "))
	_logger.Infof("PROCESSING")
	for i, arg := range os.Args{
		if i > 0 {
			var p = pth.Path{}
			p.SetPath(arg)

			if p.Exists(){
				var f = file.File{}
				if p.IsFile() {
					f.SetPath(arg)
					_logger.Infof("-> %s", f.GetPath())
					f.ReadContent().ReplaceEnvVarsPlaceholder("{{", "}}").Save()
				}
				if p.IsDirectory(){
					files := file.FindFilesRecursively(arg)
					for _, file := range files {
						var name = file.GetPath()
						f.SetPath(name)
						_logger.Infof("-> %s", f.GetPath())
						f.ReadContent().ReplaceEnvVarsPlaceholder("{{", "}}").Save()
					}
				}
			}

		}
	}
	_logger.Infof("DONE")
}
