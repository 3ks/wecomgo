package wecom

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	DefaultHost = "https://qyapi.weixin.qq.com"
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
	agentSecret  string

	// host，默认为：https://qyapi.weixin.qq.com
	host    string
	hostURL *url.URL

	// token 通过调用 API 获取
	token    *string
	expireAt int64

	// lock
	mu *sync.RWMutex

	// 通用 service
	comm service

	// HTTP Client
	client *http.Client

	// TODO 对象
	Basic   *basicService
	Address *addressService
}

func NewClient(enterpriseID, agentSecret string, opts ...options) (client *Client, err error) {
	c := &Client{
		enterpriseID: enterpriseID,
		agentSecret:  agentSecret,
		host:         DefaultHost,
		mu:           &sync.RWMutex{},
		client:       &http.Client{},
	}

	for k := range opts {
		opts[k].applyOption(c)
	}

	var u *url.URL
	u, err = url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	c.hostURL = u

	c.comm.client = c
	c.Basic = (*basicService)(&c.comm)
	c.Address = (*addressService)(&c.comm)

	return c, nil
}

// queryString 支持两种写法："name=3ks&age=18" 或者 "name=guan", "age=18"
func (c *Client) newRequest(httpMethod, path string, body interface{}, queryString ...string) (request *http.Request, err error) {
	// base info
	newURL := *c.hostURL
	newURL.Path = path

	// qs
	if len(queryString) > 0 {
		qs, err := url.ParseQuery(strings.Join(queryString, "&"))
		if err != nil {
			return nil, err
		}
		newURL.RawQuery = qs.Encode()
	}

	// body
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

	// new request
	request, err = http.NewRequest(httpMethod, newURL.String(), buf)
	if err != nil {
		return nil, err
	}
	// set header
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return request, nil
}

func (c *Client) setAK() {

}

func (c *Client) doRequest(req *http.Request, result IBase) (err error) {
	for {
		if req.URL.Path != PathGetToken {
			q := req.URL.Query()
			q.Set("access_token", c.getAccessToken())
			req.URL.RawQuery = q.Encode()
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			_ = resp.Body.Close()
			return err
		}
		_ = resp.Body.Close()

		err = json.Unmarshal(data, result)
		if err != nil {
			return err
		}
		// token 已过期
		if c.tokenExpired(result) {
			// 刷新 token
			// TODO 锁
			c.Basic.refreshAccessToken()
			// 重试
			continue
		}
		// 请求无异常，break
		return nil
	}
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
