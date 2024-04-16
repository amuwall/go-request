package request

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Method string
	Path   string

	Host    string
	Headers http.Header

	QueryParams QueryParams
	BodyParams  BodyParams
}

type RequestOption func(*Request) error

func NewRequest(method, path string, options ...RequestOption) (request *Request, err error) {
	request = &Request{
		Method: method,
		Path:   path,

		Host:    "",
		Headers: http.Header{},
	}

	for _, option := range options {
		err = option(request)
		if err != nil {
			return
		}
	}

	return
}

func WithHost(host string) RequestOption {
	return func(r *Request) error {
		r.Host = host
		return nil
	}
}

func WithHeaders(headers map[string]string) RequestOption {
	return func(r *Request) error {
		for k, v := range headers {
			r.Headers.Set(k, v)
		}
		return nil
	}
}

func WithQueryParams(queryParams QueryParams) RequestOption {
	return func(r *Request) error {
		r.QueryParams = queryParams
		return nil
	}
}

func WithBodyParams(bodyParams BodyParams) RequestOption {
	return func(r *Request) error {
		r.BodyParams = bodyParams
		return nil
	}
}

func (req *Request) build(baseURL string) (httpRequest *http.Request, err error) {
	requestURL, err := url.JoinPath(baseURL, req.Path)
	if err != nil {
		err = fmt.Errorf("build url path error %w", err)
		return
	}

	var requestBody io.Reader
	if req.BodyParams != nil {
		var contentType string
		contentType, requestBody, err = req.BodyParams.Build()
		if err != nil {
			err = fmt.Errorf("build body params error %w", err)
			return
		}
		if contentType != "" {
			if req.Headers.Get("Content-Type") != "" {
				req.Headers.Set("Content-Type", contentType)
			}
		}
	}

	httpRequest, err = http.NewRequest(req.Method, requestURL, requestBody)
	if err != nil {
		err = fmt.Errorf("new http request error %w", err)
		return
	}

	if len(req.Host) != 0 {
		httpRequest.Host = req.Host
	}

	if req.Headers != nil {
		httpRequest.Header = req.Headers
	}

	if len(req.QueryParams) != 0 {
		httpRequest.URL.RawQuery = req.QueryParams.Encode()
	}

	return
}
