package kubeclient

import (
	"net/http"
	"net/url"
	"time"

	"github.com/leonardodelira/go-lib-http/deployment"
)

const urlDefault = "http://localhost:3333"
const timeoutDefault = 2500

type Client struct {
	url        string
	timeout    time.Duration
	httpClient *http.Client

	Deployment deployment.Service
}

type option func(*Client) error

func NewClient(options ...option) (*Client, error) {
	c := &Client{
		url:        urlDefault,
		timeout:    timeoutDefault,
		httpClient: &http.Client{},
	}
	for _, o := range options {
		err := o(c)
		if err != nil {
			return nil, err
		}
	}

	c.Deployment = deployment.NewService(c.httpClient, c.url)

	return c, nil
}

func WithURL(u string) option {
	return func(c *Client) error {
		if _, err := url.ParseRequestURI(u); err != nil {
			return err
		}
		c.url = u
		return nil
	}
}

func WithTimeOut(timeout time.Duration) option {
	return func(c *Client) error {
		c.timeout = timeout
		return nil
	}
}

func WithHttpClient(httpClient *http.Client) option {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}
