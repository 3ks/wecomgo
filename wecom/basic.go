// basic.go 对应的是 https://work.weixin.qq.com/api/doc/90000/90135/90664 文档内容
// 主要实现了是获取 access_token 的 API
package wecom

import (
	"fmt"
	"net/http"
	"time"

	"github.com/3ks/wecomgo/tools"
)

type basicService service

// https://work.weixin.qq.com/api/doc/90000/90135/91039
type Basic struct {
	Base
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

const (
	PathGetToken = "cgi-bin/gettoken" // 获取 token path
)

// 一般仅在 token 过期的情况下刷新 token
// 参数要求及含义参考：https://work.weixin.qq.com/api/doc/90000/90135/91039
// TODO 在 expires_in 前 10 秒刷新 token？
func (b *basicService) refreshAccessToken() {
	// 调用 API 获取 token
	refresh := func() (*Basic, error) {
		req, err := b.client.newRequest(http.MethodGet, fmt.Sprintf("gettoken?corpid=%s&corpsecret=%s", b.client.enterpriseID, b.client.agentSecret), nil)
		if err != nil {
			return nil, err
		}
		result := new(Basic)
		// TODO Context
		err = b.client.doRequest(req, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	for i := 0; i < 3; i++ {
		t, err := refresh()
		if err != nil {
			fmt.Printf("time: %s,refresh token err: %s\n", time.Now().String(), err.Error())
			time.Sleep(time.Millisecond * 200)
			continue
		}
		// 不能再加锁
		b.client.token = tools.StringPoint(t.AccessToken)
		// 验证 token 是否已过期
		if b.client.tokenExpired(t) {
			b.client.resetToken()
		}
	}
}
