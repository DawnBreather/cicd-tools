package twillio_to_telegram_webhook_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/schema"
	"net/url"
	"strings"
)

const (
	_TELEGRAM_CHAT_ID = 409143978
	_TELEGRAM_BOT_TOKEN = `1894826274:AAEuhe-yRpiIdi0pd_tcj1UcV8u9xgiF6Fg`
	_TELEGRAM_API_URL = `https://api.telegram.org/bot` + _TELEGRAM_BOT_TOKEN + "/sendMessage"

)

type SampleEvent struct {
	ID   string `json:"id"`
	Val  int    `json:"val"`
	Flag bool   `json:"flag"`
}

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//fmt.Printf("%+v", event)
	payload := getTwillioWebhookPayload(event.Body)
	messageForTelegram := fmt.Sprintf("%s\n\n%s | %s", payload.Body, payload.From, payload.FromCountry)
	sendMessageToTelegramBot(messageForTelegram)

	return events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers: map[string]string{
			"body": event.Body,
			"path": event.Path,
			"method": event.HTTPMethod,
			"Content-Type": "application/xml",
		},
		MultiValueHeaders: nil,
		Body:              `<?xml version="1.0" encoding="UTF-8"?><Response></Response>`, //event.Body,//fmt.Sprintf("%+v", event),
		IsBase64Encoded:   false,
	}, nil
}

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func sendMessageToTelegramBot(message string){

	msgObject := sendMessageReqBody{
		ChatID: _TELEGRAM_CHAT_ID,
		Text:   message,
	}
	reqBytes, _ := json.Marshal(msgObject)

	client := resty.New()
	client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(reqBytes)).
		Post(_TELEGRAM_API_URL)
}

func ExecuteLambda() {
	lambda.Start(HandleRequest)
}

const (
	twillio_test_data = `ToCountry=GB&ToState=&SmsMessageSid=SM31d88016f544ab1431e7dc6ce6b4e880&NumMedia=0&ToCity=&FromZip=&SmsSid=SM31d88016f544ab1431e7dc6ce6b4e880&FromState=&SmsStatus=received&FromCity=&Body=Hello7&FromCountry=BY&To=%2B447782827454&ToZip=&NumSegments=1&MessageSid=SM31d88016f544ab1431e7dc6ce6b4e880&AccountSid=AC3d990a55a1f67f8f168ae6217922e42d&From=%2B375445809867&ApiVersion=2010-04-01`
)

type FromTwillioWebhookPayloadStruct struct {
	Body        string `schema:"Body"`
	FromCountry string `schema:"FromCountry"`
	To          string `schema:"To"`
	From        string `schema:"From"`
}

var decoder  = schema.NewDecoder()

func getTwillioWebhookPayload(payloadRaw string) FromTwillioWebhookPayloadStruct{
	//var fromTwillioWebhookPayload FromTwillioWebhookPayloadStruct

	params := map[string]string{}

	payload, _ := url.QueryUnescape(payloadRaw)
	for _, v := range strings.Split(payload, "&") {
		first := strings.Index(v, "=")
		name := v[:first]
		value := v[first+1:]
		params[name] = value
	}


	return FromTwillioWebhookPayloadStruct{
		Body:        params["Body"],
		FromCountry: params["FromCountry"],
		To:          params["To"],
		From:        params["From"],
	}
}

func ExecuteTest(){
	res := getTwillioWebhookPayload(twillio_test_data)
	fmt.Printf(`%+v`, res)
}