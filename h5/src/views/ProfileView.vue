<template>
  <div class="page">
    <header class="nav">
      <button class="back-btn" @click="router.back()">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M15 18l-6-6 6-6" stroke="#1A1A1A" stroke-width="2.2" stroke-linecap="round"/>
        </svg>
      </button>
      <span class="nav-title">个人资料</span>
      <div style="width:36px" />
    </header>

    <div class="page-body">
      <!-- 头像 + 基本信息 -->
      <div class="profile-hero card">
        <div class="avatar-wrap">
          <div class="avatar-lg">
            <img v-if="auth.user?.avatar" :src="auth.user.avatar" />
            <span v-else>{{ auth.user?.nickname?.[0] || '?' }}</span>
          </div>
        </div>
        <div class="user-name">{{ auth.user?.nickname || '—' }}</div>
        <div class="user-uid">ID: {{ auth.user?.uid || '—' }}</div>

        <div class="chips-row">
          <div class="chip-item">
            <div class="chip-val">{{ formatChips(auth.user?.chips) }}</div>
            <div class="chip-key">筹码</div>
          </div>
        </div>
      </div>

      <!-- 信息列表 -->
      <div class="info-section card">
        <div class="info-row">
          <span class="info-label">手机号</span>
          <span class="info-val">{{ maskedPhone }}</span>
        </div>
        <div class="divider" />
        <div class="info-row" @click="router.push('/career')">
          <span class="info-label">生涯统计</span>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M9 18l6-6-6-6" stroke="#999" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </div>
        <div class="divider" />
        <div class="info-row" @click="router.push('/home')">
          <span class="info-label">我的牌局</span>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M9 18l6-6-6-6" stroke="#999" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </div>
      </div>

      <!-- 退出登录 -->
      <div style="padding: 24px 16px 0">
        <button class="logout-btn" @click="logout">退出登录</button>
      </div>
    </div>
    <TabBar />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import TabBar from '@/components/TabBar.vue'

const router = useRouter()
const auth = useAuthStore()

const maskedPhone = computed(() => {
  const p = auth.user?.phone || ''
  if (p.length < 7) return p
  return p.slice(0, 3) + '****' + p.slice(-4)
})

function formatChips(c) {
  if (!c && c !== 0) return '—'
  return c >= 10000 ? (c / 10000).toFixed(1) + 'w' : c
}

async function logout() {
  if (!confirm('确认退出登录？')) return
  auth.logout()
  router.replace('/login')
}

onMounted(() => auth.fetchProfile())
</script>

<style scoped>
.nav {
  height: 48px; padding: 0 16px;
  display: flex; align-items: center; justify-content: space-between;
  background: var(--surface); border-bottom: 1px solid var(--border); flex-shrink: 0;
}
.back-btn { background: none; border: none; cursor: pointer; padding: 4px; display: flex; }
.nav-title { font-size: 17px; font-weight: 700; }

.profile-hero {
  margin: 16px; padding: 24px 16px;
  display: flex; flex-direction: column; align-items: center; gap: 8px;
}
.avatar-lg {
  width: 80px; height: 80px; border-radius: 50%;
  background: linear-gradient(135deg, #0EC4B0, #0aaa98);
  display: flex; align-items: center; justify-content: center;
  font-size: 32px; font-weight: 700; color: #fff; overflow: hidden;
}
.avatar-lg img { width: 100%; height: 100%; object-fit: cover; }
.user-name { font-size: 20px; font-weight: 700; }
.user-uid { font-size: 13px; color: var(--text-3); }
.chips-row { display: flex; gap: 24px; margin-top: 8px; }
.chip-item { text-align: center; }
.chip-val { font-size: 22px; font-weight: 700; color: var(--primary); }
.chip-key { font-size: 12px; color: var(--text-3); margin-top: 2px; }

.info-section { margin: 0 16px; }
.info-row {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px; cursor: pointer;
}
.info-row:active { background: var(--bg); }
.info-label { font-size: 15px; }
.info-val { font-size: 14px; color: var(--text-3); }

.logout-btn {
  width: 100%; height: 50px; border-radius: 14px;
  border: 1.5px solid var(--red); background: transparent;
  color: var(--red); font-size: 16px; font-weight: 600; cursor: pointer;
}
.logout-btn:active { background: rgba(255,82,82,.08); }
</style>
