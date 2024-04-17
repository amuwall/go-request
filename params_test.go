package request

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewQueryParams(t *testing.T) {
	type args struct {
		params map[string]string
	}
	tests := []struct {
		name string
		args args
		want QueryParams
	}{
		{
			name: "new query params",
			args: args{
				params: map[string]string{
					"test": "value",
				},
			},
			want: QueryParams{
				"test": []string{"value"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryParams(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryParams_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		p    QueryParams
		args args
		want QueryParams
	}{
		{
			name: "set query params",
			p: QueryParams{
				"test": []string{"old"},
			},
			args: args{
				key:   "test",
				value: "new",
			},
			want: QueryParams{
				"test": []string{"new"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Set(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.p, tt.want) {
				t.Errorf("Set() = %v, want %v", tt.p, tt.want)
			}
		})
	}
}

func TestQueryParams_Add(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		p    QueryParams
		args args
		want QueryParams
	}{
		{
			name: "add query params",
			p: QueryParams{
				"test": []string{"old"},
			},
			args: args{
				key:   "test",
				value: "new",
			},
			want: QueryParams{
				"test": []string{"old", "new"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Add(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.p, tt.want) {
				t.Errorf("Add() = %v, want %v", tt.p, tt.want)
			}
		})
	}
}

func TestQueryParams_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		p    QueryParams
		args args
		want string
	}{
		{
			name: "get query params",
			p: QueryParams{
				"test": []string{"value"},
			},
			args: args{
				key: "test",
			},
			want: "value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Get(tt.args.key); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryParams_Encode(t *testing.T) {
	tests := []struct {
		name string
		p    QueryParams
		want string
	}{
		{
			name: "encode query params",
			p: QueryParams{
				"test":   []string{"value"},
				"test-2": []string{"value-2"},
			},
			want: "test=value&test-2=value-2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Encode(); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJsonBodyParams(t *testing.T) {
	type args struct {
		params interface{}
	}
	tests := []struct {
		name string
		args args
		want *JsonBodyParams
	}{
		{
			name: "new json body params",
			args: args{
				params: map[string]interface{}{
					"test-string": "value",
					"test-int":    1,
				},
			},
			want: &JsonBodyParams{
				params: map[string]interface{}{
					"test-string": "value",
					"test-int":    1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJsonBodyParams(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJsonBodyParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonBodyParams_Build(t *testing.T) {
	type fields struct {
		params interface{}
	}
	tests := []struct {
		name            string
		fields          fields
		wantContentType string
		wantBodyData    []byte
		wantErr         bool
	}{
		{
			name: "build json params",
			fields: fields{
				params: map[string]interface{}{
					"test-string": "value",
					"test-int":    1,
				},
			},
			wantContentType: "application/json; charset=UTF-8",
			wantBodyData:    []byte(`{"test-int":1,"test-string":"value"}`),
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JsonBodyParams{
				params: tt.fields.params,
			}
			gotContentType, gotBody, err := p.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotContentType != tt.wantContentType {
				t.Errorf("Build() gotContentType = %v, want %v", gotContentType, tt.wantContentType)
				return
			}
			gotBodyData, err := io.ReadAll(gotBody)
			if err != nil {
				t.Errorf("Build() read from body error %v", err)
				return
			}
			if !reflect.DeepEqual(gotBodyData, tt.wantBodyData) {
				t.Errorf("Build() gotBody data = %v, want %v", string(gotBodyData), string(tt.wantBodyData))
			}
		})
	}
}

func TestNewFormBodyParams(t *testing.T) {
	type args struct {
		params map[string]string
	}
	tests := []struct {
		name string
		args args
		want *FormBodyParams
	}{
		{
			name: "new form body params",
			args: args{
				params: map[string]string{
					"test": "value",
				},
			},
			want: &FormBodyParams{
				params: map[string]string{
					"test": "value",
				},
				files: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFormBodyParams(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFormBodyParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormBodyParams_AddFile(t *testing.T) {
	type fields struct {
		params map[string]string
		files  []*formFile
	}
	type args struct {
		fieldName string
		fileName  string
		reader    io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *FormBodyParams
	}{
		{
			name: "add file to form body params",
			fields: fields{
				params: map[string]string{
					"test": "value",
				},
				files: nil,
			},
			args: args{
				fieldName: "file",
				fileName:  "hello.txt",
				reader:    bytes.NewReader([]byte("hello world")),
			},
			want: &FormBodyParams{
				params: map[string]string{
					"test": "value",
				},
				files: []*formFile{
					{
						FieldName: "file",
						FileName:  "hello.txt",
						Reader:    bytes.NewReader([]byte("hello world")),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &FormBodyParams{
				params: tt.fields.params,
				files:  tt.fields.files,
			}
			p.AddFile(tt.args.fieldName, tt.args.fileName, tt.args.reader)
			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("AddFile() = %v, want %v", p, tt.want)
			}
		})
	}
}

func TestFormBodyParams_Build(t *testing.T) {
	type fields struct {
		params map[string]string
		files  []*formFile
	}
	tests := []struct {
		name            string
		fields          fields
		wantContentType string
		wantBodyData    []byte
		wantErr         bool
	}{
		{
			name: "build form params",
			fields: fields{
				params: map[string]string{
					"test": "value",
				},
				files: nil,
			},
			wantContentType: "multipart/form-data; boundary=",
			wantBodyData:    []byte("Content-Disposition: form-data; name=\"test\"\r\n\r\nvalue\r\n"),
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &FormBodyParams{
				params: tt.fields.params,
				files:  tt.fields.files,
			}
			gotContentType, gotBody, err := p.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.HasPrefix(gotContentType, tt.wantContentType) {
				t.Errorf("Build() gotContentType = %v, want %v", gotContentType, tt.wantContentType)
				return
			}
			gotBodyData, err := io.ReadAll(gotBody)
			if err != nil {
				t.Errorf("Build() read from body error %v", err)
				return
			}
			if !strings.Contains(string(gotBodyData), string(tt.wantBodyData)) {
				t.Errorf("Build() gotBody data = %v, want %v", string(gotBodyData), string(tt.wantBodyData))
			}
		})
	}
}
