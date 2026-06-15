import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, register as apiRegister, getProfile } from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(phone, password) {
    const data = await apiLogin(phone, password)
    token.value = data.token
    user.value = { id: data.user_id, nickname: data.nickname, avatar: data.avatar }
    localStorage.setItem('token', data.token)
    return data
  }

  async function register(phone, password, nickname) {
    const data = await apiRegister(phone, password, nickname)
    token.value = data.token
    user.value = { id: data.user_id, nickname, avatar: '' }
    localStorage.setItem('token', data.token)
    return data
  }

  async function fetchProfile() {
    if (!token.value) return
    const data = await getProfile()
    user.value = {
      id: data.user_id, uid: data.uid,
      nickname: data.nickname, avatar: data.avatar,
      phone: data.phone, chips: data.chips
    }
    return data
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  return { token, user, isLoggedIn, login, register, fetchProfile, logout }
})
