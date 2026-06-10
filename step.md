# 德州扑克平台 — 实现大纲

> 基于 design.md v5，记录各模块实现路径、文件结构、关键逻辑
> 最后更新：2026-06-11

---

## 当前状态

```
✅ 已完成          🔧 进行中          ⬜ 待实现
```

| 模块 | 状态 | 说明 |
|------|------|------|
| 环境搭建 | ✅ | Go 1.26.4 / MySQL 8.0.46 / Redis 7.4.9 |
| DAO/Model 生成 | ✅ | 22 张表，88 个文件，注释干净（英文） |
| 用户注册/登录/Profile | ✅ | JWT + Redis Token，3 个接口验证通过 |
| 游戏引擎 | ⬜ | Phase 1 核心 |
| WebSocket Hub | ⬜ | Phase 2 |
| 牌桌 HTTP API | ⬜ | Phase 3 |
| 游戏流程编排 | ⬜ | Phase 4 |
| 牌谱 & 统计 API | ⬜ | Phase 5 |
| 俱乐部模块 | ⬜ | Phase 6 |

---

## Phase 1：游戏引擎

> 路径：`internal/game/`
> 纯 Go，不依赖数据库 / Redis，可独立单测

### 1.1 `deck.go` — 牌组与洗牌

```go
type Card struct {
    Rank int  // 2-14（A=14）
    Suit int  // 0=♠ 1=♥ 2=♦ 3=♣
}

func NewDeck() []Card                             // 52 张标准牌
func Shuffle(deck []Card, seed []byte) []Card     // Fisher-Yates + crypto/rand
func CardToStr(c Card) string                     // "Ah" "Kd" "2c"
func StrToCard(s string) Card
```

### 1.2 `hand_eval.go` — 手牌评估（7 选 5）

```go
type HandResult struct {
    Rank  int      // 内部强度值：1(HighCard,最弱) ~ 10(RoyalFlush,最强)
    Desc  string   // "High Card" / "One Pair" / ... / "Royal Flush"
    Cards []Card   // 最优 5 张
}

func EvalBest5(hole []Card, board []Card) HandResult
    // C(7,5)=21 种组合暴力枚举，取最大值

func evalFive(cards [5]Card) HandResult
    // isRoyalFlush / isStraightFlush / isFourOfAKind /
    // isFullHouse / isFlush / isStraight /
    // isThreeOfAKind / isTwoPair / isOnePair / highCard

func CompareHands(a, b HandResult) int   // -1 / 0 / 1
```

**hand_rank 常量**（与 rules.md 展示排名方向相反，内部强度值越大越强）：

```go
const (
    HighCard      = 1   // 高牌       rules.md 排名: 10（最弱）
    OnePair       = 2   // 一对       rules.md 排名: 9
    TwoPair       = 3   // 两对       rules.md 排名: 8
    ThreeOfAKind  = 4   // 三条       rules.md 排名: 7
    Straight      = 5   // 顺子       rules.md 排名: 6
    Flush         = 6   // 同花       rules.md 排名: 5
    FullHouse     = 7   // 葫芦       rules.md 排名: 4
    FourOfAKind   = 8   // 四条/金刚  rules.md 排名: 3
    StraightFlush = 9   // 同花顺     rules.md 排名: 2
    RoyalFlush    = 10  // 皇家同花顺 rules.md 排名: 1（最强）
)
```

### 1.3 `pot.go` — 底池 / 边池计算

```go
type Pot struct {
    Amount         int64
    EligibleSeats  []int   // 可赢得该底池的座位号
}

func CalcPots(bets map[int]int64, foldedSeats []int) []Pot
    // All-In 分层算法：
    //   1. 筛出全押玩家，按全押金额升序排列
    //   2. 逐层切割：每层 = 当层最小全押额 × 该层参与人数
    //   3. 弃牌玩家投入已进底池，但不参与分配
    //   4. 退还多余筹码（有效注之外的超额部分）

func SplitPot(amount int64, winnerCount int, dealerLeftFirstSeat int) []int64
    // 整除：每人 amount/winnerCount
    // 奇数筹码：庄家左侧第一赢家多得 1 枚
```

**边池示例**：
```
A 全押 100，B 全押 300，C 有 500
  主池  = 100 × 3 = 300   → A / B / C 均可赢
  边池1 = (300-100) × 2 = 400 → B / C 可赢
  退还  = 500 - 300 = 200  → 归还 C
```

### 1.4 `fsm.go` — 手牌状态机

```
Blinds(0) → PreFlop(1) → Flop(2) → Turn(3) → River(4) → Showdown(5)
```

```go
type HandFSM struct {
    state     *GameState
    actionCh  chan PlayerAction
    timerCh   chan struct{}
    callbacks FSMCallbacks    // 数据库落库回调（由 logic/game 注入）
}

func (f *HandFSM) Run(ctx context.Context)
func (f *HandFSM) nextStage()          // 阶段推进，发公共牌
func (f *HandFSM) processAction(a PlayerAction)
    // 合法性校验：轮到该玩家 / 筹码充足 / 动作合法
    // fold / check / call / raise / allin / bet
func (f *HandFSM) isRoundEnd() bool    // 本轮所有人行动完毕
func (f *HandFSM) isHandEnd() bool     // 只剩 1 人（提前结算）
func (f *HandFSM) dealCommunity()      // 翻/转/河 发公共牌
func (f *HandFSM) settle()
    // 1. 计算最优 5 张：EvalBest5
    // 2. 分底池：CalcPots → SplitPot
    // 3. 写回调：OnHandEnd
```

**行动顺序规则**：
- PreFlop：从 UTG（庄家左3）开始顺时针
- Flop / Turn / River：从庄家左侧第一位活跃玩家（SB 或 BB）开始
- 所有人弃牌 / 只剩 1 人 → 提前结束，跳过 Showdown

**特殊场景**：
- Straddle：UTG 主动投入 2×BB，行动权后移
- All-In 后 Run-Twice：双方同意后发两次 Turn+River，分别比牌
- Muck：获胜者可选不亮牌（`is_show_card=0`）

### 1.5 `replay.go` — 回放快照

```go
type StageSnapshot struct {
    Stage          int
    CommunityCards []string
    Pot            int64
    PlayersState   []PlayerSnapshot  // 各座位：chips/bet/status/hole_cards
    ActionSeqStart int
    ActionSeqEnd   int
}

func BuildStageSnapshot(state *GameState, stage int, actionSeqStart, actionSeqEnd int) StageSnapshot
```

### 1.6 `engine.go` — 对外入口

```go
type Engine struct {
    tables map[int64]*TableRoom
    mu     sync.RWMutex
}

func (e *Engine) StartTable(tableID int64, cfg TableConfig)
func (e *Engine) StopTable(tableID int64)
func (e *Engine) SubmitAction(tableID int64, action PlayerAction) error
func (e *Engine) AddPlayer(tableID int64, p PlayerInfo) error
func (e *Engine) RemovePlayer(tableID int64, userID int64) error
```

**单测优先覆盖的边界场景**：
- [ ] 边池分层（3人全押不同金额）
- [ ] Split Pot 奇数筹码分配
- [ ] Run-Twice 双板比牌
- [ ] A-2-3-4-5 低顺子判定
- [ ] 弃牌到只剩 1 人提前结算
- [ ] Straddle 行动顺序变化

---

## Phase 2：WebSocket Hub

> 路径：`utility/ws/hub.go`

```go
type Client struct {
    userID  int64
    tableID int64
    conn    *websocket.Conn
    send    chan []byte
}

type Hub struct {
    tables map[int64]map[int64]*Client   // tableID → userID → Client
    mu     sync.RWMutex
}

func (h *Hub) Register(c *Client)
func (h *Hub) Unregister(c *Client)
func (h *Hub) BroadcastTable(tableID int64, msg []byte)       // 广播全桌
func (h *Hub) SendToUser(tableID, userID int64, msg []byte)   // 私发（底牌）
func (h *Hub) SendToObservers(tableID int64, msg []byte)      // 广播旁观者
```

### 消息类型定义

**服务端 → 客户端**：

| type | 触发时机 | 说明 |
|------|---------|------|
| `game_state` | 每次行动后 | 全桌广播，含底池/公共牌/各玩家状态 |
| `deal` | 发牌时 | 仅私发本人底牌 |
| `action_result` | 玩家行动后 | 广播行动结果 |
| `hand_result` | 本手结算 | 含 split_pot / run_twice / 各赢家手牌 |
| `rank_update` | 每手结束 | 实时排名面板数据 |
| `chat` | 聊天消息 | 广播全桌 |
| `buyin_request` | 补码申请（审核模式） | 推给管理员 |
| `session_end` | 场次结束 | 结算汇总 |

**客户端 → 服务端**：

| type | 说明 |
|------|------|
| `action` | 行动：fold/check/call/raise/allin/bet |
| `chat` | 发送聊天 |
| `ping` | 心跳 |

---

## Phase 3：牌桌 HTTP API

> `api/table/v1/` + `internal/controller/table/` + `internal/logic/table/`

| 方法 | 路径 | 核心逻辑文件 |
|------|------|------------|
| POST | `/table/create` | `logic/table/create.go` |
| POST | `/table/join` | `logic/table/create.go`（验密码，加入旁观） |
| POST | `/table/seat/take` | `logic/table/seat.go` |
| POST | `/table/seat/leave` | `logic/table/seat.go` |
| POST | `/table/start` | `logic/table/create.go` → 启动 Engine |
| POST | `/table/buyin` | `logic/table/buyin.go` |
| POST | `/table/rebuy` | `logic/table/buyin.go` |
| GET  | `/table/{id}/rank` | `logic/table/rank.go` |
| GET  | `/lobby/tables` | `logic/table/lobby.go`（公开桌列表） |
| WS   | `/ws/table/{id}` | `controller/game/ws.go`（升级 WebSocket） |

### 入座关键流程（`logic/table/seat.go`）

```
1. Redis 分布式锁占座（防并发抢座）
   SET table:{tableID}:seat:{seatNo} {userID} NX EX 10

2. BEGIN TRANSACTION
   INSERT table_seats (table_id, user_id, seat_no, chips=buyin)
   UPDATE user_wallets
      SET chips = chips - buyin,
          frozen_chips = frozen_chips + buyin,
          version = version + 1
    WHERE user_id=? AND version=? AND chips >= buyin   ← 乐观锁
   INSERT buyin_records (type=1首次, status=2已批准)
   COMMIT

3. 乐观锁失败（version 变化）→ 重试最多 3 次，仍失败则回滚释放 Redis 锁
4. WebSocket 广播座位更新
```

### 补码流程（`logic/table/buyin.go`）

```
INSERT buyin_records (type=2补码, status=1待审核)

if tables.buyin_approval = 0（免审核）:
    直接 status=2，执行筹码变更：
    UPDATE table_seats SET chips += amount
    UPDATE user_wallets SET chips -= amount, frozen_chips += amount  ← 乐观锁
    UPDATE session_players SET total_buyin += amount
else:
    WS 推送给管理员 buyin_request
    管理员批准 → 执行上方变更
    管理员拒绝 → status=3

累计上限校验（tables.max_buyin_total > 0 时）：
    SELECT SUM(amount) FROM buyin_records
     WHERE session_id=? AND user_id=? AND status=2
    if 累计+本次 > max_buyin_total → 拒绝
```

---

## Phase 4：游戏流程编排

> `internal/logic/game/`，连接引擎 ↔ 数据库 ↔ WebSocket

### `session.go` — 场次生命周期

```go
func StartSession(tableID int64) (*entity.RoomSessions, error)
    // 1. INSERT room_sessions (status=1, started_at=now)
    // 2. engine.StartTable(tableID, cfg)
    // 3. 启动定时器（tables.duration 到期后触发 EndSession）
    // 4. 广播 session_start

func EndSession(sessionID int64, reason int) error
    // reason: 1=时间到 2=手动 3=全部离座
    // 1. UPDATE room_sessions(status=2, ended_at, end_reason)
    // 2. 逐玩家结算（见下方）
    // 3. engine.StopTable
    // 4. 广播 session_end
```

### `hand.go` — 单手牌编排

```go
func StartHand(sessionID int64) error
    // 1. 生成 hand_no（雪花 ID）
    // 2. 生成 shuffle_seed（crypto/rand 32字节，局结束前保密）
    // 3. INSERT games (session_id, hand_no, shuffle_seed, dealer_seat, started_at)
    // 4. INSERT game_players × N（每个在座未离桌玩家）
    // 5. 调用 engine 开始手牌，注入 FSMCallbacks

func OnAction(gameID, userID int64, action int, amount int64) error
    // 1. 合法性校验（轮到该玩家 / 筹码充足 / 动作合法）
    // 2. INSERT game_actions（含 action_ms）
    // 3. 更新 Redis GameState（hash）
    // 4. Hub.BroadcastTable(action_result)
    // 5. 判断本轮是否结束 → 调 engine 推进阶段
```

### `settle.go` — 结算落库

```go
func OnHandEnd(gameID int64, result HandEndResult) error
    // 1. UPDATE game_players（hand_rank/result/is_winner/fold_stage/is_show_card）
    // 2. INSERT pot_distributions（汇总行）
    // 3. INSERT pot_winner_details × 赢家数（Split Pot 每赢家一行，精确金额）
    // 4. UPDATE games(ended_at, duration_ms, community_cards, shuffle_seed 公开, status=2)
    // 5. UPDATE table_seats.chips（每人桌上筹码更新）
    // 6. UPDATE session_players（total_hands / result / vpip / win_rate）
    // 7. UPDATE room_sessions（total_hands / total_flow / avg_pot / max_pot）
    // 8. INSERT hand_replays × 最多5条（preflop/flop/turn/river/showdown 阶段快照）
    // 9. Hub 推送 hand_result（含 split_pot / run_twice 信息）
    // 10. Hub 推送 rank_update
    // 11. 延迟 2s → StartHand（下一手）

func OnSessionSettle(sessionID int64) error
    // 逐玩家：
    //   chips_final = table_seats.chips
    //   total_buyin = SUM(buyin_records.amount WHERE status=2)
    //   result      = chips_final - total_buyin
    //
    // UPDATE session_players(chips_final, result, rank, is_mvp)
    // UPDATE user_wallets(chips += chips_final, frozen_chips -= total_buyin)
    // INSERT chip_transactions(type=3赢 or 4输, amount=|result|)
    // UPDATE user_stats（按 stat_type + stat_date 预聚合）
```

**Redis GameState 结构**（key = `table_state:{tableID}`）：

```go
type GameState struct {
    SessionID      int64
    GameID         int64
    Stage          int                   // 0-5
    Pot            int64
    SidePots       []SidePot
    CommunityCards []string              // ["Ah","Kd","Qc"]
    Players        map[int]*PlayerState  // seat_no → 状态
    CurrentSeat    int
    ActionDeadline int64                 // Unix 毫秒，行动截止时间
    DealerSeat     int
    HandIndex      int                   // 本场第几手（回放进度分子）
    TotalHands     int                   // 本场总手数（回放进度分母）
}
```

---

## Phase 5：牌谱 & 统计 API

> `api/stats/v1/` + `internal/logic/stats/`

| 方法 | 路径 | 数据来源 | 逻辑文件 |
|------|------|---------|---------|
| GET | `/stats/sessions` | `room_sessions + session_players` | `session.go` |
| GET | `/stats/sessions/{id}` | 结算页：`session_players` JOIN `users` | `session.go` |
| GET | `/stats/hands` | `games + game_players` | `replay.go` |
| GET | `/stats/hands/{id}/replay` | `hand_replays + game_actions` | `replay.go` |
| POST | `/stats/hands/{id}/favorite` | `hand_favorites` | `replay.go` |
| GET | `/stats/overview` | `user_stats` | `overview.go` |

### 统计查询策略

```sql
-- 今日（预聚合，stat_type=1）
SELECT * FROM user_stats
 WHERE user_id=? AND game_type=? AND stat_type=1 AND stat_date=CURDATE()

-- 周/月/自定义（实时聚合 session_players）
SELECT
  COUNT(DISTINCT rs.id)   AS sessions,
  SUM(sp.total_hands)     AS hands,
  SUM(sp.result)          AS profit,
  SUM(sp.total_buyin)     AS buyin,
  AVG(sp.vpip)            AS vpip,
  AVG(sp.win_rate)        AS win_rate
FROM session_players sp
JOIN room_sessions rs ON rs.id = sp.session_id
WHERE sp.user_id = ?
  AND rs.game_type = ?
  AND rs.started_at >= ?   -- stat_date
  AND rs.started_at <  ?   -- stat_end（自定义）或 +1day/+7day/+30day

-- 索引保障：room_sessions(started_at) + session_players(user_id, joined_at)
```

### 回放数据组装

```
GET /stats/hands/{id}/replay 返回：
{
  "hand_no": "...",
  "shuffle_seed": "...",   // 局已结束，公开种子（可用于公平性验证）
  "total_stages": 5,
  "stages": [              // 来自 hand_replays
    {
      "stage": 1,          // preflop
      "community_cards": [],
      "pot": 30,
      "players_state": [...],
      "actions": [...]     // 来自 game_actions WHERE action_seq BETWEEN start AND end
    },
    ...
  ]
}
```

---

## Phase 6：俱乐部模块

> `api/club/v1/` + `internal/logic/club/club.go`

| 方法 | 路径 | 逻辑 |
|------|------|------|
| POST | `/club/create` | INSERT clubs + club_members(role=1 创建者)，member_count=1 |
| POST | `/club/join` | INSERT club_members(role=3)，UPDATE clubs.member_count++ |
| POST | `/club/invite` | INSERT club_members(status=pending) or 直接加入 |
| GET  | `/club/{id}/members` | SELECT club_members JOIN users |
| GET  | `/club/{id}/tables` | SELECT tables WHERE club_id=? |

**俱乐部独立筹码**：`club_members.chips` 与 `user_wallets.chips` 完全隔离，俱乐部桌使用俱乐部筹码体系。

---

## 完整文件目录（实现后目标状态）

```
claude-test/
├── main.go                               ✅
├── api/
│   ├── user/v1/user.go                   ✅
│   ├── table/v1/table.go                 ⬜ Phase 3
│   ├── stats/v1/stats.go                 ⬜ Phase 5
│   └── club/v1/club.go                   ⬜ Phase 6
├── internal/
│   ├── cmd/cmd.go                        ✅（需追加路由）
│   ├── game/                             ⬜ Phase 1（游戏引擎，纯逻辑）
│   │   ├── deck.go                       洗牌
│   │   ├── hand_eval.go                  7选5比牌
│   │   ├── pot.go                        主池/边池/split
│   │   ├── fsm.go                        手牌状态机（最复杂）
│   │   ├── replay.go                     回放快照
│   │   └── engine.go                     总入口，管理所有牌桌 goroutine
│   ├── controller/
│   │   ├── user/                         ✅
│   │   ├── table/table.go                ⬜ Phase 3
│   │   ├── game/ws.go                    ⬜ Phase 4（WebSocket 升级入口）
│   │   ├── stats/stats.go                ⬜ Phase 5
│   │   └── club/club.go                  ⬜ Phase 6
│   ├── logic/
│   │   ├── user/user.go                  ✅
│   │   ├── table/
│   │   │   ├── create.go                 ⬜ 建桌/加入/开局
│   │   │   ├── seat.go                   ⬜ 入座/离座（分布式锁+乐观锁）
│   │   │   └── buyin.go                  ⬜ 补码/退码/审核
│   │   ├── game/
│   │   │   ├── session.go                ⬜ 场次开启/结束
│   │   │   ├── hand.go                   ⬜ 单手牌编排（DB ↔ 引擎）
│   │   │   └── settle.go                 ⬜ 结算写库
│   │   ├── stats/
│   │   │   ├── session.go                ⬜ 历史牌局查询
│   │   │   ├── replay.go                 ⬜ 回放数据组装
│   │   │   └── overview.go               ⬜ 生涯统计
│   │   └── club/club.go                  ⬜ Phase 6
│   ├── model/
│   │   ├── entity/                       ✅（22 个文件，gf gen dao 生成）
│   │   └── do/                           ✅（22 个文件，gf gen dao 生成）
│   └── dao/                              ✅（22 个文件，gf gen dao 生成）
├── utility/
│   ├── jwt/jwt.go                        ✅
│   └── ws/hub.go                         ⬜ Phase 2
├── manifest/
│   ├── config/config.yaml                ✅
│   └── sql/poker_schema.sql              ✅（22 张表，英文注释）
└── hack/config.yaml                      ✅（gf gen dao 配置）
```

---

## 开发优先级与时间估算

```
┌─────────────────────────────────────────────────────┐
│ P0 MVP：能打一局牌                                    │
├─────────────────────────────────────────────────────┤
│ Week 1  游戏引擎（单测驱动）                          │
│   deck.go → hand_eval.go → pot.go → fsm.go           │
│                                                     │
│ Week 2  实时对战骨架                                  │
│   ws/hub.go → table create/seat → ws 升级联调        │
│                                                     │
│ Week 3  完整对局闭环                                  │
│   hand.go + settle.go → session → 端到端测试         │
├─────────────────────────────────────────────────────┤
│ P1 完整体验                                          │
├─────────────────────────────────────────────────────┤
│ Week 4  牌谱 & 统计                                  │
│   stats 所有接口 + user_stats 预聚合写入              │
│                                                     │
│ Week 5  社交功能                                     │
│   桌内聊天 + 邀请好友 + 旁观者 + 实时排名面板          │
├─────────────────────────────────────────────────────┤
│ P2 增值功能                                          │
├─────────────────────────────────────────────────────┤
│ Week 6+  俱乐部 / 补码审核 / 活跃度积分               │
│          VPIP/PFR/WTSD 详细统计 / GPS/IP 防作弊       │
├─────────────────────────────────────────────────────┤
│ P3 运营增长                                          │
├─────────────────────────────────────────────────────┤
│          锦标赛（MTT/SNG）/ 暴击玩法 / Run-Twice      │
│          好友系统 / 排行榜 / 充值提现                  │
└─────────────────────────────────────────────────────┘
```

---

## 关键风险点与对策

| 风险 | 对策 |
|------|------|
| **fsm.go 状态机** 最复杂，行动合法性判断边界多 | 先写完整单测用例再实现，覆盖所有边界场景 |
| **边池切割** 多人全押时分层逻辑容易出 bug | pot.go 单独测试，穷举 2/3/4 人全押组合 |
| **Split Pot 奇数筹码** 分配归属 | 严格按"庄家左侧第一赢家多得 1"规则，单测验证 |
| **WebSocket 断线重连** 客户端断开后游戏状态丢失 | 重连时从 Redis GameState 全量下发当前状态 |
| **乐观锁重试** 高并发入座/补码冲突 | 最多重试 3 次，失败返回错误提示用户 |
| **定时器超时弃牌** goroutine 泄漏 | 每手牌用 context.WithDeadline，cancel 及时清理 |

---

## 单测用例清单（Phase 1 优先）

### hand_eval_test.go
- [ ] 皇家同花顺识别（A K Q J T 同花）
- [ ] 低顺子识别（A 2 3 4 5）
- [ ] 7 选 5 最优组合（测试公共牌更强于底牌的情况）
- [ ] 两人比牌：同牌型比点数
- [ ] 两人比牌：相同牌型+点数 → Split Pot

### pot_test.go
- [ ] 2 人全押等额 → 只有主池
- [ ] 3 人全押不同金额 → 主池 + 边池
- [ ] 有人弃牌 → 弃牌玩家不参与分配
- [ ] Split Pot 偶数均分
- [ ] Split Pot 奇数：庄家左侧第一赢家多 1

### fsm_test.go
- [ ] 正常一手牌全流程（preflop→showdown）
- [ ] 所有人弃牌，最后一人赢，不亮牌（Muck）
- [ ] All-In 后跳过后续下注轮直接发牌
- [ ] Straddle 行动顺序变化
- [ ] 行动超时 → 自动弃牌

---

*关联文档：`design.md`（完整技术方案）| `manifest/sql/poker_schema.sql`（DDL）| `rules.md`（德州规则）*
