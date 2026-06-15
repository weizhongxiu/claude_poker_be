<template>
  <div class="page">
    <header class="nav">
      <button class="back-btn" @click="router.back()">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M15 18l-6-6 6-6" stroke="#1A1A1A" stroke-width="2.2" stroke-linecap="round"/>
        </svg>
      </button>
      <span class="nav-title">牌局列表</span>
      <button class="filter-btn" @click="showFilter = true">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M3 6h18M7 12h10M11 18h2" stroke="#1A1A1A" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </button>
    </header>

    <div class="filter-tabs">
      <span
        v-for="f in filters"
        :key="f.key"
        class="chip"
        :class="{ active: activeFilter === f.key }"
        @click="activeFilter = f.key"
      >{{ f.label }}</span>
    </div>

    <div class="page-body">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="!tables.length" class="empty-state">
        <div class="empty-icon">🃏</div>
        <p>暂无牌局</p>
      </div>
      <div
        v-for="t in filteredTables"
        :key="t.table_id"
        class="table-item card"
      >
        <div class="ti-avatar">{{ t.name?.[0] || '?' }}</div>
        <div class="ti-info">
          <div class="ti-creator">{{ t.name || '—' }}</div>
          <div class="ti-name">{{ t.name }}</div>
          <div class="ti-tags">
            <span class="ti-tag blind">{{ t.small_blind }}/{{ t.big_blind }}</span>
            <span class="ti-tag">30m</span>
            <span class="ti-tag">{{ t.current_players }}/{{ t.max_seats }}</span>
          </div>
        </div>
        <button class="enter-btn" @click="enter(t)">进局</button>
      </div>
    </div>

    <!-- Filter sheet -->
    <transition name="modal">
      <div v-if="showFilter" class="overlay" @click.self="showFilter = false">
        <div class="sheet">
          <div class="sheet-title">筛选</div>
          <div class="filter-group">
            <div class="fg-label">游戏类型</div>
            <div class="fg-chips">
              <span v-for="g in gameTypes" :key="g.key"
                class="chip" :class="{ active: filterGame === g.key }"
                @click="filterGame = g.key">{{ g.label }}</span>
            </div>
          </div>
          <button class="btn btn-primary" style="width:100%;margin-top:16px"
            @click="applyFilter">确定</button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getLobbyTables } from '@/api'

const router = useRouter()
const loading = ref(false)
const showFilter = ref(false)
const tables = ref([])
const activeFilter = ref('all')
const filterGame = ref(0)

const filters = [
  { key: 'all', label: '全部' },
  { key: 'password', label: '密码局' },
]
const gameTypes = [
  { key: 0, label: '全部' }, { key: 1, label: '德州' },
  { key: 2, label: '短牌' }, { key: 3, label: 'PLO' },
  { key: 4, label: 'SNG' },
]

const filteredTables = computed(() => {
  if (activeFilter.value === 'password') return tables.value.filter(t => t.has_password)
  return tables.value
})

async function load(gameType = 0) {
  loading.value = true
  try {
    const data = await getLobbyTables(gameType, 1, 50)
    tables.value = data?.list || []
  } catch {
    tables.value = []
  } finally {
    loading.value = false
  }
}

function applyFilter() {
  showFilter.value = false
  load(filterGame.value)
}

function enter(t) {
  router.push(`/table/${t.table_id}`)
}

onMounted(() => load())
</script>

<style scoped>
.nav {
  height: 48px; padding: 0 16px;
  display: flex; align-items: center; justify-content: space-between;
  background: var(--surface); border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.back-btn, .filter-btn { background: none; border: none; cursor: pointer; padding: 4px; display: flex; }
.nav-title { font-size: 17px; font-weight: 700; }

.filter-tabs {
  display: flex; gap: 10px; padding: 12px 16px;
  background: var(--surface); border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.loading { padding: 40px; text-align: center; color: var(--text-3); }
.empty-state {
  display: flex; flex-direction: column; align-items: center;
  padding: 60px 0; gap: 12px;
}
.empty-icon { font-size: 48px; }
.empty-state p { font-size: 14px; color: var(--text-3); }

.table-item {
  display: flex; align-items: center; gap: 12px;
  padding: 14px 16px; margin: 10px 16px 0; cursor: pointer;
}
.ti-avatar {
  width: 46px; height: 46px; border-radius: 12px; flex-shrink: 0;
  background: linear-gradient(135deg, #0EC4B0, #0aaa98);
  display: flex; align-items: center; justify-content: center;
  font-size: 18px; font-weight: 700; color: #fff;
}
.ti-info { flex: 1; }
.ti-creator { font-size: 13px; color: var(--text-3); margin-bottom: 2px; }
.ti-name    { font-size: 15px; font-weight: 600; margin-bottom: 5px; }
.ti-tags    { display: flex; gap: 6px; }
.ti-tag     { background: var(--bg); border-radius: 4px; padding: 2px 7px; font-size: 11px; color: var(--text-2); }
.ti-tag.blind { color: var(--primary); background: rgba(14,196,176,.1); }
.enter-btn  {
  background: transparent; color: var(--primary);
  border: 1.5px solid var(--primary);
  border-radius: 16px; height: 34px; padding: 0 16px;
  font-size: 13px; font-weight: 600; cursor: pointer; flex-shrink: 0;
}
.enter-btn:active { background: var(--primary); color: #fff; }

/* Sheet */
.overlay {
  position: fixed; inset: 0; z-index: 500;
  background: rgba(0,0,0,.5); display: flex; align-items: flex-end;
}
.sheet {
  width: 100%; background: var(--surface);
  border-radius: 20px 20px 0 0;
  padding: 20px 16px calc(24px + env(safe-area-inset-bottom));
  animation: slideUp .3s ease;
}
.sheet-title { font-size: 16px; font-weight: 600; text-align: center; margin-bottom: 16px; }
.filter-group { margin-bottom: 12px; }
.fg-label { font-size: 13px; color: var(--text-3); margin-bottom: 8px; }
.fg-chips { display: flex; flex-wrap: wrap; gap: 8px; }
</style>
