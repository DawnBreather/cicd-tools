package k8sdeploy

import (
	"fmt"
	"github.com/DawnBreather/go-commons/executor"
	"github.com/DawnBreather/go-commons/file"
	"github.com/DawnBreather/go-commons/logger"
	"github.com/iancoleman/strcase"
	"github.com/jessevdk/go-flags"
	"strings"
)

var _logger = logger.New()
/*var homeDir, _ = os.UserHomeDir()*/
var e = executor.Executor{/*WorkingDirectory: homeDir*/}

var opts = optsT{}

type optsT struct {

	EksClusterName string `short:"c" long:"eks-cluster-name" description:"EKS cluster name" required:"true"`
	K8sManifests []string `short:"m" long:"k8s-manifest" description:"K8s manifest path" required:"false"`
	HelmCharts map[string]string `short:"t" long:"helm-chart" description:"K8s helm name and path" required:"false"`
	HelmChartValues map[string]string `short:"v" long:"helm-chart-values" description:"K8s helm values" required:"false"`
	HelmChartSetValueOverrides map[string]string `short:"s" long:"helm-chart-set-value" description:"K8s helm value overrides" required:"false"`
	Namespace string `short:"n" long:"namespace" description:"K8s namespace" required:"true"`

}

func (o *optsT) ReadCliOptions() *optsT{
	//_logger.SetLevel(logrus.DebugLevel)
	var err error
	_, err = flags.Parse(o)

	if err != nil {
		_logger.Fatalf("Unable to read cli options: %v", err)
	}

	return o
}


func (o *optsT) GenerateKubeconfig() *optsT{
	_logger.Infof("GENERATING KUBECONFIG")
	arguments := []string{ "eks", "update-kubeconfig", "--name", o.EksClusterName }
	//_logger.Debugf("-> { %s }", strings.Join(append([]string{"aws"}, arguments...), " "))
	e.Execute("aws", arguments...)

	return o
}

func (o *optsT) CreateK8sNamespace() *optsT{
	_logger.Infof("CREATING K8S NAMESPACE")
	//_logger.Infof("-> %s", o.Namespace)
	arguments := []string{ "-c", fmt.Sprintf("kubectl create namespace %s --dry-run -o yaml | kubectl apply -f -", o.Namespace) }
	//_logger.Debugf("-> { %s }", strings.Join(append([]string{"/bin/sh"}, arguments...), " "))
	e.Execute("/bin/sh", arguments...)

	return o
}

func (o *optsT) DeployK8sManifests() *optsT{
	_logger.Infof("DEPLOYING K8S MANIFESTS")
	for _, m := range o.K8sManifests {
		//_logger.Infof("-> %s", m)
		arguments := []string{ "apply", "-f", m, "-n", o.Namespace }
		//_logger.Debugf("-> { %s }", strings.Join(append([]string{"kubectl"}, arguments...), " "))
		e.Execute("kubectl", arguments... )
	}

	return o
}

func (o *optsT) DeployHelmCharts() *optsT{
	_logger.Infof("DEPLOYING HELM CHARTS")
	for name, path := range o.HelmCharts {
		//_logger.Infof("-> %s:%s", name, path)
		executionArgs := []string{"upgrade", "--install", name, path, "--namespace", o.Namespace}
		if val, ok := o.HelmChartValues[name]; ok {

			val = strcase.ToScreamingSnake(val)
			executionArgs = append(executionArgs, []string{"--values", val}...)
			f := file.File{}
			f.SetPath(val).
				FillContentFromEnvironmentVariable(val).
				Save()
		}
		if val, ok := o.HelmChartSetValueOverrides[name]; ok {
			//_logger.Debugf("-> { %v }", val)
			for _, s := range strings.Split(val, ";"){
				//_logger.Debugf("-> { %v }", []string{"--set", s})
				executionArgs = append(executionArgs, []string{"--set", s}...)
			}
		}
		//_logger.Debugf("-> { %s }", strings.Join(append([]string{"helm"}, executionArgs...), " "))
		e.Execute("helm", executionArgs...)
	}

	return o
}