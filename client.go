package request

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Scheme string
	Host   string
	Port   uint16

	instance  *http.Client
	transport *http.Transport
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

		instance:  &http.Client{},
		transport: &http.Transport{},
	}

	for _, option := range options {
		err = option(client)
		if err != nil {
			return
		}
	}

	client.instance.Transport = client.transport

	return
}

func WithScheme(scheme string) ClientOption {
	return func(c *Client) error {
		c.Scheme = scheme
		return nil
	}
}

func WithPort(port uint16) ClientOption {
	return func(c *Client) error {
		c.Port = port
		return nil
	}
}

func WithTransport(transport *http.Transport) ClientOption {
	return func(c *Client) error {
		c.transport = transport
		return nil
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.instance.Timeout = timeout
		return nil
	}
}

func WithClientCertificateBlock(clientCrtBlock, clientKeyBlock []byte) ClientOption {
	return func(c *Client) error {
		if c.transport.TLSClientConfig == nil {
			c.transport.TLSClientConfig = &tls.Config{}
		}

		certificate, err := tls.X509KeyPair(clientCrtBlock, clientKeyBlock)
		if err != nil {
			return err
		}

		c.transport.TLSClientConfig.Certificates = append(
			c.transport.TLSClientConfig.Certificates, certificate,
		)

		return nil
	}
}

func WithClientCertificateFile(clientCrtFile, clientKeyFile string) ClientOption {
	return func(c *Client) error {
		if c.transport.TLSClientConfig == nil {
			c.transport.TLSClientConfig = &tls.Config{}
		}

		certificate, err := tls.LoadX509KeyPair(clientCrtFile, clientKeyFile)
		if err != nil {
			return err
		}

		c.transport.TLSClientConfig.Certificates = append(
			c.transport.TLSClientConfig.Certificates, certificate,
		)

		return nil
	}
}

func WithTLSServerName(serverName string) ClientOption {
	return func(c *Client) error {
		if c.transport.TLSClientConfig == nil {
			c.transport.TLSClientConfig = &tls.Config{}
		}

		c.transport.TLSClientConfig.ServerName = serverName

		return nil
	}
}

func WithSkipVerifyCertificates() ClientOption {
	return func(c *Client) error {
		if c.transport.TLSClientConfig == nil {
			c.transport.TLSClientConfig = &tls.Config{}
		}

		c.transport.TLSClientConfig.InsecureSkipVerify = true

		return nil
	}
}

func (c *Client) BaseURL() string {
	return fmt.Sprintf("%s://%s:%d", c.Scheme, c.Host, c.Port)
}

func (c *Client) Do(req *Request) (*Response, error) {
	httpRequest, err := req.build(c.BaseURL())
	if err != nil {
		return nil, err
	}

	httpResponse, err := c.instance.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	return parseResponse(httpResponse)
}
