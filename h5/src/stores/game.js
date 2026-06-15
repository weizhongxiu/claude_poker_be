import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'

export const useGameStore = defineStore('game', () => {
  const connected = ref(false)
  const ws = ref(null)
  let reconnectTimer = null
  let reconnectDelay = 1000
  let currentTableId = null

  const state = reactive({
    tableId: 0,
    tableNo: '',
    sessionId: 0,
    gameId: 0,
    stage: 0,           // 0=等待 1=preflop 2=flop 3=turn 4=river 5=showdown
    smallBlind: 0,
    bigBlind: 0,
    pot: 0,
    sidePots: [],
    communityCards: [],
    players: {},        // seatNo → PlayerState
    currentSeat: -1,
    actionDeadline: 0,  // unix ms timestamp
    handIndex: 0,
    totalHands: 0,
    dealerSeat: -1,
    isCreator: false,
    sessionStatus: 0,   // 0=waiting 1=running 2=ended
  })

  const myHoleCards    = ref([])
  const lastDealCards  = ref([]) // snapshot saved at deal time, persists through hand_result
  const showdownCards  = ref([])  // [{seat_no, hole_cards, hand_rank, hand_desc}] during showdown
  const lastResult     = ref(null)
  const lastError      = ref('')
  const messages       = ref([])
  const rankList       = ref([])
  const seatLastAction = ref({}) // seatNo → { label, type } for action badge display

  // YOU WIN 动画期间缓冲下一手牌消息
  let _holdNewHand = false
  const _pendingMsgs = []

  function holdNewHand() { _holdNewHand = true }
  function flushPendingHand() {
    _holdNewHand = false
    const msgs = _pendingMsgs.splice(0)
    msgs.forEach(m => handleMsg(m))
  }

  // ── WebSocket ──────────────────────────────────────────────────

  function connect(tableId) {
    currentTableId = tableId
    _createWs(tableId)
  }

  function _createWs(tableId) {
    if (ws.value) {
      ws.value.onclose = null
      ws.value.close()
    }
    const token = localStorage.getItem('token') || ''
    const proto = location.protocol === 'https:' ? 'wss' : 'ws'
    const url = `${proto}://${location.host}/ws/table/${tableId}?token=${token}`
    const socket = new WebSocket(url)
    ws.value = socket

    socket.onopen = () => {
      connected.value = true
      reconnectDelay = 1000
    }
    socket.onclose = () => {
      connected.value = false
      _scheduleReconnect()
    }
    socket.onerror = () => {
      connected.value = false
    }
    socket.onmessage = e => {
      try { handleMsg(JSON.parse(e.data)) } catch {}
    }
  }

  function _scheduleReconnect() {
    clearTimeout(reconnectTimer)
    if (!currentTableId) return
    reconnectTimer = setTimeout(() => {
      if (!connected.value && currentTableId) {
        reconnectDelay = Math.min(reconnectDelay * 2, 30000)
        _createWs(currentTableId)
      }
    }, reconnectDelay)
  }

  function disconnect() {
    currentTableId = null
    clearTimeout(reconnectTimer)
    if (ws.value) {
      ws.value.onclose = null
      ws.value.close()
      ws.value = null
    }
    connected.value = false
    // reset state
    Object.assign(state, {
      tableId: 0, tableNo: '', sessionId: 0, gameId: 0,
      stage: 0, pot: 0, communityCards: [], players: {},
      currentSeat: -1, dealerSeat: -1, sessionStatus: 0,
      handIndex: 0, totalHands: 0, actionDeadline: 0,
    })
    myHoleCards.value = []
    lastResult.value = null
    messages.value = []
    rankList.value = []
  }

  function send(type, data) {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify({ type, data }))
    }
  }

  function sendAction(gameId, action, amount = 0) {
    send('action', { game_id: gameId, action, amount })
  }

  function sendChat(sessionId, content, msgType = 1) {
    send('chat', { session_id: sessionId, msg_type: msgType, content })
    messages.value.push({ type: 'mine', content, time: Date.now() })
  }

  // ── Message Handler ────────────────────────────────────────────

  function handleMsg(msg) {
    const { type, data } = msg

    // 新手牌开始的消息（game_state hand_index变化 或 deal）在 YOU WIN 期间缓冲
    if (_holdNewHand) {
      const isNewHand = (type === 'game_state' && (data.hand_index || 0) > state.handIndex)
                     || type === 'deal'
      if (isNewHand) { _pendingMsgs.push(msg); return }
    }

    switch (type) {
      case 'table_info':
        // 桌信息：table_no, small_blind, big_blind, session_id 等
        if (data.table_no)   state.tableNo    = data.table_no
        if (data.table_id)   state.tableId    = data.table_id
        if (data.small_blind) state.smallBlind = data.small_blind
        if (data.big_blind)  state.bigBlind   = data.big_blind
        if (data.session_id) state.sessionId  = data.session_id
        if (data.session_status != null) state.sessionStatus = data.session_status
        break

      case 'session_started':
        if (data.session_id) state.sessionId = data.session_id
        state.sessionStatus = 1
        state.totalHands    = 0
        break

      case 'session_end':
        // Only mark ended if it matches the current session (avoid stale broadcasts)
        if (!data.session_id || data.session_id === state.sessionId) {
          state.sessionStatus = 2
        }
        break

      case 'player_joined':
      case 'player_left':
      case 'seat_update': {
        // 重刷 players 对应座位
        const p = data.player || data
        if (p && p.seat != null) {
          if (type === 'player_left') {
            const copy = { ...state.players }
            delete copy[p.seat]
            state.players = copy
          } else {
            state.players = { ...state.players, [p.seat]: p }
          }
        }
        break
      }

      case 'game_state':
        seatLastAction.value = {}
        Object.assign(state, {
          gameId:         data.game_id   || state.gameId,
          sessionId:      data.session_id || state.sessionId,
          stage:          stageNum(data.stage),
          pot:            data.pot       || 0,
          sidePots:       data.side_pots || [],
          communityCards: data.community_cards || [],
          currentSeat:    data.current_seat ?? -1,
          actionDeadline: data.deadline  || 0,
          handIndex:      data.hand_index || state.handIndex,
          totalHands:     data.total_hands || state.totalHands,
          dealerSeat:     data.dealer_seat ?? state.dealerSeat,
          smallBlind:     data.small_blind || state.smallBlind,
          bigBlind:       data.big_blind   || state.bigBlind,
          players:        {}
        })
        ;(data.players || []).forEach(p => {
          state.players[p.seat] = p
        })
        break

      case 'deal':
        myHoleCards.value   = data.hole_cards || []
        lastDealCards.value = data.hole_cards || [] // keep for YOU WIN overlay
        break

      case 'action_result': {
        const actionLabels = { 1:'弃牌', 2:'过牌', 3:'跟注', 4:'加注', 5:'全押', 6:'下注', 7:'小盲', 8:'前注', 9:'跨骑' }
        const actionTypes  = { 1:'fold', 2:'check', 3:'call', 4:'raise', 5:'allin', 6:'bet', 7:'blind', 8:'ante', 9:'straddle' }
        if (data.seat != null) {
          const label = actionLabels[data.action] || ''
          const atype = actionTypes[data.action] || ''
          const showAmount = data.amount > 0 && ![1,2,7,8].includes(data.action)
          seatLastAction.value = {
            ...seatLastAction.value,
            [data.seat]: { label: showAmount ? `${label} ${data.amount}` : label, type: atype }
          }
          if (state.players[data.seat]) {
            const isFold = data.action === 1
            state.players[data.seat] = {
              ...state.players[data.seat],
              bet:    isFold ? 0 : (state.players[data.seat].bet || 0) + (data.amount || 0),
              chips:  state.players[data.seat].chips - (data.amount || 0),
              status: isFold ? 2 : state.players[data.seat].status,
            }
          }
        }
        if (data.pot != null)          state.pot            = data.pot
        if (data.current_seat != null) state.currentSeat    = data.current_seat
        if (data.deadline != null)     state.actionDeadline = data.deadline
        break
      }

      case 'error':
        console.warn('[WS error]', data?.msg)
        lastError.value = data?.msg || '操作失败'
        break

      case 'new_cards':
        // 新翻公共牌
        state.communityCards = data.community_cards || state.communityCards
        state.currentSeat    = data.current_seat ?? state.currentSeat
        state.actionDeadline = data.deadline || 0
        break

      case 'showdown':
        showdownCards.value = (data.players || []).map(p => ({
          seatNo:    p.seat_no,
          holeCards: p.hole_cards ? p.hole_cards.trim().split(/\s+/).filter(Boolean) : [],
          handRank:  p.hand_rank,
          handDesc:  p.hand_desc,
          isWinner:  p.is_winner,
          winAmount: p.win_amount,
        }))
        break

      case 'hand_result':
        showdownCards.value  = []
        seatLastAction.value = {}
        state.currentSeat    = -1   // ← clear BEFORE setting lastResult so watcher sees -1
        lastResult.value     = data // ← watcher fires after this tick
        myHoleCards.value    = []
        state.handIndex   = (state.handIndex || 0) + 1
        state.totalHands  = (state.totalHands || 0) + 1
        // Update each player's chip count immediately
        ;(data.players || []).forEach(p => {
          if (state.players[p.seat_no]) {
            state.players[p.seat_no] = { ...state.players[p.seat_no], chips: p.chips_end, bet: 0 }
          }
        })
        break

      case 'rank_update':
        rankList.value = data.players || []
        break

      case 'chat':
        messages.value.push({ type: 'other', ...data })
        break

      case 'pong':
        break
    }
  }

  function stageNum(s) {
    if (typeof s === 'number') return s
    return { preflop: 1, flop: 2, turn: 3, river: 4, showdown: 5 }[s] || 0
  }

  // ── Heartbeat ─────────────────────────────────────────────────
  let pingTimer = null
  function startPing() {
    pingTimer = setInterval(() => send('ping', {}), 20000)
  }
  function stopPing() {
    clearInterval(pingTimer)
  }

  return {
    connected, state, myHoleCards, lastDealCards, showdownCards, lastResult, lastError, messages, rankList, seatLastAction,
    connect, disconnect, send, sendAction, sendChat, startPing, stopPing,
    holdNewHand, flushPendingHand,
  }
})
