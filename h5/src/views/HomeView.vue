<template>
  <div class="page">
    <!-- 顶部状态栏 -->
    <header class="header">
      <div class="chips-info" @click="router.push('/profile')">
        <div class="avatar-sm">
          <img v-if="auth.user?.avatar" :src="auth.user.avatar" />
          <span v-else>{{ auth.user?.nickname?.[0] || '?' }}</span>
        </div>
        <div class="chips-wrap">
          <div class="nickname">{{ auth.user?.nickname || '加载中' }}</div>
          <div class="chips-num">🪙 {{ userChips }}</div>
        </div>
      </div>
      <div class="header-right">
        <button class="icon-btn" @click="router.push('/tables')">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
            <circle cx="11" cy="11" r="8" stroke="#1A1A1A" stroke-width="2"/>
            <path d="m21 21-4.35-4.35" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </button>
      </div>
    </header>

    <div class="page-body">
      <!-- 创建/加入 banner -->
      <div class="banner-row">
        <div class="banner banner-create" @click="router.push('/create')">
          <div class="banner-text">
            <div class="banner-title">创建牌局</div>
            <div class="banner-sub">和朋友一起玩</div>
          </div>
          <div class="banner-img">🃏</div>
        </div>
        <div class="banner banner-join" @click="showJoin = true">
          <div class="banner-text">
            <div class="banner-title">加入牌局</div>
            <div class="banner-sub">输入房间号快速入局</div>
          </div>
          <div class="banner-img">🎰</div>
        </div>
      </div>

      <!-- 俱乐部区域 -->
      <div class="club-section card">
        <div class="club-header" @click="router.push('/club')">
          <span class="club-title">俱乐部</span>
          <span class="club-sub" v-if="myClub">{{ myClub.name }}</span>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M9 18l6-6-6-6" stroke="#999" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </div>
        <div v-if="myClub" class="club-info-row">
          <div class="club-stat">
            <div class="cs-val">{{ myClub.member_count }}</div>
            <div class="cs-key">成员</div>
          </div>
          <div class="club-stat">
            <div class="cs-val">{{ myClub.my_chips }}</div>
            <div class="cs-key">俱乐部筹码</div>
          </div>
        </div>
        <div class="club-btns">
          <button class="btn btn-outline" style="flex:1" @click="router.push('/club')">创建/管理</button>
          <button class="btn btn-primary" style="flex:1" @click="showJoinClub = true">加入俱乐部</button>
        </div>
      </div>

      <!-- 牌局列表 -->
      <div class="table-list-section">
        <div class="section-header">
          <div class="filter-tabs">
            <span
              v-for="f in filters"
              :key="f.key"
              class="chip"
              :class="{ active: currentFilter === f.key }"
              @click="currentFilter = f.key; filterTables()"
            >{{ f.label }}</span>
          </div>
          <button class="section-more" @click="router.push('/tables')">更多 ›</button>
        </div>

        <div v-if="loadingTables" class="loading-hint">加载中...</div>
        <div v-else-if="filteredTables.length === 0" class="empty-state">
          <div class="empty-icon">🃏</div>
          <p>暂无进行中的牌局</p>
        </div>

        <div
          v-for="t in filteredTables"
          :key="t.table_id"
          class="table-card card"
          @click="enterTable(t)"
        >
          <div class="tc-avatar">{{ t.name?.[0] || '?' }}</div>
          <div class="table-card-info">
            <div class="tc-name">{{ t.name || `牌桌 #${t.table_id}` }}</div>
            <div class="tc-tags">
              <span class="tc-tag primary">{{ t.small_blind }}/{{ t.big_blind }}</span>
              <span class="tc-tag">{{ gameTypeName(t.game_type) }}</span>
              <span class="tc-tag">{{ t.current_players }}/{{ t.max_seats }} 人</span>
              <span v-if="t.has_password" class="tc-tag lock">🔒</span>
            </div>
          </div>
          <button class="btn-enter" @click.stop="enterTable(t)">进局</button>
        </div>
      </div>
    </div>

    <TabBar />
    <JoinTableModal v-model:visible="showJoin" />

    <!-- 加入俱乐部弹窗 -->
    <transition name="modal">
      <div v-if="showJoinClub" class="overlay" @click.self="showJoinClub = false">
        <div class="sheet">
          <div class="sheet-title">加入俱乐部</div>
          <input v-model="clubNo" class="text-input" placeholder="输入俱乐部编号" />
          <button class="btn btn-primary" style="width:100%;margin-top:14px"
            :disabled="joinClubLoading" @click="doJoinClub">
            {{ joinClubLoading ? '加入中...' : '确认加入' }}
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getLobbyTables, getClubInfo, joinClub } from '@/api'
import TabBar from '@/components/TabBar.vue'
import JoinTableModal from '@/components/JoinTableModal.vue'

const router = useRouter()
const auth   = useAuthStore()

const showJoin      = ref(false)
const showJoinClub  = ref(false)
const loadingTables = ref(false)
const joinClubLoading = ref(false)
const tables        = ref([])
const currentFilter = ref('all')
const myClub        = ref(null)
const clubNo        = ref('')

const filters = [
  { key: 'all',      label: '全部' },
  { key: 'password', label: '密码局' },
]

const gameTypeMap = { 1: '德州', 2: '短牌', 3: 'PLO', 4: 'SNG', 5: '十三水' }
function gameTypeName(t) { return gameTypeMap[t] || '德州' }

const userChips = computed(() => {
  const c = auth.user?.chips
  if (c == null) return '—'
  return c >= 10000 ? (c / 10000).toFixed(1) + 'w' : String(c)
})

const filteredTables = computed(() => {
  if (currentFilter.value === 'password') return tables.value.filter(t => t.has_password)
  return tables.value
})

function filterTables() {}

async function loadTables() {
  loadingTables.value = true
  try {
    const data = await getLobbyTables(0, 1, 20)
    tables.value = data?.list || []
  } catch {
    tables.value = []
  } finally {
    loadingTables.value = false
  }
}

async function loadMyClub() {
  const clubId = Number(localStorage.getItem('my_club_id') || 0)
  if (!clubId) return
  try {
    myClub.value = await getClubInfo(clubId)
  } catch {
    myClub.value = null
  }
}

async function doJoinClub() {
  if (!clubNo.value.trim()) return
  joinClubLoading.value = true
  try {
    const data = await joinClub(clubNo.value.trim())
    localStorage.setItem('my_club_id', String(data.club_id))
    showJoinClub.value = false
    clubNo.value = ''
    await loadMyClub()
  } catch (e) {
    alert(e.message || '加入失败')
  } finally {
    joinClubLoading.value = false
  }
}

function enterTable(t) {
  router.push(`/table/${t.table_id}`)
}

onMounted(async () => {
  await auth.fetchProfile()
  await Promise.all([loadTables(), loadMyClub()])
})
</script>

<style scoped>
.header {
  height: var(--navbar-h); padding: 0 16px;
  display: flex; align-items: center; justify-content: space-between;
  background: var(--surface); border-bottom: 1px solid var(--border); flex-shrink: 0;
}
.chips-info { display: flex; align-items: center; gap: 10px; cursor: pointer; }
.avatar-sm {
  width: 34px; height: 34px; border-radius: 50%;
  background: var(--primary); overflow: hidden;
  display: flex; align-items: center; justify-content: center;
  font-size: 13px; font-weight: 700; color: #fff; flex-shrink: 0;
}
.avatar-sm img { width: 100%; height: 100%; object-fit: cover; }
.nickname { font-size: 13px; font-weight: 500; color: var(--text-2); }
.chips-num { font-size: 15px; font-weight: 700; color: var(--text-1); }
.icon-btn { background: none; border: none; cursor: pointer; padding: 6px; display: flex; }

/* Banner */
.banner-row { padding: 12px 16px; display: flex; flex-direction: column; gap: 10px; }
.banner {
  border-radius: var(--radius-md); height: 72px;
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 20px; cursor: pointer; position: relative; overflow: hidden;
}
.banner:active { opacity: .9; }
.banner-create { background: linear-gradient(135deg, #e0f7f5 0%, #b2ece8 100%); }
.banner-join   { background: linear-gradient(135deg, #e8f5e0 0%, #c8eab2 100%); }
.banner-title { font-size: 17px; font-weight: 700; color: #1A1A1A; }
.banner-sub   { font-size: 12px; color: var(--text-3); margin-top: 2px; }
.banner-img   { font-size: 42px; line-height: 1; }

/* Club */
.club-section { margin: 0 16px 12px; }
.club-header {
  display: flex; align-items: center; gap: 6px;
  padding: 14px 16px 10px; cursor: pointer;
}
.club-title { font-size: 15px; font-weight: 600; flex: 1; }
.club-sub { font-size: 12px; color: var(--primary); }
.club-info-row {
  display: flex; gap: 24px; padding: 0 16px 12px;
}
.club-stat { text-align: center; }
.cs-val { font-size: 18px; font-weight: 700; color: var(--primary); }
.cs-key { font-size: 11px; color: var(--text-3); margin-top: 2px; }
.club-btns { display: flex; gap: 10px; padding: 0 16px 14px; }

/* Table list */
.table-list-section { padding: 0 16px; }
.section-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px;
}
.filter-tabs { display: flex; gap: 8px; }
.section-more { font-size: 13px; color: var(--primary); background: none; border: none; cursor: pointer; }
.loading-hint { text-align: center; padding: 20px; color: var(--text-3); font-size: 13px; }
.empty-state {
  display: flex; flex-direction: column; align-items: center;
  padding: 40px 0; gap: 10px;
}
.empty-icon { font-size: 40px; }
.empty-state p { font-size: 13px; color: var(--text-3); }

.table-card {
  display: flex; align-items: center; gap: 12px;
  padding: 14px 16px; margin-bottom: 10px; cursor: pointer;
}
.table-card:active { opacity: .85; }
.tc-avatar {
  width: 44px; height: 44px; border-radius: 12px;
  background: linear-gradient(135deg, #0EC4B0, #0aaa98);
  display: flex; align-items: center; justify-content: center;
  font-size: 18px; font-weight: 700; color: #fff; flex-shrink: 0;
}
.table-card-info { flex: 1; }
.tc-name { font-size: 15px; font-weight: 600; margin-bottom: 6px; }
.tc-tags { display: flex; gap: 6px; flex-wrap: wrap; }
.tc-tag {
  background: var(--bg); border-radius: 4px;
  padding: 2px 7px; font-size: 11px; color: var(--text-2);
}
.tc-tag.primary { color: var(--primary); background: rgba(14,196,176,.1); }
.tc-tag.lock { background: transparent; padding: 0; }
.btn-enter {
  background: transparent; color: var(--primary);
  border: 1.5px solid var(--primary);
  border-radius: 16px; height: 32px; padding: 0 16px;
  font-size: 13px; font-weight: 600; cursor: pointer;
}
.btn-enter:active { background: var(--primary); color: #fff; }

/* Sheet */
.overlay {
  position: fixed; inset: 0; z-index: 500;
  background: rgba(0,0,0,.5); display: flex; align-items: flex-end;
}
.sheet {
  width: 100%; background: var(--surface);
  border-radius: 20px 20px 0 0;
  padding: 20px 16px calc(28px + env(safe-area-inset-bottom));
  animation: slideUp .3s ease;
}
.sheet-title { font-size: 16px; font-weight: 600; text-align: center; margin-bottom: 16px; }
.text-input {
  width: 100%; height: 48px; border-radius: 12px;
  border: 1.5px solid var(--border); padding: 0 16px;
  font-size: 16px; outline: none; background: var(--bg); color: var(--text-1);
}
.text-input:focus { border-color: var(--primary); }
.modal-enter-active { animation: slideUp .3s ease; }
.modal-leave-active { animation: slideUp .25s ease reverse; }
</style>
