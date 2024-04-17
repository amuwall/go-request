package request

import (
	"net/http"
	"reflect"
	"testing"
	"time"
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

func TestWithScheme(t *testing.T) {
	type args struct {
		scheme string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "init client with scheme",
			args: args{
				scheme: "http",
			},
			want: "http",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := NewClient("127.0.0.1", WithScheme(tt.args.scheme))
			if client.Scheme != tt.want {
				t.Errorf("WithScheme() = %v, want %v", client.Scheme, tt.want)
			}
		})
	}
}

func TestWithPort(t *testing.T) {
	type args struct {
		port uint16
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{
			name: "init client with port",
			args: args{
				port: 80,
			},
			want: 80,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := NewClient("127.0.0.1", WithPort(tt.args.port))
			if client.Port != tt.want {
				t.Errorf("WithPort() = %v, want %v", client.Port, tt.want)
			}
		})
	}
}

func TestWithTransport(t *testing.T) {
	type args struct {
		transport *http.Transport
	}
	tests := []struct {
		name string
		args args
		want *http.Transport
	}{
		{
			name: "init client with transport",
			args: args{
				transport: &http.Transport{
					DisableKeepAlives:  true,
					DisableCompression: true,
					ForceAttemptHTTP2:  true,
				},
			},
			want: &http.Transport{
				DisableKeepAlives:  true,
				DisableCompression: true,
				ForceAttemptHTTP2:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := NewClient("127.0.0.1", WithTransport(tt.args.transport))
			if !reflect.DeepEqual(client.transport, tt.want) {
				t.Errorf("WithTransport() = %v, want %v", client.transport, tt.want)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "init client with timeout",
			args: args{
				timeout: time.Minute,
			},
			want: time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := NewClient("127.0.0.1", WithTimeout(tt.args.timeout))
			if client.instance.Timeout != tt.want {
				t.Errorf("WithTimeout() = %v, want %v", client.instance.Timeout, tt.want)
			}
		})
	}
}

func TestWithTLSServerName(t *testing.T) {
	type args struct {
		serverName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "init client with serverName",
			args: args{
				serverName: "example.com",
			},
			want: "example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := NewClient("127.0.0.1", WithTLSServerName(tt.args.serverName))
			if client.transport.TLSClientConfig.ServerName != tt.want {
				t.Errorf(
					"WithTLSServerName() = %v, want %v",
					client.transport.TLSClientConfig.ServerName, tt.want,
				)
			}
		})
	}
}

func TestWithSkipVerifyCertificates(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "init client with skipVerify",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := NewClient("127.0.0.1", WithSkipVerifyCertificates())
			if client.transport.TLSClientConfig.InsecureSkipVerify != tt.want {
				t.Errorf(
					"WithSkipVerifyCertificates() = %v, want %v",
					client.transport.TLSClientConfig.InsecureSkipVerify, tt.want,
				)
			}
		})
	}
}
