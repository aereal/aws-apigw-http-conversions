package apigwv2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func TestFromHTTPResponse(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		handler http.Handler
		want    *events.APIGatewayV2HTTPResponse
		wantErr bool
	}{
		{
			name: "OK",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("date", now.Format(http.TimeFormat))
				fmt.Fprintln(w, "OK")
			}),
			want: &events.APIGatewayV2HTTPResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Length": "3",
					"Content-Type":   "text/plain; charset=utf-8",
					"Date":           now.Format(http.TimeFormat),
				},
				MultiValueHeaders: map[string][]string{},
				Body:              "OK\n",
			},
			wantErr: false,
		},
		{
			name: "header",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("date", now.Format(http.TimeFormat))
				w.Header().Set("cache-control", "private")
				fmt.Fprintln(w, "OK")
			}),
			want: &events.APIGatewayV2HTTPResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Length": "3",
					"Content-Type":   "text/plain; charset=utf-8",
					"Cache-Control":  "private",
					"Date":           now.Format(http.TimeFormat),
				},
				MultiValueHeaders: map[string][]string{},
				Body:              "OK\n",
			},
			wantErr: false,
		},
		{
			name: "multiple header",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("date", now.Format(http.TimeFormat))
				w.Header().Add("cache-control", "private")
				w.Header().Add("cache-control", "no-cache")
				fmt.Fprintln(w, "OK")
			}),
			want: &events.APIGatewayV2HTTPResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Length": "3",
					"Content-Type":   "text/plain; charset=utf-8",
					"Date":           now.Format(http.TimeFormat),
				},
				MultiValueHeaders: map[string][]string{
					"Cache-Control": {"private", "no-cache"},
				},
				Body: "OK\n",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()
			resp, err := srv.Client().Get(srv.URL)
			if err != nil {
				t.Fatalf("! %v", err)
			}

			got, err := FromHTTPResponse(resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromHTTPResponse() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromHTTPResponse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
