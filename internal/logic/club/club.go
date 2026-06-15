package club

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "claude-test/api/club/v1"
)

// CreateClub creates a new club and sets the creator as owner (role=1).
func CreateClub(ctx context.Context, ownerID int64, req *v1.CreateClubReq) (clubID int64, clubNo string, err error) {
	clubNo = fmt.Sprintf("C%d", time.Now().UnixMilli()%1000000000)

	result, e := g.DB().Model("clubs").Data(g.Map{
		"club_no":      clubNo,
		"name":         req.Name,
		"logo":         req.Logo,
		"owner_id":     ownerID,
		"member_count": 1,
		"status":       1,
	}).Insert()
	if e != nil {
		err = e
		return
	}
	clubID, err = result.LastInsertId()
	if err != nil {
		return
	}

	// Add owner as member (role=1)
	_, err = g.DB().Model("club_members").Data(g.Map{
		"club_id": clubID,
		"user_id": ownerID,
		"role":    1,
		"status":  1,
	}).Insert()
	return
}

// JoinClub adds a user to a club (role=3 member).
func JoinClub(ctx context.Context, userID int64, clubNo string) (clubID int64, err error) {
	type clubRow struct {
		Id     uint64
		Status int
	}
	var club clubRow
	if e := g.DB().Model("clubs").Fields("id,status").Where("club_no", clubNo).Scan(&club); e != nil || club.Id == 0 {
		err = gerror.New("俱乐部不存在")
		return
	}
	if club.Status != 1 {
		err = gerror.New("俱乐部已解散")
		return
	}
	clubID = int64(club.Id)

	// Check already member
	count, _ := g.DB().Model("club_members").
		Where("club_id", clubID).
		Where("user_id", userID).
		Count()
	if count > 0 {
		err = gerror.New("您已是该俱乐部成员")
		return
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, e := tx.Model("club_members").Data(g.Map{
			"club_id": clubID,
			"user_id": userID,
			"role":    3,
			"status":  1,
		}).Insert()
		if e != nil {
			return e
		}
		_, e = tx.Model("clubs").Where("id", clubID).
			Data(g.Map{"member_count": gdb.Raw("member_count + 1")}).Update()
		return e
	})
	return
}

// GetClubInfo returns club info with current user's membership context.
func GetClubInfo(ctx context.Context, userID, clubID int64) (*v1.ClubInfoRes, error) {
	type clubRow struct {
		Id           uint64
		ClubNo       string
		Name         string
		Logo         string
		MemberCount  int
		Announcement string
	}
	var club clubRow
	if e := g.DB().Model("clubs").Where("id", clubID).Where("status", 1).Scan(&club); e != nil || club.Id == 0 {
		return nil, gerror.New("俱乐部不存在")
	}

	type memberRow struct {
		Role  int
		Chips int64
	}
	var mem memberRow
	_ = g.DB().Model("club_members").Fields("role,chips").
		Where("club_id", clubID).
		Where("user_id", userID).
		Where("status", 1).
		Scan(&mem)

	return &v1.ClubInfoRes{
		ClubID:       clubID,
		ClubNo:       club.ClubNo,
		Name:         club.Name,
		Logo:         club.Logo,
		MemberCount:  club.MemberCount,
		Announcement: club.Announcement,
		MyRole:       mem.Role,
		MyChips:      mem.Chips,
	}, nil
}

// ListMembers returns paginated member list for a club.
func ListMembers(ctx context.Context, clubID int64, page, pageSize int) (list []v1.ClubMember, total int, err error) {
	m := g.DB().Model("club_members cm").
		LeftJoin("users u", "u.id = cm.user_id").
		Fields("cm.user_id, u.nickname, u.avatar, cm.role, cm.chips, cm.joined_at").
		Where("cm.club_id", clubID).
		Where("cm.status", 1)

	total, err = m.Count()
	if err != nil {
		return
	}

	type row struct {
		UserID   int64       `json:"user_id"`
		Nickname string      `json:"nickname"`
		Avatar   string      `json:"avatar"`
		Role     int         `json:"role"`
		Chips    int64       `json:"chips"`
		JoinedAt *gtime.Time `json:"joined_at"`
	}
	var rows []*row
	err = m.OrderAsc("cm.role").OrderDesc("cm.chips").Page(page, pageSize).Scan(&rows)
	if err != nil {
		return
	}

	for _, r := range rows {
		item := v1.ClubMember{
			UserID:   r.UserID,
			Nickname: r.Nickname,
			Avatar:   r.Avatar,
			Role:     r.Role,
			Chips:    r.Chips,
		}
		if r.JoinedAt != nil {
			item.JoinedAt = r.JoinedAt.Format("Y-m-d H:i:s")
		}
		list = append(list, item)
	}
	return
}

// ListClubTables returns tables belonging to a club that have an active, non-expired session.
func ListClubTables(ctx context.Context, clubID int64, page, pageSize int) (list []v1.ClubTableItem, total int, err error) {
	// Only show tables with a running session that hasn't exceeded its time limit
	m := g.DB().Model("tables t").
		InnerJoin("room_sessions rs",
			"rs.table_id = t.id AND rs.status = 1 AND (rs.duration = 0 OR DATE_ADD(rs.started_at, INTERVAL rs.duration HOUR) > NOW())").
		Where("t.club_id", clubID).
		Where("t.status", 2)

	total, err = m.Count()
	if err != nil {
		return
	}

	type row struct {
		Id             uint64 `json:"id"`
		TableNo        string `json:"table_no"`
		Name           string `json:"name"`
		GameType       int    `json:"game_type"`
		SmallBlind     int64  `json:"small_blind"`
		BigBlind       int64  `json:"big_blind"`
		CurrentPlayers int    `json:"current_players"`
		MaxSeats       int    `json:"max_seats"`
		Status         int    `json:"status"`
	}
	var rows []*row
	err = m.Fields("t.id, t.table_no, t.name, t.game_type, t.small_blind, t.big_blind, t.current_players, t.max_seats, t.status").
		Page(page, pageSize).OrderDesc("t.id").Scan(&rows)
	if err != nil {
		return
	}

	for _, r := range rows {
		list = append(list, v1.ClubTableItem{
			TableID:        int64(r.Id),
			TableNo:        r.TableNo,
			Name:           r.Name,
			GameType:       r.GameType,
			SmallBlind:     r.SmallBlind,
			BigBlind:       r.BigBlind,
			CurrentPlayers: r.CurrentPlayers,
			MaxSeats:       r.MaxSeats,
			Status:         r.Status,
		})
	}
	return
}
