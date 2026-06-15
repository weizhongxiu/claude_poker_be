<template>
  <teleport to="body">
    <transition name="modal">
      <div v-if="visible" class="overlay" @click.self="close">
        <div class="modal">
          <div class="modal-header">
            <span class="modal-title">加入牌局</span>
            <button class="close-btn" @click="close">✕</button>
          </div>

          <div class="modal-body">
            <img src="/chip-pile.svg" class="chip-img" alt="" />

            <!-- 输入框 -->
            <div class="input-wrap">
              <input
                ref="inputEl"
                v-model="code"
                class="code-input"
                type="text"
                maxlength="20"
                placeholder="输入房间号"
                inputmode="numeric"
                @keyup.enter="submit"
              />
            </div>

            <!-- 数字键盘 -->
            <div class="keypad">
              <button
                v-for="key in keys"
                :key="key"
                class="key-btn"
                @click="tap(key)"
              >
                <span v-if="key === '⌫'">⌫</span>
                <span v-else>{{ key }}</span>
              </button>
            </div>

            <div v-if="errMsg" class="err-msg">{{ errMsg }}</div>

            <button
              class="btn btn-primary submit-btn"
              :disabled="!code || submitting"
              @click="submit"
            >{{ submitting ? '加入中...' : '加入牌局' }}</button>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { joinTable } from '@/api'

const props = defineProps({ visible: Boolean })
const emit  = defineEmits(['update:visible', 'joined'])
const router = useRouter()

const code      = ref('')
const submitting = ref(false)
const errMsg    = ref('')
const inputEl   = ref(null)

const keys = ['1','2','3','4','5','6','7','8','9','c','0','⌫']

function tap(key) {
  errMsg.value = ''
  if (key === 'c') { code.value = ''; return }
  if (key === '⌫') { code.value = code.value.slice(0, -1); return }
  if (code.value.length < 20) code.value += key
}

function close() {
  code.value = ''
  errMsg.value = ''
  emit('update:visible', false)
}

async function submit() {
  const raw = code.value.trim()
  if (!raw || submitting.value) return
  submitting.value = true
  errMsg.value = ''
  try {
    // table_id 是数字格式，table_no 是字母+数字格式；两种都支持
    const tableId = Number(raw)
    if (!tableId || isNaN(tableId)) {
      errMsg.value = '房间号格式错误'
      return
    }
    await joinTable(tableId, '')
    close()
    router.push(`/table/${tableId}`)
    emit('joined', tableId)
  } catch (e) {
    errMsg.value = e.message || '房间不存在或已结束'
    code.value = ''
  } finally {
    submitting.value = false
  }
}

// 打开时自动聚焦
watch(() => props.visible, v => {
  if (v) {
    setTimeout(() => inputEl.value?.focus(), 300)
    code.value = ''
    errMsg.value = ''
  }
})
</script>

<style scoped>
.overlay {
  position: fixed; inset: 0; z-index: 500;
  background: rgba(0,0,0,.5); display: flex; align-items: flex-end;
}
.modal {
  width: 100%; background: var(--surface);
  border-radius: 20px 20px 0 0;
  padding-bottom: env(safe-area-inset-bottom);
  animation: slideUp .3s ease;
}
.modal-header {
  display: flex; align-items: center; justify-content: center;
  position: relative; padding: 18px 16px 0;
}
.modal-title { font-size: 17px; font-weight: 600; }
.close-btn {
  position: absolute; right: 16px; top: 18px;
  width: 28px; height: 28px; border-radius: 50%;
  background: var(--bg); border: none;
  display: flex; align-items: center; justify-content: center;
  font-size: 14px; color: var(--text-3); cursor: pointer;
}
.modal-body { padding: 12px 16px 0; }
.chip-img { display: block; margin: 0 auto 16px; width: 80px; height: 64px; object-fit: contain; }

.input-wrap { margin-bottom: 14px; }
.code-input {
  width: 100%; height: 52px; border-radius: 14px;
  border: 1.5px solid var(--border); padding: 0 16px;
  font-size: 22px; font-weight: 700; text-align: center;
  outline: none; background: var(--bg); color: var(--text-1);
  letter-spacing: 3px;
}
.code-input:focus { border-color: var(--primary); }
.code-input::placeholder { font-size: 14px; font-weight: 400; letter-spacing: 0; color: var(--text-3); }

.keypad {
  display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px;
  margin-bottom: 14px;
}
.key-btn {
  height: 52px; border-radius: 10px;
  background: var(--bg); border: none;
  font-size: 20px; font-weight: 500; color: var(--text-1);
  cursor: pointer; transition: all .1s;
  display: flex; align-items: center; justify-content: center;
}
.key-btn:active { background: #ddd; transform: scale(.97); }

.err-msg {
  text-align: center; font-size: 13px; color: var(--red);
  background: rgba(255,82,82,.08); border-radius: 8px;
  padding: 8px; margin-bottom: 12px;
}

.submit-btn {
  width: 100%; height: 50px; font-size: 16px; font-weight: 700;
  border-radius: 14px; margin-bottom: 16px;
}
.submit-btn:disabled { opacity: .4; }

.modal-enter-active { animation: slideUp .3s ease; }
.modal-leave-active { animation: slideUp .25s ease reverse; }
</style>
