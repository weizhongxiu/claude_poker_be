package v1

import "github.com/gogf/gf/v2/frame/g"

// --- 注册 ---

type RegisterReq struct {
	g.Meta   `path:"/user/register" method:"post" tags:"用户" summary:"注册"`
	Phone    string `json:"phone" v:"required|phone#手机号不能为空|手机号格式错误"`
	Password string `json:"password" v:"required|min-length:6#密码不能为空|密码不能少于6位"`
	Nickname string `json:"nickname" v:"required#昵称不能为空"`
}

type RegisterRes struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

// --- 登录 ---

type LoginReq struct {
	g.Meta   `path:"/user/login" method:"post" tags:"用户" summary:"登录"`
	Phone    string `json:"phone" v:"required#手机号不能为空"`
	Password string `json:"password" v:"required#密码不能为空"`
}

type LoginRes struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Token    string `json:"token"`
}

// --- 获取个人信息 ---

type ProfileReq struct {
	g.Meta `path:"/user/profile" method:"get" tags:"用户" summary:"获取个人信息"`
}

type ProfileRes struct {
	UserID   int64  `json:"user_id"`
	UID      string `json:"uid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Chips    int64  `json:"chips"`
}
