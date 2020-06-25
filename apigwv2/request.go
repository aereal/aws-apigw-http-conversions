package apigwv2

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func NewHTTPRequest(host string, apiReq *events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	ht := apiReq.RequestContext.HTTP
	reqURL := fmt.Sprintf("%s://%s%s", ht.Protocol, host, ht.Path)

	var body io.Reader
	if apiReq.Body != "" {
		body = strings.NewReader(apiReq.Body) // TODO: handle base64
	}

	req, err := http.NewRequest(ht.Method, reqURL, body)
	if err != nil {
		return nil, err
	}

	for k, v := range apiReq.Headers {
		if k == http.CanonicalHeaderKey("host") {
			continue
		}
		req.Header.Add(k, v)
	}

	return req, nil
}
