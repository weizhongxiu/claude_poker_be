package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "claude-test/api/user/v1"
	"claude-test/internal/logic/user"
	"claude-test/utility/jwt"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	userID, token, err := user.Register(ctx, req.Phone, req.Password, req.Nickname)
	if err != nil {
		return
	}
	res = &v1.RegisterRes{UserID: userID, Token: token}
	return
}

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	u, wallet, token, err := user.Login(ctx, req.Phone, req.Password)
	if err != nil {
		return
	}
	_ = wallet
	res = &v1.LoginRes{
		UserID:   u.Id,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Token:    token,
	}
	return
}

func (c *ControllerV1) Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	tokenStr := r.GetHeader("Authorization")
	if tokenStr == "" {
		err = gerror.NewCode(gcode.CodeNotAuthorized, "请先登录")
		return
	}
	userID, e := jwt.Parse(tokenStr)
	if e != nil {
		err = gerror.NewCode(gcode.CodeNotAuthorized, "Token无效")
		return
	}

	u, wallet, err := user.Profile(ctx, userID)
	if err != nil {
		return
	}
	res = &v1.ProfileRes{
		UserID:   u.Id,
		UID:      u.Uid,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Phone:    u.Phone,
		Chips:    wallet.Chips,
	}
	return
}
