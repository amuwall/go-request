package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int
	Header     http.Header
	RawBody    []byte
}

func parseResponse(httpResponse *http.Response) (response *Response, err error) {
	defer httpResponse.Body.Close()

	rawBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	response = &Response{
		StatusCode: httpResponse.StatusCode,
		Header:     httpResponse.Header,
		RawBody:    rawBody,
	}

	return
}

func (resp *Response) UnmarshalJSONBody(val interface{}) (err error) {
	if !strings.Contains(resp.Header.Get(contentTypeHeader), contentTypeJson) {
		return fmt.Errorf("response content-type not json, it is %s", resp.Header.Get("Content-Type"))
	}
	return json.Unmarshal(resp.RawBody, val)
}
