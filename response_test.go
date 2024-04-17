package request

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func Test_parseResponse(t *testing.T) {
	type args struct {
		httpResponse *http.Response
	}
	tests := []struct {
		name         string
		args         args
		wantResponse *Response
		wantErr      bool
	}{
		{
			name: "parse response",
			args: args{
				httpResponse: &http.Response{
					Status:        "200 OK",
					StatusCode:    http.StatusOK,
					Proto:         "HTTP/1.0",
					ProtoMajor:    1,
					ProtoMinor:    0,
					Header:        map[string][]string{"Content-Type": {"application/json"}},
					Body:          io.NopCloser(strings.NewReader(`{"hello": "world"}`)),
					ContentLength: 18,
				},
			},
			wantResponse: &Response{
				StatusCode: http.StatusOK,
				Header:     map[string][]string{"Content-Type": {"application/json"}},
				RawBody:    []byte(`{"hello": "world"}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := parseResponse(tt.args.httpResponse)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("parseResponse() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestResponse_UnmarshalJSONBody(t *testing.T) {
	type fields struct {
		StatusCode int
		Header     http.Header
		RawBody    []byte
	}
	type args struct {
		val *map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal *map[string]interface{}
		wantErr bool
	}{
		{
			name: "unmarshal response json body",
			fields: fields{
				StatusCode: http.StatusOK,
				Header:     map[string][]string{"Content-Type": {"application/json"}},
				RawBody:    []byte(`{"hello": "world"}`),
			},
			args: args{
				val: &map[string]interface{}{},
			},
			wantVal: &map[string]interface{}{"hello": "world"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &Response{
				StatusCode: tt.fields.StatusCode,
				Header:     tt.fields.Header,
				RawBody:    tt.fields.RawBody,
			}
			err := resp.UnmarshalJSONBody(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSONBody() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.val, tt.wantVal) {
				t.Errorf("UnmarshalJSONBody() val = %v, wantVal %v", tt.args.val, tt.wantVal)
			}
		})
	}
}
