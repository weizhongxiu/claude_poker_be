package club

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "claude-test/api/club/v1"
	clublogic "claude-test/internal/logic/club"
	"claude-test/utility/jwt"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 { return &ControllerV1{} }

func currentUserID(ctx context.Context) (int64, error) {
	r := ghttp.RequestFromCtx(ctx)
	token := r.GetHeader("Authorization")
	if token == "" {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "请先登录")
	}
	uid, err := jwt.Parse(token)
	if err != nil {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "Token无效")
	}
	return uid, nil
}

func (c *ControllerV1) CreateClub(ctx context.Context, req *v1.CreateClubReq) (res *v1.CreateClubRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	clubID, clubNo, err := clublogic.CreateClub(ctx, userID, req)
	if err != nil {
		return
	}
	res = &v1.CreateClubRes{ClubID: clubID, ClubNo: clubNo}
	return
}

func (c *ControllerV1) JoinClub(ctx context.Context, req *v1.JoinClubReq) (res *v1.JoinClubRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	clubID, err := clublogic.JoinClub(ctx, userID, req.ClubNo)
	if err != nil {
		return
	}
	res = &v1.JoinClubRes{ClubID: clubID}
	return
}

func (c *ControllerV1) ClubInfo(ctx context.Context, req *v1.ClubInfoReq) (res *v1.ClubInfoRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	res, err = clublogic.GetClubInfo(ctx, userID, req.ID)
	return
}

func (c *ControllerV1) MemberList(ctx context.Context, req *v1.MemberListReq) (res *v1.MemberListRes, err error) {
	_, err = currentUserID(ctx)
	if err != nil {
		return
	}
	list, total, err := clublogic.ListMembers(ctx, req.ID, req.Page, req.PageSize)
	if err != nil {
		return
	}
	res = &v1.MemberListRes{List: list, Total: total}
	return
}

func (c *ControllerV1) ClubTables(ctx context.Context, req *v1.ClubTablesReq) (res *v1.ClubTablesRes, err error) {
	_, err = currentUserID(ctx)
	if err != nil {
		return
	}
	list, total, err := clublogic.ListClubTables(ctx, req.ID, req.Page, req.PageSize)
	if err != nil {
		return
	}
	res = &v1.ClubTablesRes{List: list, Total: total}
	return
}
