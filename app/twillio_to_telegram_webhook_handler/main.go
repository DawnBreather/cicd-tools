package twillio_to_telegram_webhook_handler

import (
	"fmt"
	. "github.com/DawnBreather/go-commons/api_server"
	"github.com/DawnBreather/go-commons/logger"
	"io/ioutil"
	"net/http"
	"os"
)

var _logger = logger.New()

func Execute(){
	_logger.Infof("Listening to %s", os.Getenv(`LISTEN_TO`))
	var server = ApiServer{}
	server.
		Initialize(nil).
		Post("/sms", handleSms).
		Get("/sms", handleSms).
		Run(os.Getenv("LISTEN_TO"))
}

func handleSms(w http.ResponseWriter, r *http.Request) {

	err := r.Body.Close()
	if err != nil {
		_logger.Errorf("Unable to close the body of SMS webhook request")
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_logger.Errorf("Unable to read body of SMS webhook request")
	}

	fmt.Println(string(bodyBytes))
}