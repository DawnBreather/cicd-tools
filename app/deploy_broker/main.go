package deploy_broker

import (
	"encoding/json"
	"fmt"
	. "github.com/DawnBreather/go-commons/api_server"
	. "github.com/DawnBreather/go-commons/app/common_app_models"
	"github.com/DawnBreather/go-commons/file"
	"github.com/DawnBreather/go-commons/logger"
	path2 "github.com/DawnBreather/go-commons/path"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)
var basePath, _ = os.UserHomeDir()

var _logger = logger.New()
var md5s = map[string]string{}

func main(){
	Execute()
}

func Execute(){
	_logger.Infof("Base path: %s", basePath)

	var server = ApiServer{}
	server.
		Initialize(nil).
		Post("/files/{env}/{service}", handlePostFiles).
		Get("/files/{env}/{service}", handleGetFiles).
		Get("/md5/{env}/{service}", handleGetMD5).
		RunSslSelfSigned(os.Getenv("LISTEN_TO"), []string{os.Getenv("DOMAIN_NAME")})
		//Run(os.Getenv("LISTEN_TO"))
}

func compilePayload(env, service string) (payload *FilesPayload, fileNames []string) {
	var resFiles  = map[string]string{}

	p := path2.Path{}
	p.SetCompositePath(basePath, env, service)
	if p.Exists() && p.IsDirectory() {
		files := file.FindFilesRecursively(p.GetPath())
		for _, f := range files {
			fName := f.GetFileName()
			fBase64 := f.ReadContent().ParseContentToBase64().GetBase64()
			resFiles[fName] = fBase64
			fileNames = append(fileNames, fName)
		}

		payload = &FilesPayload{Files: resFiles}

		return payload, fileNames
	}

	return nil, nil
}

func handleGetFiles(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	env := vars["env"]
	service := vars["service"]

	payload, fileNames := compilePayload(env, service)
	_logger.Infof("Serving files from { %s/%s/%s }: %s", basePath, env, service, fileNames)
	if payload != nil {
		RespondJSON(w, http.StatusOK, payload)
	} else {
		RespondJSON(w, http.StatusNotFound, nil)
	}
}

func handleGetMD5(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	env := vars["env"]
	service := vars["service"]

	md5Path := fmt.Sprintf("%s/%s", env, service)
	if val, ok := md5s[md5Path]; ok {
		RespondPlaintext(w, http.StatusOK, val)
	} else {
		payload, fileNames := compilePayload(env, service)
		_logger.Infof("Serving MD5 for { %s/%s/%s } containing files %s", basePath, env, service, fileNames)
		if payload != nil {
			md5s[md5Path] = payload.ToMD5()
			RespondPlaintext(w, http.StatusOK, md5s[md5Path])
		} else {
			RespondPlaintext(w, http.StatusNotFound, "")
		}
	}

}


func handlePostFiles(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	env := vars["env"]
	service := vars["service"]


	var payload FilesPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	md5Path := fmt.Sprintf("%s/%s", env, service)
	md5s[md5Path] = payload.ToMD5()

	// Open Handler Body

	files := payload.Files

	dstFolderPath := path2.Path{}
	dstFolderPath.
		SetCompositePath(basePath, env, service).
		MkdirAll(0644)

	// Saving
	file.RemoveFilesRecursively(dstFolderPath.GetPath())
	var fileNames []string
	for fileName, fileContentBase64 := range files {
		fileNames = append(fileNames, fileName)
		f := file.File{}
		f.SetBase64(fileContentBase64).
			ParseBase64ToContent().
			ParseContentToMd5().
			SaveTo(dstFolderPath.GetPath() + "/" + fileName)
	}
	_logger.Infof("Saving files into { %s }: %s", dstFolderPath.GetPath(), fileNames)

	// Close Handler Body
	RespondJSON(w, http.StatusCreated, nil)
}