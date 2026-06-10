package table

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "claude-test/api/table/v1"
	tablelogic "claude-test/internal/logic/table"
	"claude-test/utility/jwt"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func currentUserID(ctx context.Context) (int64, error) {
	r := ghttp.RequestFromCtx(ctx)
	token := r.GetHeader("Authorization")
	if token == "" {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "请先登录")
	}
	userID, err := jwt.Parse(token)
	if err != nil {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "Token无效")
	}
	return userID, nil
}

func (c *ControllerV1) CreateTable(ctx context.Context, req *v1.CreateTableReq) (res *v1.CreateTableRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	tableID, tableNo, err := tablelogic.CreateTable(ctx, userID, req)
	if err != nil {
		return
	}
	res = &v1.CreateTableRes{TableID: tableID, TableNo: tableNo}
	return
}

func (c *ControllerV1) JoinTable(ctx context.Context, req *v1.JoinTableReq) (res *v1.JoinTableRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	table, e := tablelogic.GetTable(ctx, req.TableID)
	if e != nil {
		err = e
		return
	}
	if table.HasPassword == 1 && table.Password != req.Password {
		err = gerror.New("密码错误")
		return
	}
	_ = userID
	res = &v1.JoinTableRes{TableID: req.TableID}
	return
}

func (c *ControllerV1) TakeSeat(ctx context.Context, req *v1.TakeSeatReq) (res *v1.TakeSeatRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	err = tablelogic.TakeSeat(ctx, userID, req.TableID, req.SeatNo, req.BuyIn)
	if err != nil {
		return
	}
	res = &v1.TakeSeatRes{SeatNo: req.SeatNo, Chips: req.BuyIn}
	return
}

func (c *ControllerV1) LeaveSeat(ctx context.Context, req *v1.LeaveSeatReq) (res *v1.LeaveSeatRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	err = tablelogic.LeaveSeat(ctx, userID, req.TableID)
	if err != nil {
		return
	}
	res = &v1.LeaveSeatRes{}
	return
}

func (c *ControllerV1) StartSession(ctx context.Context, req *v1.StartSessionReq) (res *v1.StartSessionRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	sessionID, sessionNo, err := tablelogic.StartSession(ctx, req.TableID, userID)
	if err != nil {
		return
	}
	res = &v1.StartSessionRes{SessionID: sessionID, SessionNo: sessionNo}
	return
}

func (c *ControllerV1) BuyIn(ctx context.Context, req *v1.BuyInReq) (res *v1.BuyInRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	recordID, status, err := tablelogic.RebuyRequest(ctx, userID, req.SessionID, req.Amount)
	if err != nil {
		return
	}
	res = &v1.BuyInRes{RecordID: recordID, Status: status}
	return
}

func (c *ControllerV1) ApproveBuyIn(ctx context.Context, req *v1.ApproveBuyInReq) (res *v1.ApproveBuyInRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	err = tablelogic.ApproveBuyin(ctx, userID, req.RecordID, req.Approve)
	if err != nil {
		return
	}
	res = &v1.ApproveBuyInRes{}
	return
}

func (c *ControllerV1) TableRank(ctx context.Context, req *v1.TableRankReq) (res *v1.TableRankRes, err error) {
	_, err = currentUserID(ctx)
	if err != nil {
		return
	}
	data, err := tablelogic.GetTableRank(ctx, req.ID)
	if err != nil {
		return
	}
	_ = data
	// TODO: cast data to TableRankRes properly
	res = &v1.TableRankRes{}
	return
}

func (c *ControllerV1) LobbyTables(ctx context.Context, req *v1.LobbyTablesReq) (res *v1.LobbyTablesRes, err error) {
	list, total, err := tablelogic.LobbyTables(ctx, req.GameType, req.Page, req.PageSize)
	if err != nil {
		return
	}
	res = &v1.LobbyTablesRes{Total: total}
	for _, t := range list {
		res.List = append(res.List, v1.TableInfo{
			TableID:        int64(t.Id),
			TableNo:        t.TableNo,
			Name:           t.Name,
			GameType:       t.GameType,
			SmallBlind:     t.SmallBlind,
			BigBlind:       t.BigBlind,
			MinBuyin:       t.MinBuyin,
			MaxBuyin:       t.MaxBuyin,
			HasPassword:    t.HasPassword == 1,
			CurrentPlayers: int(t.CurrentPlayers),
			MaxSeats:       int(t.MaxSeats),
			Status:         t.Status,
			CreatorID:      int64(t.CreatorId),
		})
	}
	return
}
