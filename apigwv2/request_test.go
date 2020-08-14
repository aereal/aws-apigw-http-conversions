package apigwv2

import (
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestNewHTTPRequest(t *testing.T) {
	type args struct {
		host   string
		apiReq *events.APIGatewayV2HTTPRequest
	}
	tests := []struct {
		name    string
		args    args
		want    func() *http.Request
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				host: "example.com",
				apiReq: &events.APIGatewayV2HTTPRequest{
					RequestContext: events.APIGatewayV2HTTPRequestContext{
						HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
							Method:    "GET",
							Path:      "/",
							Protocol:  "HTTP/1.1",
							UserAgent: "robot/1.0.0",
						},
					},
				},
			},
			want: func() *http.Request {
				r, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
				if err != nil {
					panic(err)
				}
				return r
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHTTPRequest(tt.args.host, tt.args.apiReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotDump, err := httputil.DumpRequest(got, true)
			if err != nil {
				t.Fatal(err)
			}
			wantDump, err := httputil.DumpRequest(tt.want(), true)
			if err != nil {
				t.Fatal(err)
			}
			if string(gotDump) != string(wantDump) {
				t.Errorf("NewHTTPRequest() = %#v, want %#v", got, tt.want())
			}
		})
	}
}
