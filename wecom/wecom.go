package wecom

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const (
	DefaultBaseURL = "https://qyapi.weixin.qq.com/cgi-bin/"
)

// 企业微信所有接口的 Response 均有这两个字段，用于判断请求结果
type IBase interface {
	GetErrCode() int
	GetErrMsg() string
}

// IBase 接口的基础实现
type Base struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (b Base) GetErrCode() int {
	return b.ErrCode
}

func (b Base) GetErrMsg() string {
	return b.ErrMsg
}

// 客户端
type Client struct {
	// 关于 access token 的生成可参考：https://work.weixin.qq.com/api/doc/90000/90135/91039
	// 需要在初始化客户端时传入
	enterpriseID string
	agentID      string
	agentSecret  string

	// base url，默认为：https://qyapi.weixin.qq.com/cgi-bin/
	baseURL string

	// token 通过调用 API 获取
	token *string

	// lock
	mu sync.Mutex

	// 通用 service
	comm service

	// HTTP Client
	client *http.Client

	// 对象
	Basic   *basicService
	Address *addressService
}

func NewClient(enterpriseID, agentID, agentSecret string) *Client {
	c := &Client{
		enterpriseID: enterpriseID,
		agentID:      agentID,
		agentSecret:  agentSecret,
		baseURL:      DefaultBaseURL,
		client:       &http.Client{},
	}
	c.comm.client = c

	c.Basic = (*basicService)(&c.comm)
	return c
}

func NewClientWithBaseURL(enterpriseID, agentID, agentSecret, baseURL string) *Client {
	c := &Client{
		enterpriseID: enterpriseID,
		agentID:      agentID,
		agentSecret:  agentSecret,
		baseURL:      baseURL,
		client:       &http.Client{},
	}
	c.comm.client = c

	c.Basic = (*basicService)(&c.comm)
	return c
}

func (c *Client) newRequest(httpMethod, url string, body interface{}) (*http.Request, error) {
	if strings.HasSuffix(c.baseURL, "/") && strings.HasPrefix(url, "/") {
		url = strings.TrimPrefix(url, "/")
	}
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(httpMethod, url, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) doRequest(req *http.Request, result IBase) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		return err
	}
	return nil
}

// 获取 token，如果 token 无效，则调用 API 获取 token
func (c *Client) getAccessToken() string {
	c.mu.Lock()
	if c.token == nil || *c.token == "" {
		c.Basic.refreshAccessToken()
	}
	c.mu.Unlock()

	return *c.token
}

// 判断错误码是否为 token 已过期
// errcode: 42001 token 已过期
// 企业微信错误码查询页面：https://open.work.weixin.qq.com/devtool/query?e=42001
// 企业微信全局错误码：https://open.work.weixin.qq.com/api/doc/90000/90139/90313
func (c *Client) tokenExpired(result IBase) bool {
	if result.GetErrCode() == 42001 {
		return true
	}
	return false
}

func (c *Client) resetToken() {
	c.mu.Lock()
	c.token = nil
	c.mu.Unlock()
}
