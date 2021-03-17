// basic.go 对应的是 https://work.weixin.qq.com/api/doc/90000/90135/90193 文档内容
package wecom

import (
	"fmt"
	"net/http"
)

type addressService service

// https://work.weixin.qq.com/api/doc/90000/90135/90195
type User struct {
	Base
	Userid         string `json:"userid"` // 必填参数
	Name           string `json:"name"`   // 必填参数
	Alias          string `json:"alias,omitempty"`
	Mobile         string `json:"mobile,omitempty"`
	Department     []int  `json:"department,omitempty"`
	Order          []int  `json:"order,omitempty"`
	Position       string `json:"position,omitempty"`
	Gender         string `json:"gender,omitempty"`
	Email          string `json:"email,omitempty"`
	IsLeaderInDept []int  `json:"is_leader_in_dept,omitempty"`
	Enable         int    `json:"enable,omitempty"`
	AvatarMediaid  string `json:"avatar_mediaid,omitempty"`
	Telephone      string `json:"telephone,omitempty"`
	Address        string `json:"address,omitempty"`
	MainDepartment int    `json:"main_department,omitempty"`
	Extattr        struct {
		Attrs []struct {
			Type int    `json:"type,omitempty"`
			Name string `json:"name,omitempty"`
			Text struct {
				Value string `json:"value,omitempty"`
			} `json:"text,omitempty,omitempty"`
			Web struct {
				URL   string `json:"url,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"web,omitempty,omitempty"`
		} `json:"attrs,omitempty"`
	} `json:"extattr,omitempty"`
	ToInvite         bool   `json:"to_invite,omitempty"`
	ExternalPosition string `json:"external_position,omitempty"`
	ExternalProfile  struct {
		ExternalCorpName string `json:"external_corp_name,omitempty"`
		ExternalAttr     []struct {
			Type int    `json:"type,omitempty"`
			Name string `json:"name,omitempty"`
			Text struct {
				Value string `json:"value,omitempty"`
			} `json:"text,omitempty,omitempty"`
			Web struct {
				URL   string `json:"url,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"web,omitempty,omitempty"`
			Miniprogram struct {
				Appid    string `json:"appid,omitempty"`
				Pagepath string `json:"pagepath,omitempty"`
				Title    string `json:"title,omitempty"`
			} `json:"miniprogram,omitempty,omitempty"`
		} `json:"external_attr,omitempty"`
	} `json:"external_profile,omitempty"`
}

type UserResp struct {
	Base
}

// 通讯录：创建成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90195
func (b *addressService) CreateMember(user *User) (result *UserResp, err error) {
	// 调用 API
	fnReq := func() (*UserResp, error) {
		req, err := b.client.newRequest(http.MethodPost, "cgi-bin/user/create", user)
		if err != nil {
			return nil, err
		}
		result := new(UserResp)
		err = b.client.doRequest(req, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	for i := 0; i < 3; i++ {
		result, err = fnReq()
		// 请求出错
		if err != nil {
			// 直接 break
			break
		}
		// token 已过期
		if b.client.tokenExpired(result) {
			// 对于 token 过期的情况，重新获取 token
			b.client.resetToken()
			// 再次尝试
			i--
			continue
		}
		// 无错误，break
		break
	}
	if err != nil {
		return nil, fmt.Errorf("can not create member, err: %s\n", err.Error())
	}
	if result != nil {
		return result, nil
	} else {
		return result, fmt.Errorf("get result with nil")
	}
}

// 通讯录：读取（单个）成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90196
func (b *addressService) GetMember(userID string) (result *User, err error) {
	// 调用 API
	fnReq := func() (*User, error) {
		req, err := b.client.newRequest(http.MethodGet, "cgi-bin/user/get", nil, "userid="+userID)
		if err != nil {
			return nil, err
		}
		result := new(User)
		err = b.client.doRequest(req, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	for i := 0; i < 3; i++ {
		result, err = fnReq()
		// 请求出错
		if err != nil {
			// 直接 break
			break
		}
		// token 已过期
		if b.client.tokenExpired(result) {
			// 对于 token 过期的情况，重新获取 token
			b.client.resetToken()
			// 再次尝试
			i--
			continue
		}
		// 无错误，break
		break
	}
	// 出错，返回 err
	if err != nil {
		return nil, fmt.Errorf("can not get member, err: %s\n", err.Error())
	}
	// result 不为空，return
	if result != nil {
		return result, nil
	} else {
		// 极端情况，不太可能出现
		return result, fmt.Errorf("get result with nil")
	}
}

// 通讯录：更新成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90197
func (b *addressService) UpdateMember(user *User) (result *UserResp, err error) {
	// 调用 API
	fnReq := func() (*UserResp, error) {
		req, err := b.client.newRequest(http.MethodPost, "cgi-bin/user/update", user)
		if err != nil {
			return nil, err
		}
		result := new(UserResp)
		err = b.client.doRequest(req, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	for i := 0; i < 3; i++ {
		result, err = fnReq()
		// 请求出错
		if err != nil {
			// 直接 break
			break
		}
		// token 已过期
		if b.client.tokenExpired(result) {
			// 对于 token 过期的情况，重新获取 token
			b.client.resetToken()
			// 再次尝试
			i--
			continue
		}
		// 无错误，break
		break
	}
	// 出错，返回 err
	if err != nil {
		return nil, fmt.Errorf("can not update member, err: %s\n", err.Error())
	}
	// result 不为空，return
	if result != nil {
		return result, nil
	} else {
		// 极端情况，不太可能出现
		return result, fmt.Errorf("get result with nil")
	}
}

// 通讯录：删除成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90197
func (b *addressService) DeleteMember(userID string) (result *UserResp, err error) {
	// 调用 API
	fnReq := func() (*UserResp, error) {
		req, err := b.client.newRequest(http.MethodGet, "cgi-bin/user/delete", nil, "userid="+userID)
		if err != nil {
			return nil, err
		}
		result := new(UserResp)
		err = b.client.doRequest(req, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	for i := 0; i < 3; i++ {
		result, err = fnReq()
		// 请求出错
		if err != nil {
			// 直接 break
			break
		}
		// token 已过期
		if b.client.tokenExpired(result) {
			// 对于 token 过期的情况，重新获取 token
			b.client.resetToken()
			// 再次尝试
			i--
			continue
		}
		// 无错误，break
		break
	}
	// 出错，返回 err
	if err != nil {
		return nil, fmt.Errorf("can not delete member, err: %s\n", err.Error())
	}
	// result 不为空，return
	if result != nil {
		return result, nil
	} else {
		// 极端情况，不太可能出现
		return result, fmt.Errorf("get result with nil")
	}
}
