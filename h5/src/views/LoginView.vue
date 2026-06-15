<template>
  <div class="login-page">
    <div class="login-bg" />
    <div class="login-card">
      <div class="login-logo">🃏</div>
      <h1 class="login-title">好友桌</h1>
      <p class="login-sub">德州扑克 · 和朋友一起玩</p>

      <div class="tabs">
        <span :class="{ active: mode === 'login' }" @click="mode = 'login'">登录</span>
        <span :class="{ active: mode === 'register' }" @click="mode = 'register'">注册</span>
      </div>

      <div class="form">
        <div class="input-group">
          <span class="input-prefix">🇨🇳 +86</span>
          <input v-model="phone" type="tel" maxlength="11" placeholder="手机号" />
        </div>
        <div v-if="mode === 'register'" class="input-group">
          <span class="input-prefix">👤</span>
          <input v-model="nickname" type="text" placeholder="昵称" />
        </div>
        <div class="input-group">
          <span class="input-prefix">🔒</span>
          <input v-model="password" :type="showPwd ? 'text' : 'password'" placeholder="密码（6位以上）" />
          <button class="eye-btn" @click="showPwd = !showPwd">{{ showPwd ? '🙈' : '👁' }}</button>
        </div>

        <div v-if="error" class="error-msg">{{ error }}</div>

        <button class="btn btn-primary submit-btn" :disabled="loading" @click="submit">
          {{ loading ? '请稍候...' : mode === 'login' ? '登录' : '注册并登录' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const mode = ref('login')
const phone = ref('')
const password = ref('')
const nickname = ref('')
const showPwd = ref(false)
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  if (!phone.value || phone.value.length !== 11) { error.value = '请输入正确的手机号'; return }
  if (!password.value || password.value.length < 6) { error.value = '密码至少6位'; return }
  if (mode.value === 'register' && !nickname.value) { error.value = '请输入昵称'; return }

  loading.value = true
  try {
    if (mode.value === 'login') await auth.login(phone.value, password.value)
    else await auth.register(phone.value, password.value, nickname.value)
    router.replace('/home')
  } catch (e) {
    error.value = e.message || '操作失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  position: fixed; inset: 0;
  display: flex; align-items: flex-end; justify-content: center;
  background: linear-gradient(160deg, #0EC4B0 0%, #0aaa98 40%, #1a3d27 100%);
}
.login-bg {
  position: absolute; inset: 0;
  background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M30 5 L55 50 H5 Z' fill='rgba(255,255,255,0.03)'/%3E%3C/svg%3E") repeat;
}
.login-card {
  position: relative; width: 100%; z-index: 1;
  background: var(--surface);
  border-radius: 28px 28px 0 0;
  padding: 32px 24px calc(40px + env(safe-area-inset-bottom));
}
.login-logo { font-size: 52px; text-align: center; margin-bottom: 8px; }
.login-title { font-size: 24px; font-weight: 800; text-align: center; margin-bottom: 4px; }
.login-sub { font-size: 13px; color: var(--text-3); text-align: center; margin-bottom: 24px; }

.tabs {
  display: flex; gap: 24px; margin-bottom: 20px;
}
.tabs span {
  font-size: 16px; font-weight: 500; color: var(--text-3);
  cursor: pointer; padding-bottom: 8px; transition: all .15s;
}
.tabs span.active {
  color: var(--primary); font-weight: 700;
  border-bottom: 2.5px solid var(--primary);
}

.form { display: flex; flex-direction: column; gap: 14px; }
.input-group {
  display: flex; align-items: center; gap: 10px;
  background: var(--bg); border-radius: 12px;
  padding: 0 14px; height: 52px;
  border: 1.5px solid transparent; transition: border-color .15s;
}
.input-group:focus-within { border-color: var(--primary); }
.input-prefix { font-size: 15px; color: var(--text-2); flex-shrink: 0; }
.input-group input {
  flex: 1; border: none; background: transparent;
  font-size: 16px; color: var(--text-1); outline: none;
}
.eye-btn { background: none; border: none; font-size: 16px; cursor: pointer; padding: 4px; }

.error-msg {
  color: var(--red); font-size: 13px; text-align: center;
  padding: 8px; background: rgba(255,82,82,.08);
  border-radius: 8px;
}
.submit-btn {
  width: 100%; height: 52px; font-size: 16px; font-weight: 700;
  border-radius: 14px; margin-top: 4px;
}
.submit-btn:disabled { opacity: .5; }
</style>
