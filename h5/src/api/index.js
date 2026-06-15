import axios from 'axios'

const http = axios.create({ baseURL: '/api', timeout: 10000 })

// 自动附带 token
http.interceptors.request.use(cfg => {
  const token = localStorage.getItem('token')
  if (token) cfg.headers.Authorization = token
  return cfg
})

// 统一错误处理
http.interceptors.response.use(
  res => {
    if (res.data.code !== 0) return Promise.reject(new Error(res.data.message || '请求失败'))
    return res.data.data
  },
  err => Promise.reject(err)
)

// ─── 用户 ────────────────────────────────────────────────
export const register = (phone, password, nickname) =>
  http.post('/user/register', { phone, password, nickname })

export const login = (phone, password) =>
  http.post('/user/login', { phone, password })

export const getProfile = () => http.get('/user/profile')

// ─── 牌桌 ────────────────────────────────────────────────
export const createTable = data => http.post('/table/create', data)

export const joinTable = (table_id, password) =>
  http.post('/table/join', { table_id, password })

export const takeSeat = (table_id, seat_no, buyin) =>
  http.post('/table/seat/take', { table_id, seat_no, buyin })

export const leaveSeat = table_id => http.post('/table/seat/leave', { table_id })

export const startSession = table_id => http.post('/table/start', { table_id })

export const endSession = (session_id, reason = 2) =>
  http.post('/table/end', { session_id, reason })

export const buyIn = (session_id, amount) =>
  http.post('/table/buyin', { session_id, amount })

export const approveBuyIn = (record_id, approve) =>
  http.post('/table/buyin/approve', { record_id, approve })

export const inviteFriend = (table_id, session_id, invitee_id) =>
  http.post('/table/invite', { table_id, session_id, invitee_id })

export const respondInvite = (invitation_id, accept) =>
  http.post('/table/invite/respond', { invitation_id, accept })

export const getTableRank = id => http.get(`/table/${id}/rank`)

export const getTableInfo = id => http.get(`/table/${id}/info`)

export const addBots = (id, count) => http.post(`/table/${id}/bots`, { count })

export const getLobbyTables = (game_type = 0, page = 1, page_size = 20) =>
  http.get('/lobby/tables', { params: { game_type, page, page_size } })

// ─── 统计 ────────────────────────────────────────────────
export const getSessions = params => http.get('/stats/sessions', { params })

export const getSessionDetail = id => http.get(`/stats/sessions/${id}`)

export const getHands = params => http.get('/stats/hands', { params })

export const getHandReplay = id => http.get(`/stats/hands/${id}/replay`)

export const toggleFavorite = id => http.post(`/stats/hands/${id}/favorite`)

export const getOverview = params => http.get('/stats/overview', { params })

// ─── 俱乐部 ─────────────────────────────────────────────
export const createClub = (name, logo) => http.post('/club/create', { name, logo })

export const joinClub = club_no => http.post('/club/join', { club_no })

export const getClubInfo = id => http.get(`/club/${id}`)

export const getClubMembers = (id, page = 1) =>
  http.get(`/club/${id}/members`, { params: { page } })

export const getClubTables = (id, page = 1) =>
  http.get(`/club/${id}/tables`, { params: { page } })
