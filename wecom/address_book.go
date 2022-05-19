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
	pathUserSimpleList = "/cgi-bin/user/simplelist"
	pathUserList       = "/cgi-bin/user/list"
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
	Userid           string           `json:"userid"`
	Name             string           `json:"name,omitempty"`
	Alias            string           `json:"alias,omitempty"`
	Mobile           string           `json:"mobile,omitempty"`
	Department       []int            `json:"department,omitempty"`
	Order            []int            `json:"order,omitempty"`
	Position         string           `json:"position,omitempty"`
	Gender           string           `json:"gender,omitempty"`
	Email            string           `json:"email,omitempty"`
	IsLeaderInDept   []int            `json:"is_leader_in_dept,omitempty"`
	Enable           *int             `json:"enable,omitempty"`
	AvatarMediaid    string           `json:"avatar_mediaid,omitempty"`
	Telephone        string           `json:"telephone,omitempty"`
	Address          string           `json:"address,omitempty"`
	MainDepartment   int              `json:"main_department,omitempty"`
	Extattr          *Extattr         `json:"extattr,omitempty"`
	ToInvite         *bool            `json:"to_invite,omitempty"`
	ExternalPosition string           `json:"external_position,omitempty"`
	ExternalProfile  *ExternalProfile `json:"external_profile,omitempty"`
}

type Text struct {
	Value string `json:"value"`
}

type Web struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type Attrs struct {
	Type int    `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	Text Text   `json:"text,omitempty"`
	Web  Web    `json:"web,omitempty"`
}

type Extattr struct {
	Attrs []Attrs `json:"attrs,omitempty"`
}

type Miniprogram struct {
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath"`
	Title    string `json:"title"`
}

type ExternalAttr struct {
	Type        int         `json:"type,omitempty"`
	Name        string      `json:"name,omitempty"`
	Text        Text        `json:"text,omitempty"`
	Web         Web         `json:"web,omitempty"`
	Miniprogram Miniprogram `json:"miniprogram,omitempty"`
}

type ExternalProfile struct {
	ExternalCorpName string         `json:"external_corp_name,omitempty"`
	ExternalAttr     []ExternalAttr `json:"external_attr,omitempty"`
}

type UserResp struct {
	baseResponse
}

// 通讯录：创建成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90195
func (b *addressService) CreateMember(user *User) (result *UserResp, err error) {
	failCount := -1
	for failCount < b.client.maxRetryTimes {
		// 默认尝试一次，即不进行失败重试
		failCount++

		// 每次请求重新生成 request
		var req *http.Request
		req, err = b.client.newRequest(http.MethodPost, pathUserCreate, user)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		result = new(UserResp)
		err = (*service)(b).doRequest(req, result)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		// 成功，return
		return result, nil
	}
	// 失败，返回最后一次请求的 err
	return nil, err
}

// 通讯录：读取（单个）成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90196
func (b *addressService) GetMember(userID string) (result *User, err error) {
	failCount := -1
	// 默认尝试一次，即不进行失败重试
	for failCount < b.client.maxRetryTimes {
		failCount++
		var req *http.Request
		req, err = b.client.newRequest(http.MethodPost, pathUserGet, nil, "userid="+userID)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		result = new(User)
		err = (*service)(b).doRequest(req, result)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		// 成功，return
		return result, nil
	}
	// 失败，返回最后一次请求的 err
	return nil, err
}

type SimpleUserList struct {
	baseResponse
	Userlist []SimpleUser `json:"userlist"`
}

type SimpleUser struct {
	Userid     string `json:"userid"`
	Name       string `json:"name"`
	Department []int  `json:"department"`
}

// ListMember 通讯录：列出成员简单信息
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90196
// departmentID 部门ID，根部门填 1
// recursive 是否递归获取子部门成员，0 表示不需要递归获取，否则表示需要递归获取
func (b *addressService) ListMember(departmentID, recursive int) (result *SimpleUserList, err error) {
	failCount := -1
	if departmentID < 1 {
		return nil, fmt.Errorf("invalid department id: %d", departmentID)
	}
	// 只要不为0，则置为1
	if recursive != 0 {
		recursive = 1
	}

	// 默认尝试一次，即不进行失败重试
	for failCount < b.client.maxRetryTimes {
		failCount++
		var req *http.Request
		req, err = b.client.newRequest(http.MethodPost, pathUserSimpleList, nil, fmt.Sprintf("department_id=%d", departmentID), fmt.Sprintf("fetch_child=%d", recursive))
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		result = new(SimpleUserList)
		err = (*service)(b).doRequest(req, result)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		// 成功，return
		return result, nil
	}
	// 失败，返回最后一次请求的 err
	return nil, err
}

type UserList struct {
	baseResponse
	Userlist []SimpleUserWithMainDepartment `json:"userlist"`
}

type SimpleUserWithMainDepartment struct {
	Userid         string `json:"userid"`
	Name           string `json:"name"`
	Department     []int  `json:"department"`
	MainDepartment int    `json:"main_department"`
	// TODO 其他字段
}

// ListMembers 通讯录：列出成员详情
// 参考链接：https://developer.work.weixin.qq.com/document/path/90201
// departmentID 部门ID，根部门填 1
// recursive 是否递归获取子部门成员，0 表示不需要递归获取，否则表示需要递归获取
func (b *addressService) ListMembers(departmentID, recursive int) (result *UserList, err error) {
	failCount := -1
	if departmentID < 1 {
		return nil, fmt.Errorf("invalid department id: %d", departmentID)
	}
	// 只要不为0，则置为1
	if recursive != 0 {
		recursive = 1
	}

	// 默认尝试一次，即不进行失败重试
	for failCount < b.client.maxRetryTimes {
		failCount++
		var req *http.Request
		req, err = b.client.newRequest(http.MethodPost, pathUserList, nil, fmt.Sprintf("department_id=%d", departmentID), fmt.Sprintf("fetch_child=%d", recursive))
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		result = new(UserList)
		err = (*service)(b).doRequest(req, result)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		// 成功，return
		return result, nil
	}
	// 失败，返回最后一次请求的 err
	return nil, err
}

// 通讯录：更新成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90197
func (b *addressService) UpdateMember(user *User) (result *UserResp, err error) {
	failCount := -1
	// 默认尝试一次，即不进行失败重试
	for failCount < b.client.maxRetryTimes {
		failCount++
		var req *http.Request
		req, err = b.client.newRequest(http.MethodPost, pathUserUpdate, user)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		result = new(UserResp)
		err = (*service)(b).doRequest(req, result)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		// 成功，return
		return result, nil
	}
	// 失败，返回最后一次请求的 err
	return nil, err
}

// 通讯录：删除成员
// 参考链接：https://work.weixin.qq.com/api/doc/90000/90135/90197
func (b *addressService) DeleteMember(userID string) (result *UserResp, err error) {
	failCount := -1
	// 默认尝试一次，即不进行失败重试
	for failCount < b.client.maxRetryTimes {
		failCount++
		var req *http.Request
		req, err = b.client.newRequest(http.MethodPost, pathUserDelete, nil, "userid="+userID)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		result = new(UserResp)
		err = (*service)(b).doRequest(req, result)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		// 成功，return
		return result, nil
	}
	// 失败，返回最后一次请求的 err
	return nil, err
}

type invite struct {
	User []string `json:"user"`
}

// 通讯录：批量邀请成员
// 参考链接：https://open.work.weixin.qq.com/api/doc/90000/90135/90975
func (b *addressService) InviteMember(userID []string) (result *UserResp, err error) {
	failCount := -1
	// 默认尝试一次，即不进行失败重试
	for failCount < b.client.maxRetryTimes {
		failCount++
		body := invite{User: userID}
		var req *http.Request
		req, err = b.client.newRequest(http.MethodPost, pathUserInvite, body)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		result = new(UserResp)
		err = (*service)(b).doRequest(req, result)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		// 成功，return
		return result, nil
	}
	// 失败，返回最后一次请求的 err
	return nil, err
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
	failCount := -1
	// 默认尝试一次，即不进行失败重试
	for failCount < b.client.maxRetryTimes {
		failCount++
		var req *http.Request
		req, err = b.client.newRequest(http.MethodGet, pathDepartmentList, nil, fmt.Sprintf("id=%d", departmentID))
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		result = new(DepartmentList)
		err = (*service)(b).doRequest(req, result)
		if err != nil {
			continue // 如果循环结束，则会返回该 err
		}
		// 成功，return
		return result, nil
	}
	// 失败，返回最后一次请求的 err
	return nil, err
}
