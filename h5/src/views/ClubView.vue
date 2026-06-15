<template>
  <div class="page">
    <header class="nav">
      <button class="back-btn" @click="router.back()">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M15 18l-6-6 6-6" stroke="#1A1A1A" stroke-width="2.2" stroke-linecap="round"/>
        </svg>
      </button>
      <span class="nav-title">俱乐部</span>
      <button class="join-btn" @click="showJoin = true">加入</button>
    </header>

    <div class="page-body">
      <!-- 当前俱乐部信息 -->
      <div v-if="currentClub" class="club-card card">
        <div class="club-logo">{{ currentClub.name?.[0] }}</div>
        <div class="club-info">
          <div class="club-name">{{ currentClub.name }}</div>
          <div class="club-meta">
            <span>编号: {{ currentClub.club_no }}</span>
            <span>{{ currentClub.member_count }} 成员</span>
            <span>我的角色: {{ roleLabel }}</span>
          </div>
          <div class="club-chips">俱乐部筹码: {{ currentClub.my_chips }}</div>
        </div>
      </div>

      <!-- 俱乐部桌 -->
      <div v-if="currentClub" class="section-header" style="margin-top:8px">
        <span class="section-title">牌桌</span>
        <button class="create-table-btn" @click="router.push('/create')">+ 创建</button>
      </div>
      <div v-if="tables.length === 0 && currentClub" class="empty-state">
        <div class="empty-icon">🃏</div>
        <p>暂无牌桌</p>
      </div>
      <div v-for="t in tables" :key="t.table_id" class="table-item card" @click="router.push(`/table/${t.table_id}`)">
        <div class="ti-avatar">{{ t.name?.[0] || '?' }}</div>
        <div class="ti-info">
          <div class="ti-name">{{ t.name }}</div>
          <div class="ti-tags">
            <span class="ti-tag">{{ t.small_blind }}/{{ t.big_blind }}</span>
            <span class="ti-tag">{{ t.current_players }}/{{ t.max_seats }}</span>
          </div>
        </div>
        <button class="btn btn-outline" style="height:34px;padding:0 14px;font-size:13px">进局</button>
      </div>

      <!-- 未加入俱乐部 -->
      <div v-if="!currentClub && !loading" class="empty-hero">
        <div class="empty-icon">🏆</div>
        <p class="empty-text">还没有加入俱乐部</p>
        <div class="empty-btns">
          <button class="btn btn-primary" @click="showCreate = true">创建俱乐部</button>
          <button class="btn btn-outline" @click="showJoin = true">加入俱乐部</button>
        </div>
      </div>
    </div>

    <!-- 加入俱乐部 -->
    <transition name="modal">
      <div v-if="showJoin" class="overlay" @click.self="showJoin = false">
        <div class="sheet">
          <div class="sheet-title">加入俱乐部</div>
          <input v-model="joinNo" class="text-input" placeholder="输入俱乐部编号" />
          <button class="btn btn-primary" style="width:100%;margin-top:14px"
            :disabled="joinLoading" @click="doJoin">
            {{ joinLoading ? '加入中...' : '加入' }}
          </button>
        </div>
      </div>
    </transition>

    <!-- 创建俱乐部 -->
    <transition name="modal">
      <div v-if="showCreate" class="overlay" @click.self="showCreate = false">
        <div class="sheet">
          <div class="sheet-title">创建俱乐部</div>
          <input v-model="createName" class="text-input" placeholder="俱乐部名称" maxlength="20" />
          <button class="btn btn-primary" style="width:100%;margin-top:14px"
            :disabled="createLoading" @click="doCreate">
            {{ createLoading ? '创建中...' : '创建' }}
          </button>
        </div>
      </div>
    </transition>
    <TabBar />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getClubInfo, getClubTables, joinClub, createClub } from '@/api'
import TabBar from '@/components/TabBar.vue'

const router = useRouter()
const loading     = ref(false)
const showJoin    = ref(false)
const showCreate  = ref(false)
const joinLoading = ref(false)
const createLoading = ref(false)
const joinNo      = ref('')
const createName  = ref('')
const currentClub = ref(null)
const tables      = ref([])

// 从 localStorage 获取最近加入的 clubId
const myClubId    = ref(Number(localStorage.getItem('my_club_id') || 0))

const roleMap = { 1: '创始人', 2: '管理员', 3: '成员' }
const roleLabel = computed(() => roleMap[currentClub.value?.my_role] || '成员')

async function loadClub(id) {
  loading.value = true
  try {
    const data = await getClubInfo(id)
    currentClub.value = data
    const t = await getClubTables(id)
    tables.value = t?.list || []
  } catch {
    currentClub.value = null
  } finally {
    loading.value = false
  }
}

async function doJoin() {
  if (!joinNo.value.trim()) return
  joinLoading.value = true
  try {
    const data = await joinClub(joinNo.value.trim())
    localStorage.setItem('my_club_id', String(data.club_id))
    myClubId.value = data.club_id
    showJoin.value = false
    joinNo.value = ''
    await loadClub(data.club_id)
  } catch (e) {
    alert(e.message || '加入失败')
  } finally {
    joinLoading.value = false
  }
}

async function doCreate() {
  if (!createName.value.trim()) return
  createLoading.value = true
  try {
    const data = await createClub(createName.value.trim(), '')
    localStorage.setItem('my_club_id', String(data.club_id))
    myClubId.value = data.club_id
    showCreate.value = false
    createName.value = ''
    await loadClub(data.club_id)
  } catch (e) {
    alert(e.message || '创建失败')
  } finally {
    createLoading.value = false
  }
}

onMounted(() => {
  if (myClubId.value) loadClub(myClubId.value)
})
</script>

<style scoped>
.nav {
  height: 48px; padding: 0 16px;
  display: flex; align-items: center; justify-content: space-between;
  background: var(--surface); border-bottom: 1px solid var(--border); flex-shrink: 0;
}
.back-btn { background: none; border: none; cursor: pointer; padding: 4px; display: flex; }
.nav-title { font-size: 17px; font-weight: 700; }
.join-btn {
  background: var(--primary); color: #fff; border: none;
  border-radius: 14px; height: 30px; padding: 0 14px; font-size: 13px; cursor: pointer;
}

.club-card {
  margin: 12px 16px; padding: 16px;
  display: flex; align-items: center; gap: 14px;
}
.club-logo {
  width: 54px; height: 54px; border-radius: 14px; flex-shrink: 0;
  background: linear-gradient(135deg, #0EC4B0, #0aaa98);
  display: flex; align-items: center; justify-content: center;
  font-size: 22px; font-weight: 700; color: #fff;
}
.club-info { flex: 1; }
.club-name { font-size: 17px; font-weight: 700; margin-bottom: 4px; }
.club-meta { display: flex; gap: 10px; flex-wrap: wrap; font-size: 11px; color: var(--text-3); margin-bottom: 4px; }
.club-chips { font-size: 13px; color: var(--primary); font-weight: 600; }

.section-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 6px 16px 8px;
}
.section-title { font-size: 14px; font-weight: 600; color: var(--text-2); }
.create-table-btn {
  background: none; border: none; color: var(--primary);
  font-size: 14px; font-weight: 600; cursor: pointer;
}

.empty-state { text-align: center; padding: 40px 0; }
.empty-hero {
  display: flex; flex-direction: column; align-items: center;
  padding: 60px 24px; gap: 12px;
}
.empty-icon { font-size: 52px; }
.empty-text { font-size: 15px; color: var(--text-3); }
.empty-btns { display: flex; gap: 12px; margin-top: 8px; }

.table-item {
  display: flex; align-items: center; gap: 12px;
  padding: 12px 16px; margin: 0 16px 8px; cursor: pointer;
}
.ti-avatar {
  width: 42px; height: 42px; border-radius: 10px;
  background: linear-gradient(135deg, #0EC4B0, #0aaa98);
  display: flex; align-items: center; justify-content: center;
  font-size: 16px; font-weight: 700; color: #fff; flex-shrink: 0;
}
.ti-info { flex: 1; }
.ti-name { font-size: 14px; font-weight: 600; margin-bottom: 4px; }
.ti-tags { display: flex; gap: 6px; }
.ti-tag {
  background: var(--bg); border-radius: 4px;
  padding: 2px 7px; font-size: 11px; color: var(--text-2);
}

/* Overlay & Sheet */
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
