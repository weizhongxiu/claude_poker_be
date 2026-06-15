<template>
  <div class="page create-page">
    <!-- 导航栏 -->
    <header class="nav">
      <button class="back-btn" @click="router.back()">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M15 18l-6-6 6-6" stroke="#1A1A1A" stroke-width="2.2" stroke-linecap="round"/>
        </svg>
      </button>
      <span class="nav-title">德州</span>
      <div class="nav-right">
        <span class="nav-label">密码局</span>
        <div class="toggle" :class="{ on: form.hasPassword }" @click="form.hasPassword = !form.hasPassword" />
      </div>
    </header>

    <div class="page-body">
      <!-- 牌桌属性 -->
      <section class="section">
        <div class="section-title">牌桌属性</div>

        <!-- 小盲/大盲 / Straddle / Ante -->
        <div class="field-row">
          <div class="field">
            <div class="field-label">小盲/大盲</div>
            <div class="field-value">{{ form.smallBlind }}/{{ form.bigBlind }}</div>
            <input type="range" :min="1" :max="500" v-model.number="form.smallBlind"
              @input="form.bigBlind = form.smallBlind * 2"
              :style="rangeStyle(form.smallBlind, 1, 500)" />
          </div>
          <div class="field field--center">
            <div class="field-label">Straddle <span class="info-icon">ⓘ</span></div>
            <div class="field-value">{{ form.straddleEnabled ? form.bigBlind * 2 : '—' }}</div>
            <div class="toggle sm" :class="{ on: form.straddleEnabled }"
              @click="form.straddleEnabled = !form.straddleEnabled" />
          </div>
          <div class="field field--right">
            <div class="field-label">Ante</div>
            <div class="select-row" @click="showAnteSheet = true">
              <span class="field-value">{{ form.ante }}</span>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                <path d="M6 9l6 6 6-6" stroke="#999" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
          </div>
        </div>

        <!-- 买入范围 -->
        <div class="field-block">
          <div class="field-row-inline">
            <div>
              <div class="field-label">最小买入</div>
              <div class="field-value">{{ form.minBuyin }}</div>
            </div>
            <div style="text-align:right">
              <div class="field-label">最大买入</div>
              <div class="field-value">{{ form.maxBuyin }}</div>
            </div>
          </div>
          <input type="range" :min="20" :max="50000" :step="20" v-model.number="form.minBuyin"
            @input="onMinBuyinChange"
            :style="rangeStyle(form.minBuyin, 20, 50000)" />
          <div style="margin-top:10px">
            <div class="field-label">最大买入</div>
            <input type="range" :min="form.minBuyin" :max="200000" :step="20" v-model.number="form.maxBuyin"
              :style="rangeStyle(form.maxBuyin, form.minBuyin, 200000)" />
          </div>
        </div>

        <!-- 累计最大买入 -->
        <div class="field-block">
          <div class="field-row-inline">
            <div>
              <div class="field-label">累计最大买入</div>
              <div class="field-value">{{ form.maxBuyinTotal === 0 ? '无限制' : form.maxBuyinTotal }}</div>
            </div>
          </div>
          <input type="range" :min="0" :max="500000" :step="100" v-model.number="form.maxBuyinTotal"
            :style="rangeStyle(form.maxBuyinTotal, 0, 500000)" />
        </div>

        <!-- 时长 -->
        <div class="field-block">
          <div class="field-label">时长 (h)</div>
          <div class="field-value">{{ durationLabel }}</div>
          <input type="range" :min="0" :max="8" :step="1" v-model.number="form.durationIdx"
            :style="rangeStyle(form.durationIdx, 0, 8)" />
          <div class="tick-marks">
            <span v-for="t in durationTicks" :key="t">{{ t }}</span>
          </div>
        </div>

        <!-- 牌桌人数 -->
        <div class="field-row-select" @click="showSeatsSheet = true">
          <span class="field-label">牌桌人数</span>
          <div class="select-row">
            <span class="field-value">{{ form.maxSeats }}</span>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
              <path d="M6 9l6 6 6-6" stroke="#999" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </div>
        </div>

        <!-- 密码（仅密码局） -->
        <div v-if="form.hasPassword" class="field-block">
          <div class="field-label">房间密码</div>
          <input
            v-model="form.password"
            class="text-input"
            placeholder="留空自动生成"
            maxlength="20"
          />
        </div>
      </section>

      <!-- 玩法开关 -->
      <section class="section">
        <div class="toggle-row">
          <span class="field-label">All-In 后支持发两次</span>
          <div class="toggle" :class="{ on: form.runTwice }" @click="form.runTwice = !form.runTwice" />
        </div>
        <div class="divider" />
        <div class="toggle-row">
          <span class="field-label">低水保险 <span class="info-icon">ⓘ</span></span>
          <div class="toggle" :class="{ on: form.lowWater }" @click="form.lowWater = !form.lowWater" />
        </div>
      </section>

      <!-- 高级设置 -->
      <section class="section">
        <div class="section-row" @click="showAdvanced = !showAdvanced">
          <span class="section-title">高级设置</span>
          <span class="expand-btn">{{ showAdvanced ? '收起' : '展开' }}
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
              <path :d="showAdvanced ? 'M18 15l-6-6-6 6' : 'M6 9l6 6 6-6'"
                stroke="#999" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </span>
        </div>
        <transition name="expand">
          <div v-if="showAdvanced" class="advanced-chips">
            <span
              v-for="opt in advancedOpts"
              :key="opt.key"
              class="chip"
              :class="{ active: form.advanced[opt.key] }"
              @click="form.advanced[opt.key] = !form.advanced[opt.key]"
            >{{ opt.label }}</span>
          </div>
        </transition>
      </section>

      <div style="height: 80px;" />
    </div>

    <!-- 底部：桌名 + 开局按钮 -->
    <div class="bottom-bar">
      <input
        v-model="form.name"
        class="name-input"
        :placeholder="defaultName"
      />
      <button class="btn btn-primary start-btn" :disabled="loading" @click="create">
        {{ loading ? '创建中...' : '立即开局' }}
      </button>
    </div>

    <!-- Ante 选择器 -->
    <transition name="modal">
      <div v-if="showAnteSheet" class="overlay" @click.self="showAnteSheet = false">
        <div class="sheet">
          <div class="sheet-title">选择前注</div>
          <div class="sheet-options">
            <div v-for="v in [0, 1, 2, 5]" :key="v" class="sheet-opt"
              :class="{ active: form.ante === v }"
              @click="form.ante = v; showAnteSheet = false">{{ v || '无' }}</div>
          </div>
        </div>
      </div>
    </transition>

    <!-- 人数选择器 -->
    <transition name="modal">
      <div v-if="showSeatsSheet" class="overlay" @click.self="showSeatsSheet = false">
        <div class="sheet">
          <div class="sheet-title">牌桌人数</div>
          <div class="sheet-options">
            <div v-for="v in [2,3,4,5,6,7,8,9,10]" :key="v" class="sheet-opt"
              :class="{ active: form.maxSeats === v }"
              @click="form.maxSeats = v; showSeatsSheet = false">{{ v }}</div>
          </div>
        </div>
      </div>
    </transition>

    <!-- 买入弹窗（创建后入座） -->
    <transition name="modal">
      <div v-if="showBuyinDialog" class="overlay">
        <div class="sheet">
          <div class="sheet-title">选择座位和买入金额</div>
          <div class="buyin-info">
            <div class="buyin-range">买入范围：{{ form.minBuyin }} ~ {{ form.maxBuyin }}</div>
            <div class="room-info" v-if="createdTable">
              <span class="room-label">房间号</span>
              <span class="room-no">{{ createdTable.table_no }}</span>
              <button class="copy-btn-sm" @click="copyTableNo">复制</button>
            </div>
          </div>
          <!-- 座位选择 -->
          <div class="seat-grid">
            <div
              v-for="n in form.maxSeats"
              :key="n"
              class="seat-btn"
              :class="{ active: selectedSeat === n }"
              @click="selectedSeat = n"
            >{{ n }}</div>
          </div>
          <div class="field-label" style="margin:12px 0 6px">买入金额</div>
          <input type="range" :min="form.minBuyin" :max="form.maxBuyin" :step="form.minBuyin"
            v-model.number="buyinAmount"
            :style="rangeStyle(buyinAmount, form.minBuyin, form.maxBuyin)" />
          <div class="buyin-display">{{ buyinAmount }}</div>
          <button class="btn btn-primary" style="width:100%;margin-top:16px"
            :disabled="!selectedSeat || joiningLoading"
            @click="joinAndStart">
            {{ joiningLoading ? '加入中...' : `入座第 ${selectedSeat} 座` }}
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { createTable, takeSeat } from '@/api'

const router = useRouter()
const auth = useAuthStore()

const loading = ref(false)
const joiningLoading = ref(false)
const showAdvanced = ref(false)
const showAnteSheet = ref(false)
const showSeatsSheet = ref(false)
const showBuyinDialog = ref(false)
const createdTable = ref(null)
const selectedSeat = ref(1)
const buyinAmount = ref(200)

const durationOptions = [0.5, 1, 1.5, 2, 2.5, 3, 4, 6, 8]
const durationTicks = ['0.5', '1', '1.5', '2', '2.5', '3', '4', '6', '8']

const form = ref({
  name: '',
  hasPassword: true,
  password: '',
  smallBlind: 1,
  bigBlind: 2,
  straddleEnabled: false,
  ante: 0,
  minBuyin: 200,
  maxBuyin: 2000,
  maxBuyinTotal: 0,
  durationIdx: 1,
  maxSeats: 8,
  runTwice: false,
  lowWater: false,
  advanced: {
    crit: false, activity: false, autoRebuy: false,
    buyinApproval: false, delayShow: false, randomSeat: false,
    specMute: false, gpsIp: false, fullStart: false,
  },
})

const advancedOpts = [
  { key: 'crit',          label: '暴击玩法' },
  { key: 'activity',      label: '活跃度积分' },
  { key: 'autoRebuy',     label: '自动补码/退码' },
  { key: 'buyinApproval', label: '带入审核' },
  { key: 'delayShow',     label: '延迟看牌' },
  { key: 'randomSeat',    label: '随机入座' },
  { key: 'specMute',      label: '旁观者禁言' },
  { key: 'gpsIp',         label: 'GPS和IP限制' },
  { key: 'fullStart',     label: '人满开局' },
]

const durationLabel = computed(() => {
  const v = durationOptions[form.value.durationIdx]
  return v < 1 ? `${v * 60}分钟` : `${v}小时`
})

const defaultName = computed(() => `${auth.user?.nickname || ''}的牌局`)

function onMinBuyinChange() {
  if (form.value.maxBuyin < form.value.minBuyin * 2) {
    form.value.maxBuyin = form.value.minBuyin * 2
  }
}

function rangeStyle(val, min, max) {
  const pct = max === min ? 0 : ((val - min) / (max - min)) * 100
  return `background: linear-gradient(to right, #0EC4B0 ${pct}%, #EBEBEB ${pct}%)`
}

function copyTableNo() {
  navigator.clipboard?.writeText(createdTable.value?.table_no || '')
  alert('房间号已复制')
}

async function create() {
  if (loading.value) return
  loading.value = true
  try {
    const adv = form.value.advanced
    // 自动生成密码
    const pwd = form.value.hasPassword
      ? (form.value.password || Math.random().toString(36).slice(2, 8).toUpperCase())
      : ''

    const data = await createTable({
      name:               form.value.name || defaultName.value,
      game_type:          1,
      has_password:       form.value.hasPassword ? 1 : 0,
      password:           pwd,
      small_blind:        form.value.smallBlind,
      big_blind:          form.value.bigBlind,
      ante:               form.value.ante,
      straddle_enabled:   form.value.straddleEnabled ? 1 : 0,
      min_buyin:          form.value.minBuyin,
      max_buyin:          form.value.maxBuyin,
      max_buyin_total:    form.value.maxBuyinTotal,
      duration:           durationOptions[form.value.durationIdx],
      max_seats:          form.value.maxSeats,
      run_twice:          form.value.runTwice ? 1 : 0,
      low_water_insurance: form.value.lowWater ? 1 : 0,
      crit_gameplay:      adv.crit ? 1 : 0,
      activity_points:    adv.activity ? 1 : 0,
      auto_rebuy:         adv.autoRebuy ? 1 : 0,
      buyin_approval:     adv.buyinApproval ? 1 : 0,
      delay_show_card:    adv.delayShow ? 1 : 0,
      random_seat:        adv.randomSeat ? 1 : 0,
      spectator_mute:     adv.specMute ? 1 : 0,
      gps_ip_restrict:    adv.gpsIp ? 1 : 0,
      full_table_start:   adv.fullStart ? 1 : 0,
    })
    createdTable.value = { table_id: data.table_id, table_no: data.table_no }
    buyinAmount.value = form.value.minBuyin
    selectedSeat.value = 1
    showBuyinDialog.value = true
  } catch (e) {
    alert(e.message || '创建失败')
  } finally {
    loading.value = false
  }
}

async function joinAndStart() {
  if (!selectedSeat.value || joiningLoading.value) return
  joiningLoading.value = true
  try {
    await takeSeat(createdTable.value.table_id, selectedSeat.value, buyinAmount.value)
    showBuyinDialog.value = false
    // 直接进入牌桌，由 GameTableView 处理 startSession
    router.push(`/table/${createdTable.value.table_id}?creator=1`)
  } catch (e) {
    alert(e.message || '入座失败')
  } finally {
    joiningLoading.value = false
  }
}
</script>

<style scoped>
.create-page { background: var(--bg); }
.nav {
  height: 48px; padding: 0 16px;
  display: flex; align-items: center; justify-content: space-between;
  background: var(--surface); border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.back-btn { background: none; border: none; cursor: pointer; padding: 4px; display: flex; }
.nav-title { font-size: 17px; font-weight: 700; }
.nav-right { display: flex; align-items: center; gap: 8px; }
.nav-label { font-size: 14px; color: var(--text-2); }

.section {
  background: var(--surface); margin: 10px 0; padding: 14px 16px;
}
.section-title { font-size: 13px; font-weight: 600; color: var(--text-3); margin-bottom: 14px; }
.section-row {
  display: flex; align-items: center; justify-content: space-between; cursor: pointer;
}
.expand-btn { font-size: 12px; color: var(--text-3); display: flex; align-items: center; gap: 4px; }

.field-row {
  display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16px; margin-bottom: 20px;
}
.field-label { font-size: 12px; color: var(--text-3); margin-bottom: 6px; }
.info-icon { font-size: 12px; color: var(--text-3); }
.field-value { font-size: 22px; font-weight: 700; color: var(--text-1); margin-bottom: 8px; }
.select-row { display: flex; align-items: center; gap: 4px; cursor: pointer; }
.field--center { display: flex; flex-direction: column; align-items: center; }
.field--right { display: flex; flex-direction: column; align-items: flex-end; }

.field-block { margin-bottom: 20px; }
.field-row-inline { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 8px; }

.tick-marks { display: flex; justify-content: space-between; margin-top: 8px; }
.tick-marks span { font-size: 10px; color: var(--text-3); }

.field-row-select {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 0; border-top: 1px solid var(--border); cursor: pointer;
}

.text-input {
  width: 100%; height: 44px; border-radius: 10px;
  border: 1.5px solid var(--border); padding: 0 14px;
  font-size: 15px; outline: none; background: var(--bg); color: var(--text-1);
  margin-top: 4px;
}
.text-input:focus { border-color: var(--primary); }

.toggle-row {
  display: flex; align-items: center; justify-content: space-between; padding: 12px 0;
}
.toggle.sm { width: 36px; height: 22px; border-radius: 11px; }
.toggle.sm::after { width: 16px; height: 16px; top: 3px; left: 3px; }
.toggle.sm.on::after { transform: translateX(14px); }

.advanced-chips { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 12px; }

.bottom-bar {
  position: fixed; bottom: 0; left: 0; right: 0;
  display: flex; align-items: center; gap: 12px;
  padding: 10px 16px calc(10px + env(safe-area-inset-bottom));
  background: var(--surface); border-top: 1px solid var(--border);
}
.name-input {
  flex: 1; height: 44px; border-radius: 22px;
  border: 1px solid var(--border); padding: 0 16px;
  font-size: 14px; outline: none; background: var(--bg); color: var(--text-1);
}
.name-input:focus { border-color: var(--primary); }
.start-btn { height: 44px; padding: 0 24px; white-space: nowrap; }
.start-btn:disabled { opacity: .5; }

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
.sheet-title { font-size: 16px; font-weight: 600; margin-bottom: 16px; text-align: center; }
.sheet-options { display: flex; flex-wrap: wrap; gap: 10px; justify-content: center; }
.sheet-opt {
  min-width: 60px; height: 40px; border-radius: 20px;
  border: 1.5px solid var(--border);
  display: flex; align-items: center; justify-content: center;
  font-size: 15px; font-weight: 500; cursor: pointer; padding: 0 16px;
}
.sheet-opt.active { border-color: var(--primary); color: var(--primary); background: rgba(14,196,176,.08); }

/* Buyin dialog extras */
.buyin-info { background: var(--bg); border-radius: 10px; padding: 10px 14px; margin-bottom: 14px; }
.buyin-range { font-size: 13px; color: var(--text-2); margin-bottom: 6px; }
.room-info { display: flex; align-items: center; gap: 8px; }
.room-label { font-size: 12px; color: var(--text-3); }
.room-no { font-size: 18px; font-weight: 700; flex: 1; }
.copy-btn-sm {
  background: var(--primary); color: #fff; border: none;
  border-radius: 12px; height: 26px; padding: 0 12px;
  font-size: 12px; cursor: pointer;
}
.seat-grid { display: flex; flex-wrap: wrap; gap: 10px; margin-bottom: 6px; }
.seat-btn {
  width: 44px; height: 44px; border-radius: 10px;
  border: 1.5px solid var(--border); background: var(--bg);
  display: flex; align-items: center; justify-content: center;
  font-size: 15px; font-weight: 600; cursor: pointer;
}
.seat-btn.active { border-color: var(--primary); color: var(--primary); background: rgba(14,196,176,.1); }
.buyin-display { text-align: center; font-size: 26px; font-weight: 700; margin-top: 8px; }

/* Animations */
.expand-enter-active, .expand-leave-active { transition: all .2s ease; overflow: hidden; }
.expand-enter-from, .expand-leave-to { opacity: 0; max-height: 0; }
.expand-enter-to, .expand-leave-from { opacity: 1; max-height: 200px; }
.modal-enter-active { animation: slideUp .3s ease; }
.modal-leave-active { animation: slideUp .25s ease reverse; }
</style>
