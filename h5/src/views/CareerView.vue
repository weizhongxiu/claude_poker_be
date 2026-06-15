<template>
  <div class="page">
    <!-- 顶部 -->
    <header class="header">
      <span class="header-title">好友桌</span>
      <div class="header-right">
        <button class="icon-btn">🔔</button>
      </div>
    </header>

    <!-- 主 Tab -->
    <div class="main-tabs">
      <div
        v-for="t in mainTabs"
        :key="t.key"
        class="main-tab"
        :class="{ active: activeTab === t.key }"
        @click="activeTab = t.key"
      >{{ t.label }}</div>
    </div>

    <!-- 牌局 Tab 内容 -->
    <div v-if="activeTab === 'game'" class="page-body tab-body">
      <!-- 游戏类型 -->
      <div class="chip-row">
        <span
          v-for="g in gameTypes"
          :key="g.key"
          class="chip"
          :class="{ active: activeGame === g.key }"
          @click="activeGame = g.key; load()"
        >{{ g.label }}</span>
      </div>
      <!-- 时间筛选 -->
      <div class="time-row">
        <div class="time-tabs">
          <span
            v-for="t in timePeriods"
            :key="t.key"
            class="time-tab"
            :class="{ active: activeTime === t.key }"
            @click="activeTime = t.key; load()"
          >{{ t.label }}</span>
        </div>
        <!-- 迷你折线图占位 -->
        <div class="mini-chart">
          <svg width="80" height="30" viewBox="0 0 80 30">
            <polyline
              :points="chartPoints"
              fill="none" stroke="#0EC4B0" stroke-width="1.5" stroke-linecap="round"
            />
          </svg>
        </div>
      </div>

      <!-- 战绩总览入口 -->
      <div class="overview-entry" @click="showOverview = !showOverview">
        <span class="oe-label">战绩总览</span>
        <span class="oe-sub" v-if="overview">
          {{ overview.total_profit >= 0 ? '+' : '' }}{{ overview.total_profit }}
        </span>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M9 18l6-6-6-6" stroke="#999" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </div>

      <!-- 展开总览 -->
      <transition name="expand">
        <div v-if="showOverview && overview" class="overview-panel card">
          <div class="ov-grid">
            <div class="ov-item">
              <div class="ov-val" :class="overview.total_profit >= 0 ? 'green' : 'red'">
                {{ overview.total_profit >= 0 ? '+' : '' }}{{ overview.total_profit }}
              </div>
              <div class="ov-key">盈亏</div>
            </div>
            <div class="ov-item">
              <div class="ov-val">{{ overview.total_sessions }}</div>
              <div class="ov-key">场次</div>
            </div>
            <div class="ov-item">
              <div class="ov-val">{{ overview.total_hands }}</div>
              <div class="ov-key">手数</div>
            </div>
            <div class="ov-item">
              <div class="ov-val">{{ overview.total_buyin }}</div>
              <div class="ov-key">总带入</div>
            </div>
            <div class="ov-item">
              <div class="ov-val">{{ (overview.vpip || 0).toFixed(1) }}%</div>
              <div class="ov-key">VPIP</div>
            </div>
            <div class="ov-item">
              <div class="ov-val">{{ (overview.pfr || 0).toFixed(1) }}%</div>
              <div class="ov-key">PFR</div>
            </div>
          </div>
        </div>
      </transition>

      <!-- 历史牌局 -->
      <div class="section-header">
        <span class="section-title">历史牌局</span>
        <span class="section-more" @click="router.push('/tables')">全部牌局 ›</span>
      </div>

      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="sessions.length === 0" class="empty-state">
        <div class="empty-icon">📊</div>
        <p>暂无记录</p>
      </div>
      <div v-else>
        <div
          v-for="s in sessions"
          :key="s.session_id"
          class="session-card card"
          @click="router.push(`/stats/session/${s.session_id}`)"
        >
          <div class="sc-top">
            <span class="sc-no">{{ s.session_no }}</span>
            <span class="sc-profit" :class="s.result >= 0 ? 'green' : 'red'">
              {{ s.result >= 0 ? '+' : '' }}{{ s.result }}
            </span>
          </div>
          <div class="sc-meta">
            <span>{{ s.game_type_label }}</span>
            <span>{{ s.small_blind }}/{{ s.big_blind }}</span>
            <span>{{ s.total_hands }} 手</span>
            <span>{{ s.started_at?.slice(5, 16) }}</span>
          </div>
          <div class="sc-bar">
            <div class="sc-bar-fill" :style="{ width: Math.min(100, Math.abs(s.result) / Math.max(s.total_buyin, 1) * 100) + '%', background: s.result >= 0 ? '#0EC4B0' : '#FF5252' }" />
          </div>
        </div>
      </div>
    </div>

    <!-- 位置 Tab -->
    <div v-else-if="activeTab === 'position'" class="page-body tab-body">
      <div class="chip-row">
        <span
          v-for="g in ['德州','短牌','PLO','SNG']"
          :key="g"
          class="chip"
          :class="{ active: posGame === g }"
          @click="posGame = g"
        >{{ g }}</span>
      </div>
      <div class="pos-hint">
        <span class="info-icon">ⓘ</span>
        <span>仅显示近90天内9人桌数据</span>
      </div>

      <!-- 桌型选择 -->
      <div class="table-size-btns">
        <span
          v-for="sz in ['9人桌','8人桌','6人桌']"
          :key="sz"
          class="size-btn"
          :class="{ active: tableSize === sz }"
          @click="tableSize = sz"
        >{{ sz }}</span>
      </div>

      <!-- 位置图 -->
      <div class="position-diagram">
        <div class="pos-table-oval">
          <div v-for="pos in positionNodes" :key="pos.key" class="pos-node" :style="pos.style">
            <div class="pos-circle">
              <div class="pos-name">{{ pos.abbr }}</div>
              <div class="pos-cn">{{ pos.cn }}</div>
            </div>
            <div class="pos-stat">—</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 牌谱 Tab -->
    <div v-else-if="activeTab === 'hand'" class="page-body tab-body">
      <div v-if="handsLoading" class="loading-state">加载中...</div>
      <div v-else-if="hands.length === 0" class="empty-state">
        <div class="empty-icon">🃏</div>
        <p>暂无牌谱记录</p>
      </div>
      <div v-else>
        <div v-for="hand in hands" :key="hand.game_id"
          class="hand-card card"
          @click="router.push('/replay/' + hand.game_id)">
          <div class="hc-top">
            <span class="hc-id">手牌 #{{ hand.hand_no || hand.game_id }}</span>
            <span class="hc-result" :class="hand.profit >= 0 ? 'green' : 'red'">
              {{ hand.profit >= 0 ? '+' : '' }}{{ hand.profit }}
            </span>
          </div>
          <div class="hc-meta">
            <span>{{ hand.small_blind }}/{{ hand.big_blind }}</span>
            <span>{{ hand.created_at?.slice(5, 16) }}</span>
            <span v-if="hand.is_favorite">⭐</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 其他 Tab 占位 -->
    <div v-else class="page-body tab-body empty-state">
      <div class="empty-icon">🚧</div>
      <p>开发中</p>
    </div>

    <TabBar />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getSessions, getOverview, getHands } from '@/api'
import TabBar from '@/components/TabBar.vue'

const router = useRouter()
const auth = useAuthStore()

const mainTabs = [
  { key: 'game',     label: '牌局' },
  { key: 'hand',     label: '牌谱' },
  { key: 'personal', label: '个人' },
  { key: 'opponent', label: '对手' },
  { key: 'position', label: '位置' },
  { key: 'fortune',  label: '运势' },
  { key: 'cards',    label: '牌组' },
]
const activeTab = ref('game')

const gameTypes = [
  { key: 0, label: '全部' }, { key: 1, label: '德州' },
  { key: 2, label: '短牌' }, { key: 3, label: 'PLO' },
  { key: 4, label: 'SNG' },  { key: 5, label: '13水' },
]
const activeGame = ref(0)

const timePeriods = [
  { key: 'today', label: '今日' }, { key: '7',  label: '7天' },
  { key: '30',    label: '30天' }, { key: '90', label: '90天' },
]
const activeTime = ref('today')

const sessions = ref([])
const overview = ref(null)
const loading = ref(false)
const showOverview = ref(false)
const hands = ref([])
const handsLoading = ref(false)
const posGame = ref('德州')
const tableSize = ref('9人桌')

const gameTypeLabel = { 0:'全部', 1:'德州', 2:'短牌', 3:'PLO', 4:'SNG', 5:'13水' }

const chartPoints = computed(() => {
  if (!sessions.value.length) return '0,15 80,15'
  const vals = sessions.value.slice(0, 8).map(s => s.result)
  const min = Math.min(...vals, 0)
  const max = Math.max(...vals, 1)
  const range = max - min || 1
  return vals.map((v, i) => {
    const x = (i / (vals.length - 1 || 1)) * 80
    const y = 28 - ((v - min) / range) * 24
    return `${x.toFixed(1)},${y.toFixed(1)}`
  }).join(' ')
})

async function load() {
  loading.value = true
  try {
    const now = new Date()
    let dateFrom = ''
    if (activeTime.value === 'today') dateFrom = now.toISOString().slice(0, 10)
    else {
      const d = new Date(now - Number(activeTime.value) * 86400000)
      dateFrom = d.toISOString().slice(0, 10)
    }
    const data = await getSessions({ game_type: activeGame.value, date_from: dateFrom, page: 1, page_size: 20 })
    sessions.value = (data?.list || []).map(s => ({
      ...s, game_type_label: gameTypeLabel[s.game_type] || '—'
    }))
  } catch {
    sessions.value = []
  } finally {
    loading.value = false
  }
}

async function loadOverview() {
  try {
    const now = new Date()
    const today = now.toISOString().slice(0, 10)
    let statType = 4
    let dateFrom = ''
    if (activeTime.value === 'today') {
      statType = 1; dateFrom = today
    } else if (activeTime.value === '7') {
      statType = 2
      dateFrom = new Date(now - 7 * 86400000).toISOString().slice(0, 10)
    } else if (activeTime.value === '30') {
      statType = 3
      dateFrom = new Date(now - 30 * 86400000).toISOString().slice(0, 10)
    } else {
      dateFrom = new Date(now - Number(activeTime.value) * 86400000).toISOString().slice(0, 10)
    }
    const data = await getOverview({
      game_type: activeGame.value,
      stat_type: statType,
      date_from: dateFrom,
      date_to: today,
    })
    overview.value = data
  } catch {}
}

async function loadHands() {
  handsLoading.value = true
  try {
    const data = await getHands({ page: 1, page_size: 50 })
    hands.value = data?.list || []
  } catch {
    hands.value = []
  } finally {
    handsLoading.value = false
  }
}

// 位置图节点
const positionNodes = computed(() => {
  const sz = tableSize.value === '9人桌' ? 9 : tableSize.value === '8人桌' ? 8 : 6
  const nodes9 = [
    { key:'BTN', abbr:'BTN', cn:'庄位',  left:'50%', top:'82%' },
    { key:'SB',  abbr:'SB',  cn:'小盲',  left:'30%', top:'74%' },
    { key:'BB',  abbr:'BB',  cn:'大盲',  left:'14%', top:'58%' },
    { key:'UTG', abbr:'UTG', cn:'枪口',  left:'14%', top:'38%' },
    { key:'U+1', abbr:'UTG+1', cn:'枪口+1', left:'28%', top:'22%' },
    { key:'U+2', abbr:'UTG+2', cn:'枪口+2', left:'50%', top:'14%' },
    { key:'MP',  abbr:'MP',  cn:'中位',  left:'72%', top:'22%' },
    { key:'HJ',  abbr:'MP+1',cn:'中位+1',left:'86%', top:'38%' },
    { key:'CO',  abbr:'CO',  cn:'劫位',  left:'86%', top:'58%' },
  ]
  return nodes9.slice(0, sz).map(n => ({
    ...n,
    style: { position:'absolute', left:n.left, top:n.top, transform:'translate(-50%,-50%)' }
  }))
})

onMounted(async () => {
  await auth.fetchProfile()
  load()
  loadOverview()
})
watch(activeTab, () => {
  if (activeTab.value === 'game') { load(); loadOverview() }
  if (activeTab.value === 'hand') loadHands()
})
watch([activeTime, activeGame], () => {
  load()
  loadOverview()
})
</script>

<style scoped>
.header {
  height: 48px; padding: 0 16px;
  display: flex; align-items: center; justify-content: space-between;
  background: var(--surface); border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.header-title { font-size: 17px; font-weight: 700; }
.icon-btn { background: none; border: none; font-size: 20px; cursor: pointer; }

.main-tabs {
  display: flex; overflow-x: auto; background: var(--surface);
  border-bottom: 1px solid var(--border); flex-shrink: 0;
  scrollbar-width: none;
}
.main-tabs::-webkit-scrollbar { display: none; }
.main-tab {
  flex-shrink: 0; padding: 12px 16px;
  font-size: 14px; color: var(--text-2); cursor: pointer;
  position: relative; transition: color .15s;
}
.main-tab.active {
  color: var(--primary); font-weight: 600;
}
.main-tab.active::after {
  content: ''; position: absolute; bottom: 0; left: 50%; transform: translateX(-50%);
  width: 20px; height: 2.5px; background: var(--primary); border-radius: 2px;
}

.tab-body { padding: 12px 16px; }

.chip-row { display: flex; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }

.time-row {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;
}
.time-tabs { display: flex; gap: 4px; }
.time-tab {
  padding: 4px 12px; border-radius: 14px;
  font-size: 13px; color: var(--text-2); cursor: pointer; transition: all .15s;
}
.time-tab.active { background: var(--primary); color: #fff; font-weight: 600; }
.mini-chart { opacity: .7; }

.overview-entry {
  display: flex; align-items: center; gap: 6px;
  padding: 10px 0; cursor: pointer;
  border-bottom: 1px solid var(--border); margin-bottom: 12px;
}
.oe-label { font-size: 14px; font-weight: 600; flex: 1; }
.oe-sub { font-size: 14px; color: var(--primary); font-weight: 700; }

.overview-panel { padding: 14px; margin-bottom: 14px; }
.ov-grid {
  display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px;
}
.ov-item { text-align: center; }
.ov-val { font-size: 18px; font-weight: 700; margin-bottom: 4px; }
.ov-val.green { color: #0EC4B0; }
.ov-val.red   { color: #FF5252; }
.ov-key { font-size: 11px; color: var(--text-3); }

.section-header {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 10px;
}
.section-title { font-size: 15px; font-weight: 700; }
.section-more  { font-size: 13px; color: var(--text-3); cursor: pointer; }

.loading-state { padding: 40px; text-align: center; color: var(--text-3); }
.empty-state {
  display: flex; flex-direction: column; align-items: center;
  padding: 60px 0; gap: 12px;
}
.empty-icon { font-size: 48px; }
.empty-state p { font-size: 14px; color: var(--text-3); }

.session-card { padding: 14px; margin-bottom: 10px; cursor: pointer; }
.session-card:active { opacity: .85; }
.sc-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px; }
.sc-no  { font-size: 13px; color: var(--text-3); }
.sc-profit { font-size: 20px; font-weight: 700; }
.sc-profit.green { color: #0EC4B0; }
.sc-profit.red   { color: #FF5252; }
.sc-meta {
  display: flex; gap: 12px; font-size: 12px; color: var(--text-3);
  margin-bottom: 8px;
}
.sc-bar { height: 3px; background: var(--bg); border-radius: 2px; overflow: hidden; }
.sc-bar-fill { height: 100%; border-radius: 2px; transition: width .3s; }

/* 位置图 */
.pos-hint {
  display: flex; align-items: center; gap: 6px;
  font-size: 12px; color: var(--text-3); margin-bottom: 12px;
}
.info-icon { font-size: 13px; }
.table-size-btns {
  display: flex; gap: 8px; margin-bottom: 16px;
}
.size-btn {
  padding: 6px 16px; border-radius: 14px;
  font-size: 13px; color: var(--text-2); cursor: pointer;
  border: 1.5px solid var(--border); transition: all .15s;
}
.size-btn.active { border-color: var(--primary); color: var(--primary); background: rgba(14,196,176,.08); }

.position-diagram {
  display: flex; justify-content: center;
}
.pos-table-oval {
  position: relative;
  width: 300px; height: 220px;
  background: radial-gradient(ellipse, #2d6b4a 60%, #1a3d27 100%);
  border-radius: 50%; border: 6px solid #0f2418;
}
.pos-node {
  display: flex; flex-direction: column; align-items: center; gap: 3px;
  white-space: nowrap;
}
.pos-circle {
  background: rgba(255,255,255,.9); border-radius: 50%;
  width: 40px; height: 40px;
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  box-shadow: 0 2px 8px rgba(0,0,0,.3);
}
.pos-name { font-size: 9px; font-weight: 700; color: #1a1a1a; }
.pos-cn   { font-size: 8px; color: var(--text-3); }
.pos-stat { font-size: 11px; color: #fff; font-weight: 600; }

/* expand */
.expand-enter-active, .expand-leave-active { transition: all .2s; overflow: hidden; }
.expand-enter-from, .expand-leave-to { opacity: 0; max-height: 0; }
.expand-enter-to, .expand-leave-from { opacity: 1; max-height: 300px; }
.hand-card { padding: 14px; margin-bottom: 10px; cursor: pointer; }
.hand-card:active { opacity: .85; }
.hc-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 4px; }
.hc-id  { font-size: 13px; color: var(--text-3); }
.hc-result { font-size: 18px; font-weight: 700; }
.hc-result.green { color: #0EC4B0; }
.hc-result.red   { color: #FF5252; }
.hc-meta {
  display: flex; gap: 12px; font-size: 12px; color: var(--text-3);
}
</style>
