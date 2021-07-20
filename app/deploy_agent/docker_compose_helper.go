package deploy_agent

import (
	"fmt"
	path2 "github.com/DawnBreather/go-commons/path"
	"github.com/DawnBreather/go-commons/transport"
	"os"
)

func deployDockerComposeIfMissing(){
	dockerComposeExecutablePath := fmt.Sprintf("%s/%s", "/usr/local/bin", "docker-compose")
	p := path2.Path{}
	p.SetPath(dockerComposeExecutablePath)

	shouldInstall := false
	if p.Exists() {
		if ! p.IsFile() {
			shouldInstall = true
		}
	} else {
		shouldInstall = true
	}

	if shouldInstall {
		transport.DownloadFile(dockerComposeExecutablePath, config["docker_compose_download_url"].(string))
		os.Chmod(dockerComposeExecutablePath, 0755)
	}
}