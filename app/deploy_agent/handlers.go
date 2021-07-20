package deploy_agent

import (
	"encoding/json"
	. "github.com/DawnBreather/go-commons/api_server"
	. "github.com/DawnBreather/go-commons/executor"
	"github.com/DawnBreather/go-commons/file"
	path2 "github.com/DawnBreather/go-commons/path"
	"net/http"
	"os"
)

func handleDeployment(w http.ResponseWriter, r *http.Request) {
	var deploymentPayload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&deploymentPayload); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// Open Handler Body

	files := deploymentPayload["files"].(map[string]interface{})
	metadata := deploymentPayload["metadata"].(map[string]interface{})

	//basePath := metadata["base_path"].(string)
	basePath, _ := os.UserHomeDir()
	serviceName := metadata["service_name"].(string)

	dstFolderPath := path2.Path{}
	dstFolderPath.
		SetCompositePath(basePath, serviceName).
		MkdirAll(0644)

	for fileName, fileContentBase64 := range files {
		f := file.File{}
		f.SetBase64(fileContentBase64.(string)).
			ParseBase64ToContent().
			SaveTo(dstFolderPath.GetPath() + "/" + fileName)
	}

	e := Executor{WorkingDirectory: dstFolderPath.GetPath()}
	e.Execute("docker-compose", "pull")
	e.Execute("docker-compose", "up", "-d")

	// Close Handler Body

	RespondJSON(w, http.StatusCreated, nil)
}