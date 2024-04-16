package request

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		host    string
		options []ClientOption
	}
	tests := []struct {
		name       string
		args       args
		wantClient *Client
		wantErr    bool
	}{
		{
			name: "new client",
			args: args{
				host:    "127.0.0.1",
				options: nil,
			},
			wantClient: &Client{
				Scheme: "https",
				Host:   "127.0.0.1",
				Port:   443,
				instance: &http.Client{
					Transport: &http.Transport{},
				},
				transport: &http.Transport{},
			},
			wantErr: false,
		},
		{
			name: "new client with option to set host",
			args: args{
				host: "127.0.0.1",
				options: []ClientOption{
					func(client *Client) (err error) {
						client.Host = "127.0.0.2"
						return err
					},
				},
			},
			wantClient: &Client{
				Scheme: "https",
				Host:   "127.0.0.2",
				Port:   443,
				instance: &http.Client{
					Transport: &http.Transport{},
				},
				transport: &http.Transport{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, err := NewClient(tt.args.host, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotClient, tt.wantClient) {
				t.Errorf("NewClient() gotClient = %v, want %v", gotClient, tt.wantClient)
			}
		})
	}
}
