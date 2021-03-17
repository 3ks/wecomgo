// basic.go 对应的是 https://work.weixin.qq.com/api/doc/90000/90135/90664 文档内容
// 主要实现了是获取 access_token 的 API
package wecom

import (
	"net/http"
	"time"
)

const (
	pathGetToken = "cgi-bin/gettoken" // 获取 token path
)

// 获取 Access Token：https://work.weixin.qq.com/api/doc/90000/90135/91039
// 频率限制：https://open.work.weixin.qq.com/api/doc/90000/90139/90312
// 加解密方案：https://open.work.weixin.qq.com/api/doc/90000/90139/90968
type basicService service

type Basic struct {
	baseResponse
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// 一般仅在 token 过期的情况下刷新 token
// 参数要求及含义参考：https://work.weixin.qq.com/api/doc/90000/90135/91039
// TODO 优化更新时机
// TODO 处理失败
func (b *basicService) refreshAccessToken() {
	// 调用 API 获取 token
	req, err := b.client.newRequest(http.MethodGet, pathGetToken, nil, "corpid="+b.client.enterpriseID, "corpsecret="+b.client.agentSecret)
	if err != nil {
		panic(err)
	}
	result := new(Basic)
	err = b.client.do(req, result)
	if err != nil {
		panic(err)
	}
	b.client.mu.Lock()
	defer b.client.mu.Unlock()
	b.client.token = result.AccessToken
	b.client.expireAt = time.Now().Unix() + result.ExpiresIn
}
