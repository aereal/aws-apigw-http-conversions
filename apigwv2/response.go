package apigwv2

import (
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func FromHTTPResponse(resp *http.Response) (*events.APIGatewayV2HTTPResponse, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	gwresp := &events.APIGatewayV2HTTPResponse{
		Headers:           map[string]string{},
		MultiValueHeaders: map[string][]string{},
		StatusCode:        resp.StatusCode,
		Body:              string(b),
	}
	return gwresp, nil
}
