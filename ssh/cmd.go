package ssh

import (
	"fmt"
)

type cmd struct {
	envVars []string

	//cmd *goph.Cmd
	s *SSH
}

func (c *cmd) SetEnvVar(name, value string) *cmd{
	c.envVars = append (c.envVars, fmt.Sprintf("%s=%s", name, value))
	return c
}

func (c *cmd) Execute(executable string, args ...string) string{

	command, err := c.s.client.Command(executable, args...)
	if err != nil {
		_logger.Errorf("Unable to initialize command line { %s } for { %s }: %v", fmt.Sprintf("%s %v", executable, args), c.s.getSocket(), err)
		return ""
	}

	command.Env = c.envVars
	out, err := command.CombinedOutput()
	if err != nil {
		_logger.Errorf("Error executing { %s } on { %s }: { %v } { %s }", command.String(), c.s.getSocket(), err, string(out))
	}

	return string(out)
}

func (c *cmd) ExecuteBash(command string) string{
	resCommand := fmt.Sprintf(`"%s"`, command)
	out := c.Execute("/bin/bash", "-c", resCommand)

	return out
}

func (c *cmd) ExecuteSudoBash(command string) string{
	resCommand := fmt.Sprintf(`"%s"`, command)
	out := c.Execute("sudo", "/bin/bash", "-c", resCommand)

	return out
}

//func (c *cmd) ExecuteSudo(command string) string{
//	resCommand := fmt.Sprintf(`"%s"`, command)
//	out := c.Execute("/bin/bash", "-c", resCommand)
//
//	return out
//}