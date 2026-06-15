package table

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"claude-test/internal/model/entity"
)

// SendInvite creates an invitation record and notifies the invitee via WebSocket.
func SendInvite(ctx context.Context, inviterID, tableID, sessionID, inviteeID int64) (invitationID int64, err error) {
	// Verify table exists
	if _, err = GetTable(ctx, tableID); err != nil {
		return
	}

	// Verify invitee exists
	var u entity.Users
	if e := g.DB().Model("users").Where("id", inviteeID).Scan(&u); e != nil || u.Id == 0 {
		err = gerror.New("被邀请用户不存在")
		return
	}

	result, e := g.DB().Model("table_invitations").Data(g.Map{
		"table_id":   tableID,
		"session_id": sessionID,
		"inviter_id": inviterID,
		"invitee_id": inviteeID,
		"status":     1,
		"expired_at": time.Now().Add(30 * time.Minute),
	}).Insert()
	if e != nil {
		err = e
		return
	}
	invitationID, err = result.LastInsertId()
	return
}

// RespondInvite accepts or rejects an invitation.
func RespondInvite(ctx context.Context, userID, invitationID int64, accept bool) (status int, err error) {
	var inv entity.TableInvitations
	if e := g.DB().Model("table_invitations").Where("id", invitationID).Scan(&inv); e != nil || inv.Id == 0 {
		err = gerror.New("邀请不存在")
		return
	}
	if int64(inv.InviteeId) != userID {
		err = gerror.New("无权操作此邀请")
		return
	}
	if inv.Status != 1 {
		err = gerror.New("邀请已处理或已过期")
		return
	}

	status = 3
	if accept {
		status = 2
	}
	_, err = g.DB().Model("table_invitations").Where("id", invitationID).Data(g.Map{"status": status}).Update()
	return
}
