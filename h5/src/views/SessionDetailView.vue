<template>
  <div class="page">
    <header class="nav">
      <button class="back-btn" @click="router.back()">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M15 18l-6-6 6-6" stroke="#1A1A1A" stroke-width="2.2" stroke-linecap="round"/>
        </svg>
      </button>
      <span class="nav-title">{{ session?.status === 1 ? '实时战绩' : '牌局结算' }}</span>
      <div style="width:36px" />
    </header>

    <div class="page-body">
      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="!session" class="error-state">
        <p>加载失败</p>
        <button class="btn btn-primary" @click="load">重试</button>
      </div>
      <template v-else>
        <!-- 场次信息卡 -->
        <div class="session-header card">
          <div class="sh-no-row">
            <span class="sh-no">{{ session.session_no }}</span>
            <span v-if="session.status === 1" class="sh-live-badge">进行中</span>
          </div>
          <div class="sh-meta">
            <span class="sh-tag">{{ gameTypeName }}</span>
            <span class="sh-tag">{{ session.small_blind }}/{{ session.big_blind }}</span>
            <span class="sh-tag">{{ session.total_hands }} 手</span>
          </div>
          <div class="sh-time">
            {{ session.started_at?.slice(0, 16) }} — {{ session.ended_at?.slice(11, 16) || '进行中' }}
            <span class="sh-duration">{{ durationLabel }}</span>
          </div>
          <!-- 统计小图 -->
          <div class="sh-stats">
            <div class="sh-stat">
              <div class="sh-stat-val">{{ session.total_buyin }}</div>
              <div class="sh-stat-key">总带入</div>
            </div>
            <div class="sh-stat">
              <div class="sh-stat-val">{{ session.total_flow }}</div>
              <div class="sh-stat-key">总流水</div>
            </div>
            <div class="sh-stat">
              <div class="sh-stat-val">{{ session.max_pot }}</div>
              <div class="sh-stat-key">最大底池</div>
            </div>
            <div class="sh-stat">
              <div class="sh-stat-val">{{ session.avg_pot }}</div>
              <div class="sh-stat-key">平均底池</div>
            </div>
          </div>
        </div>

        <!-- 玩家排行 -->
        <div class="section-header">
          <span class="section-title">{{ session.status === 1 ? '实时排名' : '玩家结算' }}</span>
        </div>
        <div class="players-list">
          <div
            v-for="p in sortedPlayers"
            :key="p.user_id"
            class="player-row card"
          >
            <div class="pr-rank" :class="{ mvp: p.is_mvp }">
              <span v-if="p.is_mvp">MVP</span>
              <span v-else>{{ p.rank }}</span>
            </div>
            <div class="pr-avatar">{{ p.nickname?.[0] || '?' }}</div>
            <div class="pr-info">
              <div class="pr-name">{{ p.nickname }}</div>
              <div class="pr-meta">
                带入 {{ p.total_buyin }} · {{ p.total_hands }} 手 ·
                VPIP {{ p.vpip?.toFixed(0) }}%
              </div>
            </div>
            <div class="pr-result" :class="p.result >= 0 ? 'green' : 'red'">
              {{ p.result >= 0 ? '+' : '' }}{{ p.result }}
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getSessionDetail } from '@/api'

const route  = useRoute()
const router = useRouter()
const loading = ref(false)
const session = ref(null)
let refreshTimer = null

const gameTypeMap = { 1: '德州', 2: '短牌', 3: 'PLO', 4: 'SNG', 5: '十三水' }
const gameTypeName = computed(() => gameTypeMap[session.value?.game_type] || '德州')

const durationLabel = computed(() => {
  const d = session.value?.duration
  if (!d) return ''
  return d < 1 ? `${Math.round(d * 60)}m` : `${d.toFixed(1)}h`
})

const sortedPlayers = computed(() =>
  [...(session.value?.players || [])].sort((a, b) => b.result - a.result)
)

async function load(silent = false) {
  if (!silent) loading.value = true
  try {
    session.value = await getSessionDetail(route.params.id)
    // Auto-refresh every 5s while session is running
    if (session.value?.status === 1 && !refreshTimer) {
      refreshTimer = setInterval(() => load(true), 5000)
    } else if (session.value?.status !== 1 && refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
    }
  } catch {
    if (!silent) session.value = null
  } finally {
    if (!silent) loading.value = false
  }
}

onMounted(() => load())
onUnmounted(() => { if (refreshTimer) clearInterval(refreshTimer) })
</script>

<style scoped>
.nav {
  height: 48px; padding: 0 16px;
  display: flex; align-items: center; justify-content: space-between;
  background: var(--surface); border-bottom: 1px solid var(--border); flex-shrink: 0;
}
.back-btn { background: none; border: none; cursor: pointer; padding: 4px; display: flex; }
.nav-title { font-size: 17px; font-weight: 700; }

.loading-state, .error-state {
  text-align: center; padding: 60px 0; color: var(--text-3); font-size: 14px;
}
.error-state { display: flex; flex-direction: column; align-items: center; gap: 12px; }

.session-header {
  margin: 12px 16px; padding: 16px;
}
.sh-no-row { display: flex; align-items: center; gap: 8px; margin-bottom: 6px; }
.sh-no { font-size: 13px; color: var(--text-3); }
.sh-live-badge {
  font-size: 10px; font-weight: 700; color: #fff;
  background: #0EC4B0; border-radius: 10px; padding: 1px 8px;
  animation: pulse 1.5s ease-in-out infinite;
}
@keyframes pulse { 0%,100% { opacity:1 } 50% { opacity:.6 } }
.sh-meta { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 8px; }
.sh-tag {
  background: var(--bg); border-radius: 4px;
  padding: 2px 8px; font-size: 11px; color: var(--text-2);
}
.sh-time { font-size: 12px; color: var(--text-3); margin-bottom: 12px; }
.sh-duration { margin-left: 8px; color: var(--primary); font-weight: 500; }
.sh-stats { display: grid; grid-template-columns: 1fr 1fr 1fr 1fr; gap: 8px; }
.sh-stat { text-align: center; }
.sh-stat-val { font-size: 16px; font-weight: 700; }
.sh-stat-key { font-size: 10px; color: var(--text-3); margin-top: 2px; }

.section-header { padding: 12px 16px 6px; }
.section-title { font-size: 14px; font-weight: 600; color: var(--text-2); }

.players-list { padding: 0 16px; display: flex; flex-direction: column; gap: 8px; }
.player-row {
  display: flex; align-items: center; gap: 10px; padding: 12px 14px;
}
.pr-rank {
  width: 26px; height: 26px; border-radius: 50%; background: var(--bg);
  display: flex; align-items: center; justify-content: center;
  font-size: 11px; font-weight: 600; color: var(--text-3); flex-shrink: 0;
}
.pr-rank.mvp { background: #FFD700; color: #1a1a1a; font-size: 9px; font-weight: 700; }
.pr-avatar {
  width: 38px; height: 38px; border-radius: 50%;
  background: linear-gradient(135deg, #0EC4B0, #0aaa98);
  display: flex; align-items: center; justify-content: center;
  font-size: 14px; font-weight: 700; color: #fff; flex-shrink: 0;
}
.pr-info { flex: 1; }
.pr-name { font-size: 14px; font-weight: 500; }
.pr-meta { font-size: 11px; color: var(--text-3); margin-top: 2px; }
.pr-result { font-size: 18px; font-weight: 700; }
.green { color: #0EC4B0; }
.red   { color: var(--red); }
</style>
