package stats

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "claude-test/api/stats/v1"
	statslogic "claude-test/internal/logic/stats"
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

func (c *ControllerV1) SessionList(ctx context.Context, req *v1.SessionListReq) (res *v1.SessionListRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	list, total, err := statslogic.ListSessions(ctx, userID, req)
	if err != nil {
		return
	}
	res = &v1.SessionListRes{List: list, Total: total}
	return
}

func (c *ControllerV1) SessionDetail(ctx context.Context, req *v1.SessionDetailReq) (res *v1.SessionDetailRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	res, err = statslogic.GetSessionDetail(ctx, userID, req.ID)
	return
}

func (c *ControllerV1) HandList(ctx context.Context, req *v1.HandListReq) (res *v1.HandListRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	list, total, err := statslogic.ListHands(ctx, userID, req)
	if err != nil {
		return
	}
	res = &v1.HandListRes{List: list, Total: total}
	return
}

func (c *ControllerV1) HandReplay(ctx context.Context, req *v1.HandReplayReq) (res *v1.HandReplayRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	res, err = statslogic.GetHandReplay(ctx, userID, req.ID)
	return
}

func (c *ControllerV1) Favorite(ctx context.Context, req *v1.FavoriteReq) (res *v1.FavoriteRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	isFav, err := statslogic.ToggleFavorite(ctx, userID, req.ID)
	if err != nil {
		return
	}
	res = &v1.FavoriteRes{IsFavorite: isFav}
	return
}

func (c *ControllerV1) Overview(ctx context.Context, req *v1.OverviewReq) (res *v1.OverviewRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	res, err = statslogic.GetOverview(ctx, userID, req)
	return
}
