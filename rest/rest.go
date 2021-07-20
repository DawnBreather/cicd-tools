package rest

import (
	"github.com/DawnBreather/go-commons/logger"
	"github.com/go-resty/resty/v2"
)

var _logger = logger.New()

type REST struct {
	endpoint string
	payload  string
}

func (r *REST) SetEndpoint(url string) *REST {
	r.endpoint = url
	return r
}

func (r *REST) SetPayload(payload string) *REST {
	r.payload = payload
	return r
}

func (r *REST) SubmitPOST() *resty.Response {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(r.payload)).
		Post(r.endpoint)
	if err != nil {
		_logger.Errorf("Unable to submit POST request to endpoint {%s}: %v", r.endpoint, err)
	}

	return resp
}

func (r *REST) SubmitGET() *resty.Response {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(r.endpoint)
	if err != nil {
		_logger.Errorf("Unable to submit GET request to endpoint { %s }: %v", r.endpoint, err)
	}

	return resp
}
