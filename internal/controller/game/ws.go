package game

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"

	"claude-test/internal/game"
	gamelogic "claude-test/internal/logic/game"
	"claude-test/utility/jwt"
	"claude-test/utility/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type ControllerV1 struct{}

func NewV1() *ControllerV1 { return &ControllerV1{} }

// WsTable handles WebSocket connections at /ws/table/{table_id}.
func (c *ControllerV1) WsTable(r *ghttp.Request) {
	tableID := r.GetRouter("table_id").Int64()
	if tableID == 0 {
		r.Response.WriteStatus(400)
		return
	}

	// Auth: JWT from query param or header
	tokenStr := r.GetHeader("Authorization")
	if tokenStr == "" {
		tokenStr = r.GetQuery("token").String()
	}
	if tokenStr == "" {
		r.Response.WriteStatus(401)
		return
	}
	userID, err := jwt.Parse(tokenStr)
	if err != nil {
		r.Response.WriteStatus(401)
		return
	}

	// Upgrade to WebSocket
	conn, e := upgrader.Upgrade(r.Response.Writer, r.Request, nil)
	if e != nil {
		g.Log().Errorf(r.Context(), "ws upgrade error: %v", e)
		return
	}

	// Determine player or observer
	count, _ := g.DB().Model("table_seats").
		Where("table_id", tableID).
		Where("user_id", userID).
		Where("status", 1).
		Count()
	isPlayer := count > 0

	client := ws.GlobalHub().NewClient(userID, tableID, isPlayer, conn)

	// Send current game state after a brief yield so the hub registers the client first.
	go func() {
		time.Sleep(50 * time.Millisecond)
		sendCurrentState(r.Context(), client, tableID, userID)
	}()

	// Read pump blocks until connection closes
	client.ReadPump(func(uid, tid int64, msg []byte) {
		handleClientMessage(r.Context(), uid, tid, msg)
	})
}

// sendCurrentState pushes the current game state (and hole cards if in-hand) to a newly connected client.
func sendCurrentState(ctx context.Context, client *ws.Client, tableID, userID int64) {
	state := game.GlobalEngine().GetState(tableID)
	if state == nil {
		return
	}
	ws.GlobalHub().SendToUser(tableID, userID, ws.MsgTypeGameState, gamelogic.BuildGameStateForClient(state))

	// Re-send hole cards if a hand is in progress so reconnecting players see their cards.
	if state.Stage > 0 {
		for _, p := range state.Players {
			if p.UserID != userID || len(p.HoleCards) == 0 {
				continue
			}
			cards := make([]string, len(p.HoleCards))
			for i, c := range p.HoleCards {
				cards[i] = c.String()
			}
			ws.GlobalHub().SendToUser(tableID, userID, ws.MsgTypeDeal, g.Map{
				"hole_cards": cards,
			})
			break
		}
	}
}

// handleClientMessage processes incoming messages from players.
func handleClientMessage(ctx context.Context, userID, tableID int64, raw []byte) {
	var msg struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}

	switch msg.Type {
	case ws.MsgTypeAction:
		var data struct {
			GameID int64           `json:"game_id"`
			Action json.RawMessage `json:"action"` // accepts int or string
			Amount int64           `json:"amount"`
		}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			return
		}
		actionInt := parseAction(data.Action)
		if actionInt == 0 {
			ws.GlobalHub().SendToUser(tableID, userID, ws.MsgTypeError, g.Map{"msg": "无效行动"})
			return
		}
		// Find seat
		seatNo := getSeatNo(ctx, tableID, userID)
		if seatNo == 0 {
			ws.GlobalHub().SendToUser(tableID, userID, ws.MsgTypeError, g.Map{"msg": "您未在座"})
			return
		}
		err := game.GlobalEngine().SubmitAction(tableID, game.PlayerAction{
			UserID: userID,
			SeatNo: seatNo,
			Action: actionInt,
			Amount: data.Amount,
		})
		if err != nil {
			ws.GlobalHub().SendToUser(tableID, userID, ws.MsgTypeError, g.Map{"msg": err.Error()})
		}

	case ws.MsgTypeChat:
		var data struct {
			SessionID int64  `json:"session_id"`
			MsgType   int    `json:"msg_type"`
			Content   string `json:"content"`
		}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			return
		}
		handleChat(ctx, userID, tableID, data.SessionID, data.MsgType, data.Content)

	case ws.MsgTypePing:
		ws.GlobalHub().SendToUser(tableID, userID, ws.MsgTypePong, nil)
	}
}

// parseAction converts a JSON value (int or string) to a game action constant.
func parseAction(raw json.RawMessage) int {
	var n int
	if err := json.Unmarshal(raw, &n); err == nil && n > 0 {
		return n
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		switch s {
		case "fold":
			return game.ActionFold
		case "check":
			return game.ActionCheck
		case "call":
			return game.ActionCall
		case "raise":
			return game.ActionRaise
		case "bet":
			return game.ActionBet
		case "allin":
			return game.ActionAllIn
		}
	}
	return 0
}

func getSeatNo(ctx context.Context, tableID, userID int64) int {
	val, _ := g.DB().Model("table_seats").
		Fields("seat_no").
		Where("table_id", tableID).
		Where("user_id", userID).
		Where("status", 1).
		Value()
	return val.Int()
}

func handleChat(ctx context.Context, userID, tableID, sessionID int64, msgType int, content string) {
	// Persist chat message
	type userRow struct {
		Nickname string
	}
	var u userRow
	_ = g.DB().Model("users").Fields("nickname").Where("id", userID).Scan(&u)

	if len(content) > 500 {
		content = content[:500]
	}
	_, _ = g.DB().Model("table_messages").Data(g.Map{
		"session_id": sessionID,
		"user_id":    userID,
		"type":       msgType,
		"content":    content,
	}).Insert()

	// Broadcast
	ws.GlobalHub().BroadcastTable(tableID, ws.MsgTypeChat, g.Map{
		"user_id":  userID,
		"nickname": u.Nickname,
		"type":     msgType,
		"content":  content,
	})
}

// StartSessionHandler is the HTTP handler for POST /table/start (also starts game engine).
func StartSessionHandler(ctx context.Context, tableID, sessionID int64) error {
	if tableID == 0 || sessionID == 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "参数错误")
	}
	return gamelogic.StartGameEngine(ctx, sessionID)
}
