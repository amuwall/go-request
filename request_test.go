package request

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewRequest(t *testing.T) {
	type args struct {
		method  string
		path    string
		options []RequestOption
	}
	tests := []struct {
		name        string
		args        args
		wantRequest *Request
		wantErr     bool
	}{
		{
			name: "new request",
			args: args{
				method: http.MethodGet,
				path:   "/api/test",
			},
			wantRequest: &Request{
				Method:      http.MethodGet,
				Path:        "/api/test",
				Host:        "",
				Headers:     http.Header{},
				QueryParams: nil,
				BodyParams:  nil,
			},
			wantErr: false,
		},
		{
			name: "new request with option to set content type in header",
			args: args{
				method: http.MethodPost,
				path:   "/api/test/json",
				options: []RequestOption{
					func(request *Request) (err error) {
						request.Headers.Add(contentTypeHeader, contentTypeJsonWithUTF8)
						return
					},
				},
			},
			wantRequest: &Request{
				Method:      http.MethodPost,
				Path:        "/api/test/json",
				Host:        "",
				Headers:     http.Header{contentTypeHeader: []string{contentTypeJsonWithUTF8}},
				QueryParams: nil,
				BodyParams:  nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRequest, err := NewRequest(tt.args.method, tt.args.path, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRequest, tt.wantRequest) {
				t.Errorf("NewRequest() gotRequest = %v, want %v", gotRequest, tt.wantRequest)
			}
		})
	}
}

func TestWithHost(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "init request with host",
			args: args{
				host: "example.com",
			},
			want: "example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := NewRequest(http.MethodGet, "/api/test", WithHost(tt.args.host))
			if request.Host != tt.want {
				t.Errorf("WithHost() = %v, want %v", request.Host, tt.want)
			}
		})
	}
}

func TestWithHeaders(t *testing.T) {
	type args struct {
		headers map[string]string
	}
	tests := []struct {
		name string
		args args
		want http.Header
	}{
		{
			name: "init request with headers",
			args: args{
				headers: map[string]string{
					"Test-Header-1": "test-value-1",
					"Test-Header-2": "test-value-2",
				},
			},
			want: http.Header{
				"Test-Header-1": []string{"test-value-1"},
				"Test-Header-2": []string{"test-value-2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := NewRequest(http.MethodGet, "/api/test", WithHeaders(tt.args.headers))
			if !reflect.DeepEqual(request.Headers, tt.want) {
				t.Errorf("WithHeaders() = %v, want %v", request.Headers, tt.want)
			}
		})
	}
}

func TestWithQueryParams(t *testing.T) {
	type args struct {
		queryParams QueryParams
	}
	tests := []struct {
		name string
		args args
		want QueryParams
	}{
		{
			name: "init request with query params",
			args: args{
				queryParams: QueryParams{
					"key": []string{"value"},
				},
			},
			want: QueryParams{
				"key": []string{"value"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := NewRequest(http.MethodGet, "/api/test", WithQueryParams(tt.args.queryParams))
			if !reflect.DeepEqual(request.QueryParams, tt.want) {
				t.Errorf("WithQueryParams() = %v, want %v", request.QueryParams, tt.want)
			}
		})
	}
}

func TestWithBodyParams(t *testing.T) {
	type args struct {
		bodyParams BodyParams
	}
	tests := []struct {
		name string
		args args
		want BodyParams
	}{
		{
			name: "init request with body params",
			args: args{
				bodyParams: NewJsonBodyParams(map[string]string{
					"test": "test-value",
				}),
			},
			want: NewJsonBodyParams(map[string]string{
				"test": "test-value",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := NewRequest(http.MethodGet, "/api/test", WithBodyParams(tt.args.bodyParams))
			if !reflect.DeepEqual(request.BodyParams, tt.want) {
				t.Errorf("WithBodyParams() = %v, want %v", request.BodyParams, tt.want)
			}
		})
	}
}

func TestRequest_build(t *testing.T) {
	type fields struct {
		Method      string
		Path        string
		Host        string
		Headers     http.Header
		QueryParams QueryParams
		BodyParams  BodyParams
	}
	type args struct {
		baseURL string
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		wantHttpRequest         *http.Request
		wantHttpRequestBodyData []byte
		wantErr                 bool
	}{
		{
			name: "build request",
			fields: fields{
				Method: http.MethodPost,
				Path:   "/api/test",
				Host:   "example.com",
				Headers: http.Header{
					"content-type": []string{"application/json; charset=utf-8"},
				},
				QueryParams: NewQueryParams(map[string]string{"test": "value"}),
				BodyParams:  NewJsonBodyParams(map[string]interface{}{"hello": "world"}),
			},
			args: args{
				baseURL: "https://127.0.0.1",
			},
			wantHttpRequest: &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Scheme:   "https",
					Host:     "127.0.0.1",
					Path:     "/api/test",
					RawQuery: "test=value",
				},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"content-type": []string{"application/json; charset=utf-8"},
				},
				ContentLength: 17,
				Host:          "example.com",
			},
			wantHttpRequestBodyData: []byte("{\"hello\":\"world\"}"),
			wantErr:                 false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Method:      tt.fields.Method,
				Path:        tt.fields.Path,
				Host:        tt.fields.Host,
				Headers:     tt.fields.Headers,
				QueryParams: tt.fields.QueryParams,
				BodyParams:  tt.fields.BodyParams,
			}
			gotHttpRequest, err := req.build(tt.args.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHttpRequest.Method, tt.wantHttpRequest.Method) ||
				!reflect.DeepEqual(gotHttpRequest.URL, tt.wantHttpRequest.URL) ||
				!reflect.DeepEqual(gotHttpRequest.Proto, tt.wantHttpRequest.Proto) ||
				!reflect.DeepEqual(gotHttpRequest.ProtoMajor, tt.wantHttpRequest.ProtoMajor) ||
				!reflect.DeepEqual(gotHttpRequest.ProtoMinor, tt.wantHttpRequest.ProtoMinor) ||
				!reflect.DeepEqual(gotHttpRequest.Header, tt.wantHttpRequest.Header) ||
				!reflect.DeepEqual(gotHttpRequest.ContentLength, tt.wantHttpRequest.ContentLength) ||
				!reflect.DeepEqual(gotHttpRequest.Host, tt.wantHttpRequest.Host) {
				t.Errorf("build() gotHttpRequest = %v, want %v", gotHttpRequest, tt.wantHttpRequest)
				return
			}
			gotHttpRequestBodyData, _ := io.ReadAll(gotHttpRequest.Body)
			if !reflect.DeepEqual(gotHttpRequestBodyData, tt.wantHttpRequestBodyData) {
				t.Errorf("build() gotHttpRequestBodyData = %v, want %v", gotHttpRequestBodyData, tt.wantHttpRequestBodyData)
				return
			}
		})
	}
}
