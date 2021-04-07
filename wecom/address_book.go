// basic.go 对应的是 https://work.weixin.qq.com/api/doc/90000/90135/90193 文档内容
package wecom

import (
	"context"
	"fmt"
	"net/http"
)

const (
	pathUserCreate     = "/cgi-bin/user/create"
	pathUserGet        = "/cgi-bin/user/get"
	pathUserUpdate     = "/cgi-bin/user/update"
	pathUserDelete     = "/cgi-bin/user/delete"
	pathUserInvite     = "/cgi-bin/batch/invite"
	pathDepartmentList = "/cgi-bin/department/list"
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
	Userid           string          `json:"userid"`
	Name             string          `json:"name"`
	Alias            string          `json:"alias"`
	Mobile           string          `json:"mobile"`
	Department       []int           `json:"department"`
	Order            []int           `json:"order"`
	Position         string          `json:"position"`
	Gender           string          `json:"gender"`
	Email            string          `json:"email"`
	IsLeaderInDept   []int           `json:"is_leader_in_dept"`
	Enable           *int            `json:"enable"`
	AvatarMediaid    string          `json:"avatar_mediaid"`
	Telephone        string          `json:"telephone"`
	Address          string          `json:"address"`
	MainDepartment   int             `json:"main_department"`
	Extattr          Extattr         `json:"extattr"`
	ToInvite         *bool           `json:"to_invite"`
	ExternalPosition string          `json:"external_position"`
	ExternalProfile  ExternalProfile `json:"external_profile"`
}

type Text struct {
	Value string `json:"value"`
}

type Web struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type Attrs struct {
	Type int    `json:"type"`
	Name string `json:"name"`
	Text Text   `json:"text,omitempty"`
	Web  Web    `json:"web,omitempty"`
}

type Extattr struct {
	Attrs []Attrs `json:"attrs"`
}

type Miniprogram struct {
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath"`
	Title    string `json:"title"`
}

type ExternalAttr struct {
	Type        int         `json:"type"`
	Name        string      `json:"name"`
	Text        Text        `json:"text,omitempty"`
	Web         Web         `json:"web,omitempty"`
	Miniprogram Miniprogram `json:"miniprogram,omitempty"`
}

type ExternalProfile struct {
	ExternalCorpName string         `json:"external_corp_name"`
	ExternalAttr     []ExternalAttr `json:"external_attr"`
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

type invite struct {
	User []string `json:"user"`
}

// 通讯录：批量邀请成员
// 参考链接：https://open.work.weixin.qq.com/api/doc/90000/90135/90975
func (b *addressService) InviteMember(userID []string) (result *UserResp, err error) {
	body := invite{User: userID}
	req, err := b.client.newRequest(http.MethodPost, pathUserInvite, body)
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

// 部门列表
type DepartmentList struct {
	baseResponse
	Department []Department `json:"department"`
}

type Department struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Parentid int    `json:"parentid"`
	Order    int    `json:"order"`
}

// 通讯录：获取部门列表
// 参考链接：https://open.work.weixin.qq.com/api/doc/90000/90135/90208
func (b *addressService) DepartmentList(departmentID int) (result *DepartmentList, err error) {
	req, err := b.client.newRequest(http.MethodGet, pathDepartmentList, nil, fmt.Sprintf("id=%d", departmentID))
	if err != nil {
		return nil, err
	}
	result = new(DepartmentList)
	err = (*service)(b).doRequest(req, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
