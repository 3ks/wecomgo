package wecom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultHost   = "https://qyapi.weixin.qq.com"
	maxRetryTimes = 3
)

// 企业微信所有接口的 Response 均有这两个字段，用于判断请求结果
type iBaseResponse interface {
	GetErrCode() int
	GetErrMsg() string
}

// iBaseResponse 接口的基础实现
type baseResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (b baseResponse) GetErrCode() int {
	return b.ErrCode
}

func (b baseResponse) GetErrMsg() string {
	return b.ErrMsg
}

// 客户端
type client struct {
	// 关于 access token 的生成可参考：https://work.weixin.qq.com/api/doc/90000/90135/91039
	enterpriseID string
	agentSecret  string

	// host，默认为：https://qyapi.weixin.qq.com
	host    string
	hostURL *url.URL

	// token 通过调用 API 获取
	token    string
	expireAt int64

	// lock，主要用于更新 token
	mu *sync.RWMutex

	// HTTP client
	client *http.Client

	comm service

	Basic   *basicService
	Address *addressService
}

func (c client) String() string {
	return fmt.Sprintf("enterprise: %s\napi host:%s\n", c.enterpriseID, c.host)
}

func NewClient(enterpriseID, agentSecret string, opts ...options) (c *client, err error) {
	c = &client{
		enterpriseID: enterpriseID,
		agentSecret:  agentSecret,
		host:         defaultHost,
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
func (c *client) newRequest(httpMethod, path string, body interface{}, queryString ...string) (request *http.Request, err error) {
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
		request.Header.Set("User-Agent", "WecomGo")
	}

	return request, nil
}

func (c *client) do(req *http.Request, result iBaseResponse) (err error) {
	for {
		if req.URL.Path != pathGetToken {
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
			c.Basic.refreshAccessToken()
			continue
		}
		// 请求无异常，break
		return nil
	}
}

// 获取 token，如果 token 无效，则调用 API 获取 token
func (c *client) getAccessToken() string {
	if c.token == "" {
		c.Basic.refreshAccessToken()
	}
	return c.token
}

// 判断错误码是否为 token 已过期
// errcode: 42001 token 已过期
// 企业微信错误码查询页面：https://open.work.weixin.qq.com/devtool/query?e=42001
// 企业微信全局错误码：https://open.work.weixin.qq.com/api/doc/90000/90139/90313
func (c *client) tokenExpired(result iBaseResponse) bool {
	if result.GetErrCode() == 42001 {
		return true
	}
	return false
}
