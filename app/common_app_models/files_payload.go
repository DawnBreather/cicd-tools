package common_app_models

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/DawnBreather/go-commons/file"
	"github.com/DawnBreather/go-commons/logger"
)

var _logger = logger.New()

type FilesPayload struct {
	Files               map[string]string `json:"files"`
	//ExecutionDirectives []string          `json:"execution_directives"`
}

func (fp *FilesPayload) FromJson(jsonBytes []byte) *FilesPayload{
	err := json.Unmarshal(jsonBytes, fp)
	if err != nil {
		_logger.Errorf("Unable to unmarshal JSON for FilesPayload: %v", err)
	}
	return fp
}

func (fp FilesPayload) ToJsonBytes() []byte{
	jsonBytes, _ := json.Marshal(fp)
	return jsonBytes
}
func (fp FilesPayload) ToMD5() string{
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v", fp))))
}
func (fp FilesPayload) ToFiles(basePath string) (res []*file.File) {
	for fName, fBase64 := range fp.Files {
		f := file.File{}
		f.SetPath(fmt.Sprintf("%s/%s", basePath, fName)).
			SetBase64(fBase64).
			ParseBase64ToContent()
		res = append(res, &f)
	}

	return
}