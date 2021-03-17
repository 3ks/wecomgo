package wecom

import "net/http"

// 目前支持两个 options，hostURL 和 HTTP client
type options interface {
	applyOption(*client)
}

type optHost struct {
	Host string
}

func (o *optHost) applyOption(client *client) {
	client.host = o.Host
}

func NewWithHostOption(host string) options {
	return &optHost{
		Host: host,
	}
}

type optHTTPClient struct {
	Client *http.Client
}

func (o *optHTTPClient) applyOption(client *client) {
	client.client = o.Client
}

func NewWithHTTPClientOption(client *http.Client) options {
	return &optHTTPClient{
		Client: client,
	}
}
