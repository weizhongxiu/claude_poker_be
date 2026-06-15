<template>
  <div class="replay-wrap">
    <!-- 顶部栏 -->
    <div class="top-bar">
      <button class="icon-btn" @click="router.back()">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M15 18l-6-6 6-6" stroke="#fff" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </button>
      <div class="title">手牌 #{{ handNo }}</div>
      <div class="top-right">
        <button class="icon-btn" @click="toggleFav">
          <svg width="20" height="20" viewBox="0 0 24 24" :fill="isFav ? '#FFD700' : 'none'">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
              stroke="#FFD700" stroke-width="1.5"/>
          </svg>
        </button>
        <button class="speed-btn" @click="cycleSpeed">X{{ speed }}</button>
      </div>
    </div>

    <!-- 加载中 / 空状态 -->
    <div v-if="loading" class="center-state">加载中...</div>
    <div v-else-if="!replayData" class="center-state">暂无数据</div>

    <!-- 牌桌主体 -->
    <template v-else>
      <div class="table-container">
        <div class="table-felt">
          <div class="table-inner">
            <!-- 公共牌 -->
            <div class="community-cards">
              <div v-for="(card, i) in currentCommunityCards" :key="i"
                class="playing-card" :class="cardColor(card)">
                <span class="card-rank">{{ cardRank(card) }}</span>
                <span class="card-suit">{{ cardSuit(card) }}</span>
              </div>
              <div v-for="i in Math.max(0, 5 - currentCommunityCards.length)"
                :key="'ph' + i" class="playing-card placeholder" />
            </div>
            <!-- 底池 -->
            <div class="pot-area" v-if="currentPot > 0">
              <span class="pot-label">底池</span>
              <span class="pot-amount">{{ currentPot }}</span>
            </div>
          </div>
        </div>

        <!-- 座位 -->
        <div v-for="seat in seatPositions" :key="seat.no"
          class="seat-wrapper" :style="seat.style">
          <div v-if="playersMap[seat.no]" class="seat occupied"
            :class="{
              'is-active': currentAction && currentAction.seat_no === seat.no,
              'is-folded': foldedSeats.has(seat.no)
            }">
            <div v-if="replayData.dealer_seat === seat.no" class="dealer-badge">D</div>
            <div class="seat-avatar">
              <span>{{ playersMap[seat.no].nickname?.[0] || '?' }}</span>
            </div>
            <div class="seat-info">
              <div class="seat-name">{{ playersMap[seat.no].nickname }}</div>
              <div class="seat-chips">{{ currentChips[seat.no] ?? playersMap[seat.no].chips }}</div>
            </div>
            <div v-if="currentBets[seat.no]" class="seat-bet">{{ currentBets[seat.no] }}</div>
            <div v-if="lastActionText[seat.no]" class="action-text">{{ lastActionText[seat.no] }}</div>
            <div v-if="holeCardsMap[seat.no]" class="hole-cards">
              <div v-for="(c, ci) in holeCardsMap[seat.no]" :key="ci"
                class="playing-card sm" :class="cardColor(c)">
                <span class="card-rank">{{ cardRank(c) }}</span>
                <span class="card-suit">{{ cardSuit(c) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部控制区 -->
      <div class="bottom-area">
        <div class="stage-tabs">
          <div v-for="tab in stageTabs" :key="tab.stage"
            class="stage-tab"
            :class="{ active: currentStageIndex === tab.stageIdx }"
            @click="jumpToStage(tab.stageIdx)">
            {{ tab.label }}
          </div>
        </div>
        <div class="control-bar">
          <button class="ctrl-btn" @click="restartReplay">⟳</button>
          <button class="ctrl-btn" @click="goFirst">⏮</button>
          <button class="ctrl-btn" @click="goPrev">◀</button>
          <span class="step-count">{{ currentStep }}/{{ totalSteps }}</span>
          <button class="ctrl-btn" @click="goNext">▶</button>
          <button class="ctrl-btn" @click="goLast">⏭</button>
          <button class="ctrl-btn fav" @click="toggleFav">
            <svg width="16" height="16" viewBox="0 0 24 24" :fill="isFav ? '#FFD700' : 'none'">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
                stroke="#FFD700" stroke-width="1.5"/>
            </svg>
          </button>
        </div>
        <div class="hand-id">ID: {{ gameId }}</div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getHandReplay, toggleFavorite } from '@/api'

const route = useRoute()
const router = useRouter()

const gameId = computed(() => route.params.id)
const loading = ref(true)
const replayData = ref(null)
const isFav = ref(false)
const speed = ref(1)
const currentStep = ref(0)
let autoTimer = null

async function loadReplay() {
  loading.value = true
  try {
    replayData.value = await getHandReplay(gameId.value)
  } catch {
    replayData.value = null
  } finally {
    loading.value = false
  }
}

const handNo = computed(() => replayData.value?.hand_no || gameId.value)

const playersMap = computed(() => {
  if (!replayData.value) return {}
  const map = {}
  for (const stage of (replayData.value.stages || [])) {
    for (const ps of (stage.players_state || [])) {
      if (!map[ps.seat_no]) {
        map[ps.seat_no] = {
          seat_no: ps.seat_no,
          user_id: ps.user_id,
          nickname: ps.nickname || ('玩家' + ps.seat_no),
          chips: ps.chips_start ?? ps.chips ?? 0,
        }
      }
    }
  }
  return map
})

const allActions = computed(() => {
  if (!replayData.value) return []
  const result = []
  const stages = replayData.value.stages || []
  for (let si = 0; si < stages.length; si++) {
    for (const action of (stages[si].actions || [])) {
      result.push({ ...action, stageIdx: si })
    }
  }
  return result
})

const totalSteps = computed(() => allActions.value.length)

const STAGE_LABELS = { 1: '翻前', 2: '翻牌', 3: '转牌', 4: '河牌', 5: '摊牌' }

const stageTabs = computed(() => {
  if (!replayData.value) return []
  return (replayData.value.stages || []).map((s, idx) => ({
    stage: s.stage,
    stageIdx: idx,
    label: STAGE_LABELS[s.stage] || ('阶段' + s.stage),
  }))
})

const currentAction = computed(() => {
  if (currentStep.value === 0) return null
  return allActions.value[currentStep.value - 1] || null
})

const currentStageIndex = computed(() => currentAction.value?.stageIdx ?? 0)

const ACTION_TEXT = { 1: '弃牌', 2: '跟注', 3: '加注', 4: 'All-in', 5: '过牌' }

const replayState = computed(() => {
  const bets = {}
  const chips = {}
  const lastAct = {}
  const folded = new Set()
  const holeCards = {}
  if (!replayData.value) return { bets, chips, lastAct, folded, holeCards }

  const stages = replayData.value.stages || []
  let stepCount = 0

  for (let si = 0; si < stages.length; si++) {
    const stage = stages[si]
    const stageBets = {}

    if (si === 0) {
      for (const ps of (stage.players_state || [])) {
        chips[ps.seat_no] = ps.chips_start ?? ps.chips ?? 0
      }
    }
    for (const ps of (stage.players_state || [])) {
      if (ps.hole_cards) holeCards[ps.seat_no] = ps.hole_cards.split(',')
    }

    for (const action of (stage.actions || [])) {
      if (stepCount >= currentStep.value) break
      stepCount++
      const seat = action.seat_no
      if (action.action === 1) {
        folded.add(seat)
        lastAct[seat] = ACTION_TEXT[1]
        delete stageBets[seat]
      } else {
        lastAct[seat] = ACTION_TEXT[action.action] || String(action.action)
        if (action.amount > 0) {
          stageBets[seat] = (stageBets[seat] || 0) + action.amount
          chips[seat] = (chips[seat] || 0) - action.amount
        }
      }
    }

    for (const [k, v] of Object.entries(stageBets)) bets[k] = v
    if (stepCount >= currentStep.value) break
    if (si < stages.length - 1) {
      for (const k of Object.keys(bets)) delete bets[k]
    }
  }

  return { bets, chips, lastAct, folded, holeCards }
})

const currentBets      = computed(() => replayState.value.bets)
const currentChips     = computed(() => replayState.value.chips)
const lastActionText   = computed(() => replayState.value.lastAct)
const foldedSeats      = computed(() => replayState.value.folded)
const holeCardsMap     = computed(() => replayState.value.holeCards)

const currentPot = computed(() => {
  if (currentStep.value === 0) return 0
  return currentAction.value?.pot_after || 0
})

const currentCommunityCards = computed(() => {
  if (!replayData.value) return []
  const stages = replayData.value.stages || []
  const si = currentStageIndex.value
  const cc = stages[si]?.community_cards || ''
  return cc ? cc.split(',').filter(Boolean) : []
})

const seatPositions = computed(() => {
  const positions = [
    { no: 1, left: '50%', top: '82%' },
    { no: 2, left: '22%', top: '70%' },
    { no: 3, left: '5%',  top: '45%' },
    { no: 4, left: '22%', top: '18%' },
    { no: 5, left: '78%', top: '18%' },
    { no: 6, left: '95%', top: '45%' },
    { no: 7, left: '78%', top: '70%' },
  ]
  return positions.map(p => ({
    ...p,
    style: { position: 'absolute', left: p.left, top: p.top, transform: 'translate(-50%,-50%)' }
  }))
})

function goFirst()   { currentStep.value = 0 }
function goLast()    { currentStep.value = totalSteps.value }
function goNext()    { if (currentStep.value < totalSteps.value) currentStep.value++ }
function goPrev()    { if (currentStep.value > 0) currentStep.value-- }
function restartReplay() { stopAuto(); currentStep.value = 0; startAuto() }

function jumpToStage(stageIdx) {
  const idx = allActions.value.findIndex(a => a.stageIdx === stageIdx)
  currentStep.value = idx >= 0 ? idx : 0
}

function cycleSpeed() {
  const speeds = [1, 2, 4]
  const i = speeds.indexOf(speed.value)
  speed.value = speeds[(i + 1) % speeds.length]
}

function startAuto() {
  stopAuto()
  if (currentStep.value >= totalSteps.value) return
  autoTimer = setInterval(() => {
    if (currentStep.value < totalSteps.value) {
      currentStep.value++
    } else {
      stopAuto()
    }
  }, 1000 / speed.value)
}

function stopAuto() {
  if (autoTimer) { clearInterval(autoTimer); autoTimer = null }
}

watch(speed, () => { if (autoTimer) startAuto() })

async function toggleFav() {
  try {
    await toggleFavorite(gameId.value)
    isFav.value = !isFav.value
  } catch {}
}

const rankMap = { A:'A',K:'K',Q:'Q',J:'J',T:'10' }
const suitMap = { h:'♥',d:'♦',c:'♣',s:'♠' }
function cardRank(c) { return rankMap[c?.[0]] || c?.[0] || '' }
function cardSuit(c) { return suitMap[c?.[c.length-1]] || '' }
function cardColor(c) {
  const s = c?.[c.length - 1]
  return (s === 'h' || s === 'd') ? 'red' : 'black'
}

onMounted(async () => {
  await loadReplay()
  startAuto()
})
onUnmounted(() => stopAuto())
</script>

<style scoped>
.replay-wrap {
  position: fixed; inset: 0;
  background: linear-gradient(180deg, #0a1628 0%, #0d2644 30%, #0e3355 60%, #1a4a6b 100%);
  display: flex; flex-direction: column; overflow: hidden;
}

.top-bar {
  position: relative; z-index: 50;
  height: 48px; display: flex; align-items: center;
  justify-content: space-between; padding: 0 12px;
  background: rgba(0,0,0,.3); flex-shrink: 0;
}
.title { font-size: 15px; font-weight: 700; color: #fff; }
.top-right { display: flex; align-items: center; gap: 8px; }
.icon-btn {
  width: 36px; height: 36px; border-radius: 18px;
  background: rgba(255,255,255,.15); border: none;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; color: #fff;
}
.icon-btn:active { opacity: .7; }
.speed-btn {
  background: rgba(255,255,255,.2); border: none; color: #fff;
  border-radius: 14px; height: 28px; padding: 0 12px;
  font-size: 13px; font-weight: 700; cursor: pointer;
}

.center-state {
  flex: 1; display: flex; align-items: center; justify-content: center;
  color: rgba(255,255,255,.5); font-size: 15px;
}

.table-container { flex: 1; position: relative; overflow: hidden; }

.table-felt {
  position: absolute; left: 50%; top: 50%;
  transform: translate(-50%, -50%);
  width: 60%; height: 45%; border-radius: 50%;
  background: radial-gradient(ellipse at center, #1e6b8a 0%, #155070 65%, #0d3d55 100%);
  border: 8px solid #0a2d3d;
  box-shadow: 0 0 0 3px #0a2d3d, 0 0 40px rgba(0,0,0,.6), inset 0 2px 8px rgba(255,255,255,.05);
}
.table-inner {
  position: absolute; inset: 0;
  display: flex; flex-direction: column;
  align-items: center; justify-content: center; gap: 8px;
}

.community-cards { display: flex; gap: 4px; }
.playing-card {
  width: 28px; height: 40px; border-radius: 4px; background: #fff;
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  box-shadow: 0 2px 8px rgba(0,0,0,.4); font-weight: 700; font-size: 9px; line-height: 1.1;
}
.playing-card.red { color: #e53935; }
.playing-card.black { color: #1a1a1a; }
.playing-card.placeholder {
  background: rgba(255,255,255,.12);
  border: 1px dashed rgba(255,255,255,.3);
}
.card-suit { font-size: 11px; }
.playing-card.sm { width: 22px; height: 32px; font-size: 8px; }

.pot-area {
  display: flex; align-items: center; gap: 8px;
  background: rgba(0,0,0,.35); border-radius: 16px; padding: 3px 12px;
}
.pot-label { font-size: 10px; color: rgba(255,255,255,.6); }
.pot-amount { font-size: 13px; font-weight: 700; color: #fff; }

.seat-wrapper { z-index: 10; }
.seat.occupied {
  display: flex; flex-direction: column; align-items: center; gap: 2px;
  position: relative;
}
.seat.is-active .seat-avatar { border-color: #0EC4B0 !important; box-shadow: 0 0 12px rgba(14,196,176,.6); }
.seat.is-folded { opacity: .35; }

.dealer-badge {
  position: absolute; top: -6px; right: -6px; z-index: 5;
  width: 18px; height: 18px; border-radius: 50%;
  background: #fff; color: #1a1a1a; font-size: 9px; font-weight: 900;
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 2px 6px rgba(0,0,0,.4);
}

.seat-avatar {
  width: 44px; height: 44px; border-radius: 50%;
  background: linear-gradient(135deg, #1a4a6b, #0d3d55);
  border: 2px solid rgba(255,255,255,.3);
  display: flex; align-items: center; justify-content: center;
  font-size: 15px; font-weight: 700; color: #fff; overflow: hidden;
}
.seat-info { text-align: center; }
.seat-name {
  font-size: 10px; color: #fff; white-space: nowrap;
  max-width: 60px; overflow: hidden; text-overflow: ellipsis;
}
.seat-chips {
  background: rgba(0,0,0,.5); border-radius: 10px;
  padding: 1px 6px; font-size: 9px; color: #FFD700; font-weight: 600;
  display: inline-block;
}
.seat-bet {
  position: absolute; top: -18px;
  background: #FFD700; color: #1a1a1a;
  border-radius: 10px; padding: 1px 6px; font-size: 9px; font-weight: 700;
  white-space: nowrap;
}
.action-text {
  font-size: 9px; color: rgba(255,255,255,.7);
  background: rgba(0,0,0,.4); border-radius: 8px; padding: 1px 6px;
}
.hole-cards { display: flex; gap: 2px; margin-top: 2px; }

.bottom-area {
  flex-shrink: 0; padding: 6px 12px 16px;
  background: rgba(0,0,0,.5);
  display: flex; flex-direction: column; gap: 6px;
}

.stage-tabs { display: flex; gap: 4px; justify-content: center; }
.stage-tab {
  padding: 5px 14px; border-radius: 14px;
  font-size: 13px; color: rgba(255,255,255,.6);
  cursor: pointer; transition: all .15s;
  border: 1px solid rgba(255,255,255,.15);
}
.stage-tab.active {
  background: rgba(14,196,176,.8); color: #fff;
  border-color: #0EC4B0; font-weight: 600;
}

.control-bar {
  display: flex; align-items: center; justify-content: center; gap: 8px;
}
.ctrl-btn {
  width: 36px; height: 36px; border-radius: 18px;
  background: rgba(255,255,255,.15); border: none; color: #fff;
  font-size: 14px; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
}
.ctrl-btn:active { opacity: .7; }
.ctrl-btn.fav { background: rgba(255,215,0,.15); }
.step-count { font-size: 14px; font-weight: 700; color: #fff; min-width: 52px; text-align: center; }

.hand-id { text-align: center; font-size: 11px; color: rgba(255,255,255,.35); }
</style>
