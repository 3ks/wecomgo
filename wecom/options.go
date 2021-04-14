package wecom

import "net/http"

// 目前支持两个 options，hostURL 和 HTTP Client
type options interface {
	applyOption(*Client)
}

type optHost struct {
	Host string
}

func (o *optHost) applyOption(client *Client) {
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

func (o *optHTTPClient) applyOption(client *Client) {
	client.client = o.Client
}

func NewWithHTTPClientOption(client *http.Client) options {
	return &optHTTPClient{
		Client: client,
	}
}

type optPrintPayload struct {
	printPayload bool
}

func (o *optPrintPayload) applyOption(client *Client) {
	client.printPayload = o.printPayload
}

func NewWithPrintPayloadOption() options {
	return &optPrintPayload{
		printPayload: true,
	}
}
