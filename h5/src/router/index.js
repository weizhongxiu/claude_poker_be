import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  { path: '/', redirect: '/home' },

  // 主页
  {
    path: '/home',
    component: () => import('@/views/HomeView.vue'),
    meta: { tab: 'home' }
  },

  // 创建牌桌
  {
    path: '/create',
    component: () => import('@/views/CreateTableView.vue'),
    meta: { requiresAuth: true }
  },

  // 牌桌列表（大厅）
  {
    path: '/tables',
    component: () => import('@/views/TableListView.vue'),
    meta: { tab: 'home' }
  },

  // 游戏桌（横屏）
  {
    path: '/table/:id',
    component: () => import('@/views/GameTableView.vue'),
    meta: { requiresAuth: true, landscape: true }
  },

  // 生涯统计
  {
    path: '/career',
    component: () => import('@/views/CareerView.vue'),
    meta: { tab: 'career', requiresAuth: true }
  },

  // 牌局结算详情
  {
    path: '/stats/session/:id',
    component: () => import('@/views/SessionDetailView.vue'),
    meta: { requiresAuth: true }
  },

  // 个人资料
  {
    path: '/profile',
    component: () => import('@/views/ProfileView.vue'),
    meta: { requiresAuth: true }
  },

  // 俱乐部
  {
    path: '/club',
    component: () => import('@/views/ClubView.vue'),
    meta: { requiresAuth: true }
  },
  // 兼容旧跳转路径
  {
    path: '/club/create',
    redirect: '/club'
  },
  {
    path: '/club/join',
    redirect: '/club'
  },

  // 登录/注册
  {
    path: '/login',
    component: () => import('@/views/LoginView.vue')
  },

  // 发现页（暂时重定向到大厅）
  { path: '/discover', redirect: '/tables' },

  // 手牌回放
  {
    path: '/replay/:id',
    component: () => import('@/views/HandReplayView.vue'),
    meta: { requiresAuth: true }
  },

  // 404 兜底
  { path: '/:pathMatch(.*)*', redirect: '/home' }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 })
})

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth) {
    const auth = useAuthStore()
    if (!auth.isLoggedIn) return next('/login')
  }
  next()
})

export default router
