package main

import (
	"github.com/DawnBreather/go-commons/app/twillio_to_telegram_webhook_handler"
)

func main(){
	//cicd_envsubst.Execute()
	//deploy_agent.Execut()
	//deploy_invoker.Execute()
	//deploy_broker.Execute()
	//deploy_puller.Execute()
	twillio_to_telegram_webhook_handler.ExecuteLambda()
}