package table

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/gogf/gf/v2/frame/g"

	v1 "claude-test/api/table/v1"
	gamelogic "claude-test/internal/logic/game"
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
	// Start game engine asynchronously
	go func() {
		if e := gamelogic.StartGameEngine(ctx, sessionID); e != nil {
			g.Log().Errorf(ctx, "StartGameEngine error: %v", e)
		}
	}()
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
	raw, err := tablelogic.GetTableRank(ctx, req.ID)
	if err != nil {
		return
	}
	m, _ := raw.(g.Map)
	res = &v1.TableRankRes{
		SessionID:  g.NewVar(m["session_id"]).Int64(),
		TotalHands: g.NewVar(m["total_hands"]).Int(),
		TotalFlow:  g.NewVar(m["total_flow"]).Int64(),
		TotalBuyin: g.NewVar(m["total_buyin"]).Int64(),
		AvgPot:     g.NewVar(m["avg_pot"]).Int64(),
	}
	if rows, ok := m["players"]; ok {
		_ = g.NewVar(rows).Scan(&res.Players)
	}
	return
}

func (c *ControllerV1) Invite(ctx context.Context, req *v1.InviteReq) (res *v1.InviteRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	invitationID, err := tablelogic.SendInvite(ctx, userID, req.TableID, req.SessionID, req.InviteeID)
	if err != nil {
		return
	}
	res = &v1.InviteRes{InvitationID: invitationID}
	return
}

func (c *ControllerV1) RespondInvite(ctx context.Context, req *v1.RespondInviteReq) (res *v1.RespondInviteRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	status, err := tablelogic.RespondInvite(ctx, userID, req.InvitationID, req.Accept)
	if err != nil {
		return
	}
	res = &v1.RespondInviteRes{Status: status}
	return
}

func (c *ControllerV1) EndSession(ctx context.Context, req *v1.EndSessionReq) (res *v1.EndSessionRes, err error) {
	_, err = currentUserID(ctx)
	if err != nil {
		return
	}
	err = gamelogic.EndSession(ctx, req.SessionID, req.Reason)
	if err != nil {
		return
	}
	res = &v1.EndSessionRes{SessionID: req.SessionID}
	return
}

func (c *ControllerV1) TableInfoHandler(ctx context.Context, req *v1.TableInfoReq) (res *v1.TableInfoRes, err error) {
	info, e := tablelogic.TableInfo(ctx, req.ID)
	if e != nil {
		err = e
		return
	}
	t := info.Table
	res = &v1.TableInfoRes{
		TableID:   int64(t.Id),
		TableNo:   t.TableNo,
		Name:      t.Name,
		GameType:  t.GameType,
		SmallBlind: t.SmallBlind,
		BigBlind:  t.BigBlind,
		MinBuyin:  t.MinBuyin,
		MaxBuyin:  t.MaxBuyin,
		MaxSeats:  int(t.MaxSeats),
		Duration:  t.Duration,
		Status:    t.Status,
		CreatorID: int64(t.CreatorId),
		SessionID:     info.SessionID,
		StartedAt:     info.StartedAt,
		SessionStatus: info.SessionStatus,
	}
	for _, s := range info.Seats {
		res.Seats = append(res.Seats, v1.SeatInfo{
			SeatNo:   s.SeatNo,
			UserID:   s.UserID,
			Nickname: s.Nickname,
			Avatar:   s.Avatar,
			Chips:    s.Chips,
		})
	}
	return
}

func (c *ControllerV1) AddBots(ctx context.Context, req *v1.AddBotsReq) (res *v1.AddBotsRes, err error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return
	}
	// Only creator can add bots
	table, e := tablelogic.GetTable(ctx, req.ID)
	if e != nil {
		err = e
		return
	}
	if int64(table.CreatorId) != userID {
		err = gerror.NewCode(gcode.CodeNotAuthorized, "仅创建者可添加机器人")
		return
	}
	botIDs, e := tablelogic.AddBots(ctx, req.ID, req.Count)
	if e != nil {
		err = e
		return
	}
	// If the game engine is already running, register bots there too
	gamelogic.RegisterBots(req.ID, botIDs)
	res = &v1.AddBotsRes{Count: len(botIDs)}
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
