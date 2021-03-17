// basic.go 对应的是 https://work.weixin.qq.com/api/doc/90000/90135/90193 文档内容
package wecom

import (
	"context"
	"net/http"
)

const (
	pathUserCreate = "cgi-bin/user/create"
	pathUserGet    = "cgi-bin/user/get"
	pathUserUpdate = "cgi-bin/user/update"
	pathUserDelete = "cgi-bin/user/delete"
)

type addressService service

func (b *addressService) WithContext(ctx context.Context) *addressService {
	return &addressService{
		client: b.client,
		ctx:    ctx,
	}
}

// https://work.weixin.qq.com/api/doc/90000/90135/90195
type User struct {
	baseResponse
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
	baseResponse
}

// 通讯录：创建成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90195
func (b *addressService) CreateMember(user *User) (result *UserResp, err error) {
	req, err := b.client.newRequest(http.MethodPost, pathUserCreate, user)
	if err != nil {
		return nil, err
	}
	result = new(UserResp)
	err = (*service)(b).doRequest(req, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 通讯录：读取（单个）成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90196
func (b *addressService) GetMember(userID string) (result *User, err error) {
	req, err := b.client.newRequest(http.MethodPost, pathUserGet, nil, "userid="+userID)
	if err != nil {
		return nil, err
	}
	result = new(User)
	err = (*service)(b).doRequest(req, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 通讯录：更新成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90197
func (b *addressService) UpdateMember(user *User) (result *UserResp, err error) {
	req, err := b.client.newRequest(http.MethodPost, pathUserUpdate, user)
	if err != nil {
		return nil, err
	}
	result = new(UserResp)
	err = (*service)(b).doRequest(req, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 通讯录：删除成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90197
func (b *addressService) DeleteMember(userID string) (result *UserResp, err error) {
	req, err := b.client.newRequest(http.MethodPost, pathUserDelete, nil, "userid="+userID)
	if err != nil {
		return nil, err
	}
	result = new(UserResp)
	err = (*service)(b).doRequest(req, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
