package deploy_puller

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	. "github.com/DawnBreather/go-commons/app/common_app_models"
	"github.com/DawnBreather/go-commons/executor"
	"github.com/DawnBreather/go-commons/file"
	"github.com/DawnBreather/go-commons/logger"
	path2 "github.com/DawnBreather/go-commons/path"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	_logger = logger.New()

	md5s = map[string]string{}

	workDir, _ = os.UserHomeDir()

	client = resty.New().SetTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
)

func main() {

	Execute()

}


func setupDockerConfig(){
	p := path2.Path{}
	if ! p.SetCompositePath(workDir, ".docker").Exists() {
		p.MkdirAll(0644)
	}

	p.SetCompositePath(workDir, ".docker", "config.json")
	f := file.File{}
	f.SetPath(p.GetPath())
	if ! p.Exists() {
		f.
			SetContent([]byte(`{ "credsStore": "ecr-login" }`)).
			Save()
	} else {
		f.ReadContent()
		var content map[string]interface{}
		err := json.Unmarshal(f.GetContent(), &content)
		if err != nil {
			_logger.Errorf("Unable to unmarshal { %s }: %v", f.GetPath(), err)
			return
		}
		content["credStore"] = "ecr-login"
		resContent, _ := json.Marshal(content)
		f.SetContent(resContent).Save()
	}
}

func Execute() {
	if len(os.Args) < 5 {
		_logger.Fatalf("Not enough arguments provided: 1 - Bastion URL\n 2 - Environment name\n 3 - Pull interval in seconds\n 4..N - Names of services")
	}

	workDir = workDir

	bastionUrl := os.Args[1]
	env := os.Args[2]
	pullIntervalString := os.Args[3]
	services := os.Args[4:]

	pullInterval, err := strconv.Atoi(pullIntervalString)
	if err != nil {
		_logger.Errorf("Wrong format for { Pull internval } parameter provided. Setting default value of { 5 }.")
		pullInterval = 5
	}

	_logger.Infof("Workdir: %s", workDir)
	_logger.Infof("Env: %s", env)
	_logger.Infof("Pull interval: %d", pullInterval)
	_logger.Infof("Services: %s", services)


	setupDockerConfig()

	for {

		for _, s := range services {
			md5Path := fmt.Sprintf("%s/%s", env, s)
			md5 := pullMd5(bastionUrl, env, s)
			if md5 != "" {
				if val, ok := md5s[md5Path]; ok {
					if val != md5 {
						_logger.Infof("Pulled MD5 for { %s/%s }: %s", env, s, md5)
						handleService(bastionUrl, env, s)
					}
				} else {
					handleService(bastionUrl, env, s)
				}

				md5s[md5Path] = md5
			}
		}

		time.Sleep(time.Duration(pullInterval) * time.Second)
	}
}

func handleService(bastionUrl, env, service string){
	files := pullFiles(bastionUrl, env, service)
	basePath := path2.Path{}
	basePath.SetCompositePath(workDir, env, service)
	if ! basePath.Exists() {
		basePath.MkdirAll(0644)
	}
	file.RemoveFilesRecursively(basePath.GetPath())
	for _, f := range files.ToFiles(basePath.GetPath()) {
		f.Save()
	}
	e := executor.Executor{WorkingDirectory: basePath.GetPath()}
	e.Execute("docker-compose", "pull")
	e.Execute("docker-compose", "up", "-d")
}

func pullMd5(bastionUrl, env, service string) string {

	fullUrl := fmt.Sprintf("%s/md5/%s/%s", bastionUrl, env, service)

	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Get(fullUrl)
	if err != nil {
		_logger.Errorf("Unable to pull MD5 from { %s }: %v", fullUrl, err)
		return ""
	} else {
		if resp.StatusCode() >= 400 {
			_logger.Errorf("Unable to pull MD5 from { %s }: %d", fullUrl, resp.StatusCode())
			return ""
		} else {
			md5 := string(resp.Body())
			//_logger.Infof("Pulled MD5 for { %s/%s }: %s", env, service, md5)
			return md5
		}
	}

}

func pullFiles(bastionUrl, env, service string) *FilesPayload {

	fullUrl := fmt.Sprintf("%s/files/%s/%s", bastionUrl, env, service)

	//client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Get(fullUrl)
	if err != nil {
		_logger.Errorf("Unable to pull files from { %s }: %v", fullUrl, err)
		return nil
	} else {
		if resp.StatusCode() >= 400 {
			_logger.Errorf("Unable to pull files from { %s }: %d", fullUrl, resp.StatusCode())
			return nil
		} else {
			var payload FilesPayload
			payload.FromJson(resp.Body())
			if len(payload.Files) == 0 {
				_logger.Errorf("Unable to pull files for { %s/%s }: no files in payload", env, service)
				return nil
			} else {
				var fileNames []string
				for fName, _ := range payload.Files{
					fileNames = append(fileNames, fName)
				}
				_logger.Infof("Pulled files for { %s/%s }: %s", env, service, fileNames)
				return &payload
			}
		}
	}
}