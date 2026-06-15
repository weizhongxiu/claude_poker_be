<template>
  <nav class="tabbar">
    <div
      v-for="item in tabs"
      :key="item.key"
      class="tab-item"
      :class="{ active: current === item.key }"
      @click="go(item)"
    >
      <div class="tab-icon">
        <component :is="item.icon" :active="current === item.key" />
      </div>
      <span class="tab-label">{{ item.label }}</span>
    </div>
  </nav>
</template>

<script setup>
import { computed, defineComponent, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const current = computed(() => {
  if (route.path.startsWith('/career')) return 'career'
  if (route.path.startsWith('/club')) return 'club'
  if (route.path.startsWith('/discover')) return 'discover'
  if (route.path.startsWith('/profile')) return 'profile'
  return 'home'
})

const makeIcon = (paths, vboxSize = '0 0 24 24') =>
  defineComponent({
    props: ['active'],
    render() {
      return h('svg', { viewBox: vboxSize, fill: 'none', width: 22, height: 22 },
        paths.map(d => h('path', { d, fill: this.active ? '#fff' : '#BBBBBB', fillRule: 'evenodd' }))
      )
    }
  })

const IconClub = makeIcon([
  'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 14H9V8h2v8zm4 0h-2V8h2v8z'
])
const IconDiscover = makeIcon([
  'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z'
])
const IconCareer = makeIcon([
  'M3 3h18v2H3V3zm0 4h12v2H3V7zm0 4h18v2H3v-2zm0 4h12v2H3v-2z'
])
const IconProfile = makeIcon([
  'M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z'
])

// Home icon (cards icon) — center highlighted
const IconHome = defineComponent({
  props: ['active'],
  render() {
    const c = this.active ? '#fff' : '#BBBBBB'
    return h('svg', { viewBox: '0 0 24 24', fill: 'none', width: 26, height: 26 }, [
      h('rect', { x: 3, y: 5, width: 18, height: 14, rx: 3, fill: c }),
      h('rect', { x: 6, y: 8, width: 5, height: 8, rx: 1, fill: this.active ? '#0EC4B0' : '#f0f0f0', opacity: .8 }),
      h('rect', { x: 13, y: 8, width: 5, height: 8, rx: 1, fill: this.active ? '#0EC4B0' : '#f0f0f0', opacity: .8 }),
    ])
  }
})

const tabs = [
  { key: 'club',     label: '俱乐部', path: '/club',     icon: IconClub },
  { key: 'discover', label: '发现',   path: '/discover', icon: IconDiscover },
  { key: 'home',     label: '好友桌', path: '/home',     icon: IconHome },
  { key: 'career',   label: '生涯',   path: '/career',   icon: IconCareer },
  { key: 'profile',  label: '我的',   path: '/profile',  icon: IconProfile },
]

function go(item) {
  if (item.key === 'home') router.push('/home')
  else if (item.key === 'career') router.push('/career')
  else router.push(item.path)
}
</script>

<style scoped>
.tabbar {
  position: fixed; bottom: 0; left: 0; right: 0; z-index: 100;
  height: var(--tabbar-h);
  padding-bottom: env(safe-area-inset-bottom);
  display: flex; align-items: flex-start;
  background: rgba(255,255,255,.96);
  backdrop-filter: blur(10px);
  border-top: 1px solid var(--border);
}
.tab-item {
  flex: 1; display: flex; flex-direction: column; align-items: center;
  justify-content: center; padding-top: 6px; cursor: pointer;
  gap: 2px; transition: opacity .15s;
}
.tab-item:active { opacity: .6; }
.tab-label {
  font-size: 10px; color: var(--text-3); transition: color .15s;
}

/* All active tabs: green circle icon + green label */
.tab-item.active .tab-label { color: var(--primary); font-weight: 600; }
.tab-item.active .tab-icon {
  width: 54px; height: 54px; border-radius: 27px;
  background: var(--primary);
  display: flex; align-items: center; justify-content: center;
  margin-top: -24px;
  box-shadow: 0 4px 16px rgba(14,196,176,.5);
}

.tab-item:nth-child(3) { position: relative; }
</style>
