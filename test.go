package main

import (
	"fmt"
	"github.com/DawnBreather/go-commons/aws"
	"github.com/DawnBreather/go-commons/file"
	"github.com/DawnBreather/go-commons/rest"
)

func testFindFiles() {
	files := file.FindFilesByRegexpRecursively(`c:\lab\buildtools\helpers\golang\commons\go\bin\test\find_files`, `([.]env([.][a-zA-Z0-9])*)|(docker-compose.yml)`)
	for _, f := range files {
		fmt.Println(f.GetPath())
	}
}

func testBase64FilesConversion() {
	f := file.File{}
	f.SetPath(`C:\lab\buildtools\helpers\golang\commons\go\bin\test\base64\invoice.pdf`).
		ReadContent().
		ParseContentToBase64()

	b64v := f.GetBase64()

	f2 := file.File{}
	f2.SetBase64(b64v).
		ParseBase64ToContent().
		SaveTo(`C:\lab\buildtools\helpers\golang\commons\go\bin\test\base64\invoice_dst.pdf`)
}

func testReplaceEnvVarsPlaceholders() {
	f := file.File{}
	f.SetPath(`C:\lab\buildtools\helpers\golang\commons\go\bin\test\env_vars_placeholders\.env`).
		ReadContent().
		ReplaceEnvVarsPlaceholder("{{", "}}").
		SaveTo(`C:\lab\buildtools\helpers\golang\commons\go\bin\test\env_vars_placeholders\.env_res`)

	fmt.Println(string(f.GetContent()))
}

func testGetAutoscalingInstanceId() {
	a := aws.ASG{}
	a.SetName("Backend").CollectInstances()
	instances := a.GetInstances()
	for _, i := range instances {
		fmt.Println(*i.PrivateIpAddress)
	}
}

func testRest() {
	r := rest.REST{}
	r.SetEndpoint("https://3dnews.ru")
	resp := r.SubmitGET()
	fmt.Println(string(resp.Body()), resp.StatusCode())
}

func main() {
	//testFindFiles()
	//testBase64FilesConversion()
	//testReplaceEnvVarsPlaceholders()
	//testGetAutoscalingInstanceId()
	testRest()
}
