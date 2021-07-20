package deploy_invoker

import (
	"crypto/tls"
	"fmt"
	. "github.com/DawnBreather/go-commons/app/common_app_models"
	"github.com/DawnBreather/go-commons/env_var"
	"github.com/DawnBreather/go-commons/file"
	"github.com/DawnBreather/go-commons/logger"
	path2 "github.com/DawnBreather/go-commons/path"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
)

//type deployPayload struct{
//	Files map[string]string `json:"files"`
//	Metadata map[string]interface{} `json:"metadata"`
//}

var (
	_logger = logger.New()

	client = resty.New().SetTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
)

func main(){
	Execute()
}

func Execute(){
	if len(os.Args) < 3 {
		_logger.Fatalf("Not enough arguments provided: %d", len(os.Args))
	}

	var files = map[string]string{}
	var service string
	var env string
	var brokerUrl string


	f := file.File{}
	noErrors := true
	for i, arg := range os.Args{
		if i == 0{
			continue
		}
		if i == 1 {
			brokerUrl = arg
		}
		if i == 2 {
			env = arg
		}
		if i == 3 {
			service = arg
		}
		if i > 3 {
			p := path2.Path{}
			p.SetPath(arg)
			if ! p.Exists() {
				noErrors = false
				_logger.Errorf("Provided file path { %s } does not exist", arg)
			}
			if ! p.IsFile() {
				noErrors = false
				_logger.Errorf("Provided file path { %s } is not actually a file", arg)
			}
			if ! noErrors {
				os.Exit(1)
			}

			b64 := f.SetPath(arg).ReadContent().ParseContentToBase64().GetBase64()
			fName := f.GetFileName()

			files[fName] = b64
		}
	}

	files["commit_sha"] = f.SetPath("commit_sha").SetContent([]byte(getCommitShortSha())).ParseContentToBase64().GetBase64()

	payload := FilesPayload{
		Files:    files,
	}

	var fileNames []string
	for k, _ := range files {
		fileNames = append(fileNames, k)
	}

	_logger.Infof("Submitting files to the broker: %v", fileNames)

	fullUrl := fmt.Sprintf("%s/files/%s/%s", brokerUrl, env, service)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload.ToJsonBytes()).
		Post(fullUrl)

	if err != nil {
		_logger.Errorf("Error posting deployment request to { %s }: %v", fullUrl, err)
		_logger.Errorf("Response: %v", resp)
		os.Exit(1)
	}
}

func getCommitShortSha() string{
	sha := env_var.EnvVar{}

	switch {
	case sha.SetName("CI_COMMIT_SHORT_SHA").IsExist():
			return sha.Value()
	default:
		return ""
	}

}

















