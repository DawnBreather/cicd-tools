package executor

import (
	"bytes"
	"github.com/DawnBreather/go-commons/logger"
	"os/exec"
)

var _logger = logger.New()

type Executor struct {
	WorkingDirectory string
}

func (e *Executor) Execute(executablePath string, args ...string) string{
	cmd := exec.Command(executablePath, args...)
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	cmd.Dir = e.WorkingDirectory
	out, err := cmd.Output()
	_logger.Infof("Executing { %s } with args %s", executablePath, args)
	if err != nil {
		_logger.Errorf("Exit code { %s } | stdout { %s } | stderr { %s }", err, out, stderr.String())
	} else {
		_logger.Infof("%s%s", out, stderr.String())
	}

	return string(out)
}

func Execute(executablePath string, args ...string) string{
	e := Executor{}
	return e.Execute(executablePath, args...)
}

