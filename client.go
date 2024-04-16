package request

import (
	"net/http"
	"time"
)

type Client struct {
	Scheme string
	Host   string
	Port   uint16

	instance *http.Client
}

type ClientOption func(*Client) error

const (
	defaultScheme = "https"
	defaultPort   = 443
)

func NewClient(host string, options ...ClientOption) (client *Client, err error) {
	client = &Client{
		Scheme: defaultScheme,
		Host:   host,
		Port:   defaultPort,

		instance: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, option := range options {
		err = option(client)
		if err != nil {
			return
		}
	}

	return
}
