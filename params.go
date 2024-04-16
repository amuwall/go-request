package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

type QueryParams map[string][]string

func NewQueryParams(params map[string]string) QueryParams {
	p := QueryParams{}
	for key, value := range params {
		p.Set(key, value)
	}
	return p
}

func (p QueryParams) Set(key, value string) {
	p[key] = []string{value}
}

func (p QueryParams) Add(key, value string) {
	p[key] = append(p[key], value)
}

func (p QueryParams) Get(key string) string {
	value := p[key]
	if len(value) == 0 {
		return ""
	}
	return value[0]
}

func (p QueryParams) Encode() string {
	return url.Values(p).Encode()
}

type BodyParams interface {
	Build() (contentType string, body io.Reader, err error)
}

type JsonBodyParams struct {
	params interface{}
}

func NewJsonBodyParams(params interface{}) *JsonBodyParams {
	return &JsonBodyParams{
		params: params,
	}
}

func (p *JsonBodyParams) Build() (contentType string, body io.Reader, err error) {
	data, err := json.Marshal(p.params)
	if err != nil {
		err = fmt.Errorf("marshal error %w", err)
		return
	}

	contentType = contentTypeJsonWithUTF8
	body = bytes.NewReader(data)

	return
}

type FormBodyParams struct {
	params map[string]string
	files  []*formFile
}

type formFile struct {
	FieldName string
	FileName  string
	Reader    io.Reader
}

func NewFormBodyParams(params map[string]string) *FormBodyParams {
	return &FormBodyParams{
		params: params,
		files:  nil,
	}
}

func (p *FormBodyParams) AddFile(fieldName, fileName string, reader io.Reader) {
	p.files = append(p.files, &formFile{
		FieldName: fieldName,
		FileName:  fileName,
		Reader:    reader,
	})
}

func (p *FormBodyParams) Build() (contentType string, body io.Reader, err error) {
	buffer := &bytes.Buffer{}
	body = buffer

	bodyWriter := multipart.NewWriter(buffer)
	defer bodyWriter.Close()

	for key, value := range p.params {
		err = bodyWriter.WriteField(key, value)
		if err != nil {
			err = fmt.Errorf("write field %s error %w", key, err)
			return
		}
	}

	for _, file := range p.files {
		fileWriter, err := bodyWriter.CreateFormFile(file.FieldName, file.FileName)
		if err != nil {
			return "", nil, err
		}
		_, err = io.Copy(fileWriter, file.Reader)
		if err != nil {
			return "", nil, err
		}
	}

	contentType = bodyWriter.FormDataContentType()

	return
}
