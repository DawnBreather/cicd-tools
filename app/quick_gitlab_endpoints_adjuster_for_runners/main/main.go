package main

import (
	. "github.com/DawnBreather/go-commons/app/quick_gitlab_endpoints_adjuster_for_runners"
)

func main(){
	h := HostsConfig{}
	h.ParseFile().ExtractConfig().ExtractProjects()
	//projects := h.GetProjects()
	//config := h.GetConfig()
	//fmt.Printf("%v", projects)
	//fmt.Printf("%v", config)

	//h.ValidateConnectivity()
	h.InitializeSSHConnectors()

	//h.ValidateConnectivity()

	//h.RestartGitlabRunners()

	h.IdentifyConfigTomlLocation().
		ReadRemoteConfigToml().
		BackupConfigToml().
		AdjustConfigToml().
		//SaveAdjustedConfigTomlLocally().
		UploadAdjustedConfigToml().
		RestartGitlabRunners().
		CloseSSHConnections()

}