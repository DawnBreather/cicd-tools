package deploy_agent

import (
	"fmt"
	. "github.com/DawnBreather/go-commons/api_server"
	. "github.com/DawnBreather/go-commons/executor"
	"os"
	"strings"
)

var homeDir, _ = os.UserHomeDir()

var config = map[string]interface{}{
	"base_path": homeDir,
	"docker_compose_download_url": fmt.Sprintf("https://github.com/docker/compose/releases/download/1.29.1/docker-compose-%s-%s", strings.TrimSuffix(Execute("uname", "-s"), "\n"), strings.TrimSuffix(Execute("uname", "-m"), "\n")),
	"listen_to": os.Getenv("LISTEN_TO"),
}

func main(){
	var server = ApiServer{}
	deployDockerComposeIfMissing()
	server.
		Initialize(nil).
		Post("/deploy", handleDeployment).
		Run(config["listen_to"].(string))
}

func Execut() {

	homeDir, _ := os.UserHomeDir()

	var config = map[string]interface{}{
		"base_path":                   homeDir,
		"docker_compose_download_url": fmt.Sprintf("https://github.com/docker/compose/releases/download/1.29.1/docker-compose-%s-%s", Execute("uname", "-s"), Execute("uname", "-m")),
		"listen_to":                   os.Getenv("LISTEN_TO"),
	}

	var server = ApiServer{}
	deployDockerComposeIfMissing()
	server.
		Initialize(nil).
		Post("/deploy", handleDeployment).
		Run(config["listen_to"].(string))
}




