package quick_gitlab_endpoints_adjuster_for_runners

import (
	"fmt"
	"github.com/DawnBreather/go-commons/file"
	"github.com/DawnBreather/go-commons/logger"
	path2 "github.com/DawnBreather/go-commons/path"
	"github.com/DawnBreather/go-commons/ssh"
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"
	"time"
)


const (
	_hostsFilePath = "hosts.yml"
	_gitlabRunnerConfigName = "config.toml"
	_gitlabRunnerConfigRemotePath = "/etc/gitlab-runner/config.toml"
)

var _logger = logger.New()

type config struct {
	sshKeysBasedir string
}

type project struct {
	name string `yaml:"-"`
	abcloudz bool `yaml:"abcloudz"`

	hosts map[string]*host `yaml:"hosts"`
}

type host struct {
	host           string `yaml:"host"`
	port           int    `yaml:"port"`
	username       string `yaml:"username"`
	password       string `yaml:"password"`
	sshKeyName     string `yaml:"ssh_key_name"`
	sshKeyPassword string `yaml:"ssh_key_password"`
	path           string `yaml:"path"`
}

func (h *host) SetPath(path string) {
	h.path = path
}



type HostsConfig struct {
	generic  map[string]interface{}
	config   config
	projects map[string]project

	gitlabRunnersConfigs map[string]map[string]string
	gitlabRunnersResConfigs map[string]map[string]string
	sshConnectors map[string]map[string] *ssh.SSH
}
func (hc *HostsConfig) GetProjects() map[string]project{
	return hc.projects
}

func (hc *HostsConfig) GetConfig() config{
	return hc.config
}

func (hc *HostsConfig) ParseFile() *HostsConfig {
	path := path2.Path{}
	path.SetPath(_hostsFilePath)

	if path.Exists() {
		hostsFile := file.File{}
		content := hostsFile.SetPath(_hostsFilePath).ReadContent().GetContent()
		err := yaml.Unmarshal([]byte(content), &hc.generic)
		if err != nil {
			_logger.Fatalf("Unable to unmarshal HostsConfig file: %v", err)
		}
	} else {
		_logger.Fatalf("Unable to locate HostsConfig file { %s }", _hostsFilePath)
	}

	return hc
}

func (hc *HostsConfig) InitializeSSHConnectors() *HostsConfig{
	for _, p := range hc.projects {
		for _, h := range p.hosts {
			if hc.sshConnectors == nil {
				hc.sshConnectors = map[string]map[string] *ssh.SSH{}
			}
			if hc.sshConnectors[p.name] == nil {
				hc.sshConnectors[p.name] = map[string]*ssh.SSH{}
			}
			hc.sshConnectors[p.name][h.host] = sshToHost(hc, *h)
		}
	}

	return hc
}

func (hc *HostsConfig) IdentifyConfigTomlLocation() *HostsConfig{
	r := regexp.MustCompile(`--config .*config.toml`)
	for project, hostsMap := range hc.sshConnectors {
		for host, ssh := range hostsMap {
			ssh.Connect()
			out := ssh.CMD.ExecuteSudoBash(`ps aux | grep gitlab-runner`)
			configPathRaw := r.FindString(out)
			configPath := strings.ReplaceAll(configPathRaw, "--config ", "")
			if configPath != "" {
				hc.projects[project].hosts[host].SetPath(configPath)
				_logger.Infof("Identifying location of config for { %s / %s }: %s", project, host, configPath)
			} else {
				hc.projects[project].hosts[host].SetPath(_gitlabRunnerConfigRemotePath)
				_logger.Infof("Identifying location of config for { %s / %s } did not succeed, falling back to { %s }", project, host, _gitlabRunnerConfigRemotePath)
			}
		}
	}

	return hc
}

func (hc *HostsConfig) ExtractConfig() *HostsConfig{
	 hc.config.sshKeysBasedir = hc.generic["ssh"].(map[interface{}]interface{})["keys"].(map[interface{}]interface{})["basedir"].(string)
	 return hc
}

func (hc *HostsConfig) ExtractProjects() *HostsConfig{
	for k, v := range hc.generic{
		if k != "ssh" {
			p := project{
				name:     k,
				abcloudz: v.(map[interface{}]interface{})["abcloudz"].(bool),
				hosts: map[string]*host{},
			}

			for _, hostMap := range v.(map[interface{}]interface{})["hosts"].([]interface{}){

				var _host = hostMap.(map[interface{}]interface{})["host"]
				var _port = hostMap.(map[interface{}]interface{})["port"]
				var _username = hostMap.(map[interface{}]interface{})["username"]
				var _password = hostMap.(map[interface{}]interface{})["password"]
				var _sshKeyName =     hostMap.(map[interface{}]interface{})["ssh_key_name"]
				var _sshKeyPassword = hostMap.(map[interface{}]interface{})["ssh_key_password"]
				var _path = hostMap.(map[interface{}]interface{})["path"]

				hst := host{}

				if _host != nil {
					hst.host = _host.(string)
				}
				if _port != nil {
					hst.port = _port.(int)
				}
				if _username != nil {
					hst.username = _username.(string)
				}
				if _password != nil {
					hst.password = _password.(string)
				}
				if _sshKeyName != nil {
					hst.sshKeyName = _sshKeyName.(string)
				}
				if _sshKeyPassword != nil {
					hst.sshKeyPassword = _sshKeyPassword.(string)
				}
				if _path != nil {
					hst.path = _path.(string)
				}

				p.hosts[hst.host] = &hst
			}

			p.name = k
			if hc.projects == nil {
				hc.projects = map[string]project{}
			}
			hc.projects[k] = p
		}
	}

	return hc
}

func (hc *HostsConfig) AdjustConfigToml() *HostsConfig{
	for project, hostsMap := range hc.gitlabRunnersConfigs {
		for host, configContent := range hostsMap {
			if hc.gitlabRunnersResConfigs == nil {
				hc.gitlabRunnersResConfigs = map[string]map[string]string{}
			}
			if hc.gitlabRunnersResConfigs[project] == nil {
				hc.gitlabRunnersResConfigs[project] = map[string]string{}
			}

			hc.gitlabRunnersResConfigs[project][host] = strings.ReplaceAll(configContent, "https://gitlab.kharkov.dbbest.com/", "https://gitlab.abcloudz.com/")

			adjustmentsNeeded := hc.gitlabRunnersResConfigs[project][host] != hc.gitlabRunnersConfigs[project][host]
			var adjustmentsNeededString = "YES"
			if !adjustmentsNeeded {
				adjustmentsNeededString = "NO"
			 }
			_logger.Infof("Adjustments for { %s / %s } required: %s", project, host, adjustmentsNeededString)
		}
	}

	return hc
}

func (hc *HostsConfig) RestartGitlabRunners() *HostsConfig{
	for project, hostsMap := range hc.sshConnectors {
		for host, ssh := range hostsMap {
			adjustmentsNeeded := hc.gitlabRunnersResConfigs[project][host] != hc.gitlabRunnersConfigs[project][host]
			if adjustmentsNeeded {
				ssh.Connect()
				var out string

				out = ssh.CMD.ExecuteSudoBash("find /usr -type f -name gitlab-runner | xargs -I {} /bin/bash -c '{} restart'")

				//out = ssh.CMD.ExecuteSudoBash("which gitlab-runner")
				//if out != "" {
				//	out = ssh.CMD.ExecuteSudoBash("gitlab-runner restart")
				//}
				//
				//out = ssh.CMD.ExecuteSudoBash("which /usr/local/bin/gitlab-runner")
				//if out != "" {
				//	out = ssh.CMD.ExecuteSudoBash("/usr/local/bin/gitlab-runner restart")
				//}
				//
				//out = ssh.CMD.ExecuteSudoBash("which /usr/local/sbin/gitlab-runner")
				//if out != "" {
				//	out = ssh.CMD.ExecuteSudoBash("/usr/local/sbin/gitlab-runner restart")
				//}

				//ssh.Close()
				_logger.Infof("Restarting runner for { %s / %s }: %s", project, host, out)
			}
		}
	}

	return hc
}

func (hc *HostsConfig) SaveAdjustedConfigTomlLocally() *HostsConfig{
	file := file.File{}
	for project, hostsMap := range hc.gitlabRunnersResConfigs {
		for host, resConfigContent := range hostsMap {
			adjustmentsNeeded := hc.gitlabRunnersResConfigs[project][host] != hc.gitlabRunnersConfigs[project][host]
			if adjustmentsNeeded {
				file.
					SetCompositePath(project, host, _gitlabRunnerConfigName).
					SetContent([]byte(resConfigContent)).
					Save()
				_logger.Infof("Saving adjusted config for { %s / %s } into { %s%s }", project, host, file.GetBaseDir(), file.GetFileName())
			}
		}
	}

	return hc
}

func (hc *HostsConfig) UploadAdjustedConfigToml() *HostsConfig{
	path := path2.Path{}
	for project, hostsMap := range hc.gitlabRunnersResConfigs {
		for host, _ := range hostsMap {

			adjustmentsNeeded := hc.gitlabRunnersResConfigs[project][host] != hc.gitlabRunnersConfigs[project][host]
			if adjustmentsNeeded {

				path.SetCompositePath(project, host, _gitlabRunnerConfigName)

				ssh := hc.sshConnectors[project][host]
				_logger.Infof("Uploading adjusted config toml for { %s / %s } from local { %s } to remote { %s }", project, host, path.GetPath(), _gitlabRunnerConfigRemotePath)

				ssh.Connect()

				if hc.projects[project].hosts[host].path == "" {
					ssh.UploadLocalFileOverSSHWithSudo(path.GetPath(), _gitlabRunnerConfigRemotePath)
				} else {
					ssh.UploadLocalFileOverSSHWithSudo(path.GetPath(), hc.projects[project].hosts[host].path)
				}
			}

			//ssh.Close()
		}
	}

	return hc
}

func (hc *HostsConfig) BackupConfigToml() *HostsConfig{

	p := path2.Path{}
	f := file.File{}

	for project, hostsMap := range hc.gitlabRunnersConfigs {
		for host, configContent := range hostsMap {
			p.
				SetCompositePath(project, host).
				MkdirAll(0644)
			_logger.Infof("Backing up config.toml for { %s / %s }", project, host)
			f.
				SetCompositePath(project, host, fmt.Sprintf("%s.%s", _gitlabRunnerConfigName, time.Now().Format("20060102150405"))).
				SetContent([]byte(configContent)).
				Save()
		}
	}

	return hc
}

func (hc *HostsConfig) ReadRemoteConfigToml() *HostsConfig{

	for project, hostsMap := range hc.sshConnectors {
		for host, ssh := range hostsMap {

			_logger.Infof("Reading remote config.toml for { %s / %s }", project, host)
			ssh.Connect()

			var content string
			if hc.projects[project].hosts[host].path == "" {
				content = ssh.CMD.ExecuteSudoBash("cat " + _gitlabRunnerConfigRemotePath)
			} else {
				content = ssh.CMD.ExecuteSudoBash("cat " + hc.projects[project].hosts[host].path)
			}
			if hc.gitlabRunnersConfigs == nil {
				hc.gitlabRunnersConfigs = map[string]map[string]string{}
			}
			if hc.gitlabRunnersConfigs[project] == nil {
				hc.gitlabRunnersConfigs[project] = map[string]string{}
			}
			hc.gitlabRunnersConfigs[project][fmt.Sprintf("%s", host)] = content

			//ssh.Close()

		}
	}

	return hc
}

func (hc *HostsConfig) ValidateConnectivity() *HostsConfig{
	for project, hostsMap := range hc.sshConnectors {
		for host, ssh := range hostsMap {

			_logger.Infof("Validating connectivity to { %s / %s }", project, host)

			if hc.projects[project].hosts[host].password == "" {
				if ssh.CanConnectWithKey() {
					_logger.Infof("Connection established { key }")
				} else {
					_logger.Infof("Connection failed { key }")
				}
			} else {
				if ssh.CanConnectWithPassword() {
					_logger.Infof("Connection established { password }")
				} else {
					_logger.Infof("connection failed { password }")
				}
			}
		}
	}
	return hc
}

func (hc *HostsConfig) CloseSSHConnections() *HostsConfig {
	for project, hostsMap := range hc.sshConnectors {
		for host, ssh := range hostsMap {
			_logger.Infof("Closing SSH connection for { %s / %s }", project, host)
			ssh.Close()
		}
	}

	return hc
}

func sshToHost(hc *HostsConfig, h host) *ssh.SSH{
	ssh := ssh.SSH{}
	return ssh.
		SetHost(h.host).
		SetPort(h.port).
		SetKeyPassword(h.sshKeyPassword).
		SetKeyPath(fmt.Sprintf("%s/%s", hc.config.sshKeysBasedir, h.sshKeyName)).
		SetUsername(h.username).
		SetPassword(h.password)
}