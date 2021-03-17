# wecomgo

Wecom sdk for go，一个简单实用的企业微信 SDK

# Feature

- 极简的设计
- 完全按照企业微信官方文档结构对 `API` 进行封装，根据文档即可知道 `API` 用法
- 可选的 `Context`
- 支持自定义 `API Host`、`HTTP Client`
- 无需关心 `Access Token`，由 `Wecom` 自行维护，您只需要关注业务逻辑


# Quick Started

`Wecomgo` 只需要非常简单的几行代码，即可完成 API 的调用

```go
package main

import (
	"fmt"
	
	"github.com/3ks/wecomgo/wecom"
)

func main() {
	client, err := wecom.NewClient("企业 ID", "应用 Secret")
	if err != nil {
		panic(err)
	}
	user, err := client.Address.GetMember("3ks")
	if err != nil {
		fmt.Printf("get member fail, err:%v\n", err)
	} else {
		fmt.Printf("get member success, user:%v",user)
	}
}
```

# Context

`Wecomgo` 的所有 `API` 调用均支持 `Context`，并且这是可选的。

你只需要在调用具体方法之前调用 `WithContext` 即可

```go
package main

import (
	"context"
	"fmt"

	"github.com/3ks/wecomgo/wecom"
)

func main() {
	client, err := wecom.NewClient("企业 ID", "应用 Secret")
	if err != nil {
		panic(err)
	}
	ctx:=context.Background()
	user, err := client.Address.WithContext(ctx).GetMember("3ks")
	if err != nil {
		fmt.Printf("get member fail, err:%v\n", err)
	} else {
		ffmt.Printf("get member success, user:%v",user)
	}
}
```

# 自定义 Client

`Wecomgo` 支持自定义 `HTTP Client` 和 API 服务器的 `Host`，它们同样很简单

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/3ks/wecomgo/wecom"
)

func main() {
	client, err := wecom.NewClient("企业 ID", "应用 Secret",
		wecom.NewWithHostOption("http://baidu.com"),
		wecom.NewWithHTTPClientOption(&http.Client{}),
		)
	if err != nil {
		panic(err)
	}
	user, err := client.Address.GetMember("3ks")
	if err != nil {
		fmt.Printf("get member fail, err:%v\n", err)
	} else {
		fmt.Printf("get member success, user:%v",user)
	}
}
```