package k8sdeploy

import (
	"github.com/DawnBreather/go-commons/app/cicd_envsubst/playground"
	"os"
)

func Execute(){

	os.Setenv("SECRET", playground.TestEnvVar)

	opts.ReadCliOptions()

	_logger.Infof("PROCESSING")

	opts.
		GenerateKubeconfig().
		CreateK8sNamespace().
		DeployK8sManifests().
		DeployHelmCharts()

	_logger.Infof("DONE")

}