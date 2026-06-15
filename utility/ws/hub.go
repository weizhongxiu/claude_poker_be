package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 65536
)

// Message types for server → client push.
const (
	MsgTypeGameState    = "game_state"    // Full table state broadcast after each action
	MsgTypeDeal         = "deal"          // Private hole cards (sent to player only)
	MsgTypeActionResult = "action_result" // Action result broadcast
	MsgTypeHandResult   = "hand_result"   // Hand settlement (split pot / run-twice)
	MsgTypeRankUpdate   = "rank_update"   // Real-time ranking panel
	MsgTypeChat         = "chat"          // Chat message broadcast
	MsgTypeBuyinRequest = "buyin_request" // Rebuy approval request (for admins)
	MsgTypeShowdown     = "showdown"      // Showdown: reveal all hole cards
	MsgTypeSessionStarted = "session_started" // Session started
	MsgTypeSessionEnd     = "session_end"     // Session ended
	MsgTypeChipUpdate     = "chip_update"     // Chip count changed (rebuy)
	MsgTypeError          = "error"           // Error notification
	MsgTypePong           = "pong"            // Heartbeat response

	// Client → server
	MsgTypeAction = "action" // Player game action
	MsgTypePing   = "ping"   // Heartbeat
)

// Message is the standard WebSocket envelope.
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// Client represents one connected WebSocket client.
type Client struct {
	UserID    int64
	TableID   int64
	IsPlayer  bool // false = observer
	conn      *websocket.Conn
	send      chan []byte
	hub       *Hub
}

// Hub manages all WebSocket connections, grouped by table.
type Hub struct {
	// tableID → userID → *Client
	tables map[int64]map[int64]*Client
	mu     sync.RWMutex

	register   chan *Client
	unregister chan *Client
	broadcast  chan tableMsg
}

type tableMsg struct {
	tableID  int64
	userID   int64 // 0 = broadcast to all
	data     []byte
}

var globalHub = NewHub()

// GlobalHub returns the singleton hub instance.
func GlobalHub() *Hub {
	return globalHub
}

// NewHub creates a new Hub and starts its event loop.
func NewHub() *Hub {
	h := &Hub{
		tables:     make(map[int64]map[int64]*Client),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		broadcast:  make(chan tableMsg, 256),
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			if h.tables[c.TableID] == nil {
				h.tables[c.TableID] = make(map[int64]*Client)
			}
			h.tables[c.TableID][c.UserID] = c
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.tables[c.TableID]; ok {
				if _, ok := clients[c.UserID]; ok {
					delete(clients, c.UserID)
					close(c.send)
					if len(clients) == 0 {
						delete(h.tables, c.TableID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.tables[msg.tableID]
			if msg.userID == 0 {
				// Broadcast to all
				for _, c := range clients {
					select {
					case c.send <- msg.data:
					default:
						// Slow client: drop
					}
				}
			} else {
				// Send to specific user
				if c, ok := clients[msg.userID]; ok {
					select {
					case c.send <- msg.data:
					default:
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register adds a client to the hub.
func (h *Hub) Register(c *Client) {
	h.register <- c
}

// Unregister removes a client from the hub.
func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}

// BroadcastTable sends a message to all clients at a table.
func (h *Hub) BroadcastTable(tableID int64, msgType string, data interface{}) {
	b, _ := json.Marshal(Message{Type: msgType, Data: data})
	h.broadcast <- tableMsg{tableID: tableID, data: b}
}

// SendToUser sends a message to a specific user at a table.
func (h *Hub) SendToUser(tableID, userID int64, msgType string, data interface{}) {
	b, _ := json.Marshal(Message{Type: msgType, Data: data})
	h.broadcast <- tableMsg{tableID: tableID, userID: userID, data: b}
}

// NewClient creates a client and starts its read/write pumps.
func (h *Hub) NewClient(userID, tableID int64, isPlayer bool, conn *websocket.Conn) *Client {
	c := &Client{
		UserID:   userID,
		TableID:  tableID,
		IsPlayer: isPlayer,
		conn:     conn,
		send:     make(chan []byte, 256),
		hub:      h,
	}
	h.Register(c)
	go c.writePump()
	return c
}

// ReadPump reads messages from the WebSocket and calls handler.
// Should be called in a goroutine; returns when connection closes.
func (c *Client) ReadPump(handler func(userID, tableID int64, msg []byte)) {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		handler(c.UserID, c.TableID, msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
