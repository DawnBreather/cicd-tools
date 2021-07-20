package deploy_agent

var payloadJsonExample = `{
	"files": {
		".env": "${base64_string}",
		".env.application": "${base64_string}",
		"docker-compose.yml": "${base64_string}"
	},
	"metadata": {
		"bastion_endpoint": "http://10.1.2.3:4444",
		"service_name": "service1"
	}
}`