<template>
  <div class="game-wrap">

    <!-- ══ 顶部工具栏 ══ -->
    <div class="top-bar">
      <button class="icon-btn" @click="showMenu = !showMenu">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M3 6h18M3 12h18M3 18h18" stroke="#fff" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </button>
      <div class="top-center">
        <span class="top-blind">{{ blindStr }}</span>
        <span class="top-room">{{ tableNo || `#${tableId}` }}</span>
        <span v-if="sessionRemain !== null" class="top-timer" :class="{ urgent: sessionRemain < 300 }">
          {{ formatRemain(sessionRemain) }}
        </span>
      </div>
      <div class="top-right-btns">
        <button class="icon-btn" @click="showInvite = true">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M9 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8zM23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" stroke="#fff" stroke-width="1.8" stroke-linecap="round"/>
          </svg>
        </button>
        <button class="icon-btn" @click="showRank = true">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
            <path d="M18 20V10M12 20V4M6 20v-6" stroke="#fff" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- ══ 侧边菜单 ══ -->
    <transition name="menu-slide">
      <div v-if="showMenu" class="side-menu">
        <div v-for="item in menuItems" :key="item.key" class="menu-item" @click="handleMenu(item.key)">
          <span class="menu-icon">{{ item.icon }}</span>
          <span>{{ item.label }}</span>
        </div>
      </div>
    </transition>
    <div v-if="showMenu" class="menu-mask" @click="showMenu = false" />

    <!-- ══ 牌桌主体 ══ -->
    <div class="table-scene">
      <!-- 水下气泡背景动画 -->
      <div class="bubbles-bg" aria-hidden="true">
        <span v-for="n in 12" :key="n" class="bubble" :style="`--x:${(n*7+13)%100}%;--d:${2.5+n*0.4}s;--s:${6+n%8}px;--delay:${(n*0.7)%3}s`"></span>
      </div>

      <!-- 阶段标签：椭圆外，不遮挡内容 -->
      <div v-if="game.state.stage > 0" class="stage-label-outer">{{ stageLabel }}</div>

      <div class="table-oval">

        <!-- 公共牌区域（含底池横条 + 公共牌） -->
        <div class="table-center-game" v-if="game.state.stage > 0">

          <!-- 底池横条（公共牌上方，不遮挡） -->
          <div class="pot-bar" v-if="game.state.pot > 0">
            <span ref="potBubbleEl" class="pot-bar-inner">
              <span class="pot-bar-label">总底池</span>
              <span class="pot-bar-val">{{ game.state.pot }}</span>
            </span>
            <!-- 边池 -->
            <span v-for="(sp, i) in game.state.sidePots" :key="i" class="side-pot-tag">
              边池{{ i+1 }}: {{ sp.amount }}
            </span>
          </div>

          <!-- 公共牌 -->
          <div class="community-area" v-if="game.state.stage >= 2">
            <div v-for="(card, i) in game.state.communityCards" :key="card + '_' + i"
              class="playing-card" :class="[cardColor(card), { 'card-flip': newCardIndexes.has(i) }]">
              <span class="card-rank-top">{{ cardRank(card) }}</span>
              <span class="card-suit-mid">{{ cardSuit(card) }}</span>
              <span class="card-rank-bot">{{ cardRank(card) }}</span>
            </div>
          </div>

          <!-- 桌面信息行 -->
          <div class="table-info-row">
            <span class="tinfo-item">{{ blindStr }}</span>
            <span v-if="tableNo" class="tinfo-item">{{ tableNo }}</span>
            <span v-if="sessionRemain !== null" class="tinfo-item">⏱ {{ formatRemain(sessionRemain) }}</span>
          </div>

          <!-- 发牌中旋转指示器 -->
        </div>

        <!-- 开局前中央区域 -->
        <div v-if="game.state.stage === 0" class="table-center-pre">
          <!-- 总底池标签 -->
          <div class="pre-pot-label">总底池</div>

          <!-- 创建者：开始游戏 / 等待 -->
          <template v-if="isCreator && !game.state.sessionId">
            <button v-if="seated && seatedCount >= 2" class="pre-start-btn"
              :disabled="starting" @click="tryStartSession">
              {{ starting ? '开局中...' : '开始游戏' }}
            </button>
            <div v-else class="pre-wait-hint">
              <span v-if="!seated">请先入座</span>
              <span v-else>等待更多玩家加入...</span>
            </div>
          </template>
          <template v-else-if="!game.state.sessionId">
            <div class="pre-wait-hint">等待房主开局...</div>
          </template>

          <!-- 游戏信息行 -->
          <div class="pre-meta-row">
            <span class="pre-meta-item">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2"/><path d="M12 7v5l3 3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              {{ blindStr !== '—' ? blindStr : '—' }}
            </span>
            <span class="pre-meta-item">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2"/><path d="M12 8v4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              {{ seatedCount }}/{{ maxSeats }}
            </span>
            <span v-if="tableDuration" class="pre-meta-item">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2"/><path d="M12 7v5l3 3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              {{ String(Math.floor(tableDuration * 60)).padStart(2,'0') }}:00:00
            </span>
          </div>
          <div class="pre-room-name" v-if="tableNo">{{ tableNo }}</div>

          <!-- 添加机器人（仅创建者） -->
          <button v-if="isCreator && !game.state.sessionId" class="pre-bot-btn" @click="showBotSheet = true">
            🤖 添加机器人
          </button>
        </div>

        <!-- 底池 + 阶段（替代旧 dealer-chip） -->
      </div>

      <!-- ══ 座位 ══ -->
      <div v-for="seat in seatPositions" :key="seat.no"
        class="seat-pos" :style="seat.style">

        <!-- 空座 -->
        <div v-if="!mergedPlayers[seat.no]" class="seat-empty" @click="handleSeatClick(seat.no)">
          <div class="seat-empty-circle">+</div>
        </div>

        <!-- 有玩家 -->
        <div v-else class="seat-player"
          :ref="el => { if (el) seatEls[seat.no] = el; else delete seatEls[seat.no] }"
          :class="{
            'is-turn': game.state.currentSeat === seat.no,
            'is-me': seat.no === mySeat,
            'is-fold': mergedPlayers[seat.no].status === 2
          }">

          <!-- 行动标签（下注/弃牌/All in 等） -->
          <div v-if="mergedPlayers[seat.no]?.chips === 0 && mergedPlayers[seat.no]?.status !== 2 && game.state.stage > 0"
            class="action-badge act-allin allin-pulse">ALL IN</div>
          <div v-else-if="game.seatLastAction[seat.no]"
            class="action-badge"
            :class="'act-' + game.seatLastAction[seat.no].type">
            {{ game.seatLastAction[seat.no].label }}
          </div>

          <!-- 玩家名 + 位置 + 庄家 -->
          <div class="player-name">
            <span v-if="game.state.dealerSeat === seat.no && game.state.stage > 0" class="dealer-dot">D</span>
            {{ mergedPlayers[seat.no].nickname }}
            <span v-if="seatPosition(seat.no)" class="pos-badge" :class="'pos-'+seatPosition(seat.no).toLowerCase()">
              {{ seatPosition(seat.no) }}
            </span>
          </div>

          <!-- 头像 -->
          <div class="player-avatar-wrap">
            <img v-if="mergedPlayers[seat.no].avatar" :src="mergedPlayers[seat.no].avatar" class="player-avatar" />
            <span v-else class="player-avatar-txt">{{ mergedPlayers[seat.no].nickname?.[0] || '?' }}</span>
            <!-- 行动计时环 -->
            <svg v-if="game.state.currentSeat === seat.no && timerPct > 0"
              class="timer-ring" viewBox="0 0 44 44">
              <circle cx="22" cy="22" r="20" fill="none" stroke="rgba(0,0,0,.25)" stroke-width="3"/>
              <circle cx="22" cy="22" r="20" fill="none" stroke="#f5c842" stroke-width="3"
                stroke-dasharray="125.7" :stroke-dashoffset="125.7 - timerPct * 125.7"
                stroke-linecap="round" transform="rotate(-90 22 22)"/>
            </svg>
            <!-- 弃牌遮罩 -->
            <div v-if="mergedPlayers[seat.no].status === 2" class="fold-mask">弃牌</div>
          </div>

          <!-- 筹码 -->
          <div class="player-chips">{{ mergedPlayers[seat.no].chips }}</div>

          <!-- 下注额 -->
          <div v-if="mergedPlayers[seat.no].bet" class="bet-chip"
            :class="{ 'bet-chip-fly': animatingBets.has(seat.no) }">
            🪙 {{ mergedPlayers[seat.no].bet }}
          </div>
          <!-- 手牌 -->
          <template v-if="seat.no === mySeat && game.myHoleCards?.length">
            <div class="hole-cards">
              <div v-for="(c, i) in game.myHoleCards" :key="i"
                class="playing-card hole-mine" :class="cardColor(c)">
                <span class="card-rank-top">{{ cardRank(c) }}</span>
                <span class="card-suit-mid">{{ cardSuit(c) }}</span>
                <span class="card-rank-bot">{{ cardRank(c) }}</span>
              </div>
            </div>
            <div v-if="myLiveHandRank" class="my-hand-rank-badge">{{ myLiveHandRank }}</div>
          </template>
          <!-- 摊牌阶段亮出对手底牌 -->
          <template v-else-if="showdownSeat(seat.no)">
            <div class="hole-cards">
              <div v-for="(c, i) in showdownSeat(seat.no).holeCards" :key="i"
                class="playing-card hole-mine" :class="[cardColor(c), showdownSeat(seat.no).isWinner ? 'winner-card' : '']">
                <span class="card-rank-top">{{ cardRank(c) }}</span>
                <span class="card-suit-mid">{{ cardSuit(c) }}</span>
                <span class="card-rank-bot">{{ cardRank(c) }}</span>
              </div>
            </div>
            <div class="showdown-hand-label" :class="{ 'winner-label': showdownSeat(seat.no).isWinner }">
              <template v-if="showdownSeat(seat.no).isWinner">
                +{{ showdownSeat(seat.no).winAmount }} {{ handRankName(showdownSeat(seat.no).handRank) }}赢
              </template>
              <template v-else>
                {{ handRankName(showdownSeat(seat.no).handRank) }}
              </template>
            </div>
          </template>
          <div v-else-if="shouldShowBackCards(seat.no)" class="hole-back">
            <div class="card-back"/><div class="card-back"/>
          </div>

          <!-- 赢筹码浮动动画 -->
          <transition name="chip-float">
            <div v-if="winFloats[seat.no]" class="chip-win-float">
              +{{ winFloats[seat.no] }}
            </div>
          </transition>
        </div>
      </div>
    </div>

    <!-- 提示 toast -->
    <transition name="toast-slide">
      <div v-if="errorMsg" class="game-toast">{{ errorMsg }}</div>
    </transition>

    <!-- 自定义弹窗 -->
    <transition name="dialog-fade">
      <div v-if="dialogVisible" class="game-dialog-mask" @click.self="() => closeDialog(false)">
        <div class="game-dialog">
          <div class="game-dialog-body">{{ dialogMsg }}</div>
          <div class="game-dialog-actions">
            <button v-if="dialogConfirm" class="game-dialog-btn cancel" @click="closeDialog(false)">取消</button>
            <button class="game-dialog-btn confirm" @click="closeDialog(true)">确定</button>
          </div>
        </div>
      </div>
    </transition>

    <!-- 断线提示 -->
    <div v-if="!game.connected" class="disconnected-bar">连接中...</div>

    <!-- ══ 底部操作区 ══ -->
    <div class="bottom-bar">
      <!-- 左：工具图标 -->
      <div class="bottom-left">
        <button class="act-icon-btn toggle-tools-btn" @click="showTools = !showTools">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <circle cx="5" cy="12" r="2" fill="#fff"/><circle cx="12" cy="12" r="2" fill="#fff"/><circle cx="19" cy="12" r="2" fill="#fff"/>
          </svg>
        </button>
        <transition name="tools-slide">
          <div v-if="showTools" class="tools-expanded">
            <button class="act-icon-btn" @click="showChat = !showChat">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" stroke="#fff" stroke-width="1.8" stroke-linejoin="round"/>
              </svg>
            </button>
            <button class="act-icon-btn" @click="openRebuySheet()" title="补码">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="10" r="7" stroke="#fff" stroke-width="1.8"/>
                <path d="M12 7v6M9 10h6" stroke="#fff" stroke-width="1.8" stroke-linecap="round"/>
                <path d="M5 20c0-2.5 3.1-4 7-4s7 1.5 7 4" stroke="#fff" stroke-width="1.8" stroke-linecap="round"/>
              </svg>
            </button>
            <button v-if="game.state.sessionId" class="act-icon-btn"
              @click="router.push(`/sessions/${game.state.sessionId}`)">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <rect x="4" y="3" width="16" height="18" rx="2" stroke="#fff" stroke-width="1.8"/>
                <path d="M8 8h8M8 12h8M8 16h5" stroke="#fff" stroke-width="1.8" stroke-linecap="round"/>
              </svg>
            </button>
          </div>
        </transition>
      </div>

      <!-- 右：表情 -->
      <transition name="tools-slide">
        <div v-if="showTools" class="bottom-right">
          <button class="act-icon-btn" @click="sendEmoji('👍')">👍</button>
          <button class="act-icon-btn" @click="sendEmoji('😄')">😄</button>
        </div>
      </transition>
    </div>

    <!-- ══ 行动按钮（绝对定位，对齐玩家座位水平中心） ══ -->
    <!-- 筹码为0：补码提示 -->
    <div v-if="isMyTurn && myChips === 0" class="action-float" :style="{ left: mySeatLeft }">
      <div class="rebuy-prompt">
        <div class="rebuy-prompt-text">筹码不足，请补码后继续游戏</div>
        <button class="rebuy-prompt-btn" @click="openRebuySheet()">立即补码</button>
      </div>
    </div>
    <!-- 圆形行动按钮 -->
    <div v-else-if="isMyTurn" class="action-float" :style="{ left: mySeatLeft }">
      <div v-if="showRaiseSlider" class="raise-row">
        <button class="preset" @click="raiseAmount = Math.min(myChips + (game.state.players[mySeat]?.bet||0), thirdPotRaise)">1/3</button>
        <button class="preset" @click="raiseAmount = Math.min(myChips + (game.state.players[mySeat]?.bet||0), halfPotRaise)">1/2</button>
        <button class="preset" @click="raiseAmount = Math.min(myChips + (game.state.players[mySeat]?.bet||0), twoThirdPotRaise)">2/3</button>
        <button class="preset" @click="raiseAmount = Math.min(myChips + (game.state.players[mySeat]?.bet||0), pot1Raise)">底池</button>
        <button class="preset allin" @click="raiseAmount = myChips + (game.state.players[mySeat]?.bet||0)">全押</button>
        <input type="range" :min="minRaise" :max="myChips + (game.state.players[mySeat]?.bet||0)" v-model.number="raiseAmount"
          class="raise-slider" :style="rangeStyle(raiseAmount, minRaise, myChips + (game.state.players[mySeat]?.bet||0))" />
        <span class="raise-num">{{ raiseAmount }}</span>
        <button class="close-raise" @click="showRaiseSlider = false">✕</button>
      </div>
      <div class="circle-btns-row">
        <div class="circle-btn-wrap">
          <button class="circle-btn circle-fold" @click="act('fold')">弃牌</button>
        </div>
        <div class="circle-btn-wrap circle-main-wrap">
          <button class="circle-btn circle-call" :class="{ 'timer-urgent': timerPct < 0.2 }" @click="act(canCheck ? 'check' : 'call')">
            <svg class="circle-timer-svg" viewBox="0 0 100 100">
              <circle cx="50" cy="50" r="46" fill="none" stroke="rgba(255,255,255,.15)" stroke-width="4"/>
              <circle cx="50" cy="50" r="46" fill="none"
                :stroke="timerPct < 0.2 ? '#ff4444' : 'rgba(255,255,255,.5)'" stroke-width="4"
                stroke-dasharray="289" :stroke-dashoffset="289 - timerPct * 289"
                stroke-linecap="round" transform="rotate(-90 50 50)"/>
            </svg>
            <span class="circle-btn-label">{{ canCheck ? '让牌' : `跟注` }}</span>
            <span v-if="!canCheck && callAmount" class="circle-btn-sub">{{ callAmount }}</span>
            <span class="circle-timer-sec">{{ timerSec }}s</span>
          </button>
        </div>
        <div class="circle-btn-wrap">
          <button class="circle-btn circle-raise" :class="{ 'raise-confirm': showRaiseSlider }" @click="onRaiseClick">
            <span class="circle-btn-label">{{ showRaiseSlider ? `确认` : (canCheck ? '下注' : '加注') }}</span>
            <span v-if="showRaiseSlider" class="circle-btn-sub">{{ raiseAmount }}</span>
          </button>
        </div>
      </div>
      <transition name="tools-slide">
        <div v-if="showTools" class="pre-action-row">
          <button class="pre-act-btn" :class="{ active: preAction === 'fold' }" @click="togglePreAction('fold')">下手弃牌</button>
          <button class="pre-act-btn" :class="{ active: preAction === 'check' }" @click="togglePreAction('check')">{{ canCheck ? '下手过牌' : '弃牌过牌' }}</button>
          <button class="pre-act-btn" :class="{ active: preAction === 'call' }" @click="togglePreAction('call')">跟任意注</button>
        </div>
      </transition>
    </div>

    <!-- ══ 聊天面板 ══ -->
    <transition name="chat-slide">
      <div v-if="showChat" class="chat-panel">
        <div class="chat-msgs" ref="chatMsgsEl">
          <div v-for="(m, i) in game.messages" :key="i" class="chat-msg" :class="m.type">
            <span class="msg-nick" v-if="m.type !== 'mine'">{{ m.nickname }}</span>
            <span class="msg-body">{{ m.content }}</span>
          </div>
        </div>
        <div class="chat-input-row">
          <input v-model="chatInput" class="chat-input" placeholder="说点什么..."
            @keyup.enter="sendChatMsg" />
          <button class="chat-send" @click="sendChatMsg">发送</button>
        </div>
      </div>
    </transition>

    <!-- ══ YOU WIN 弹窗 ══ -->
    <transition name="youwin">
      <div v-if="showYouWin" class="youwin-overlay" @click="showYouWin = false; game.flushPendingHand()">
        <div class="youwin-content">
          <div class="youwin-title">YOU WIN!</div>
          <div class="youwin-amount">+{{ youWinAmount }}</div>
          <div class="youwin-nick">{{ mergedPlayers[mySeat]?.nickname }}</div>
          <div class="youwin-avatar-wrap">
            <img v-if="mergedPlayers[mySeat]?.avatar" :src="mergedPlayers[mySeat].avatar" class="youwin-avatar" />
            <span v-else class="youwin-avatar-txt">{{ mergedPlayers[mySeat]?.nickname?.[0] || '?' }}</span>
          </div>
          <div v-if="youWinCards.length" class="youwin-cards">
            <div v-for="(c,i) in youWinCards" :key="i"
              class="playing-card hole-mine" :class="cardColor(c)">
              <span class="card-rank-top">{{ cardRank(c) }}</span>
              <span class="card-suit-mid">{{ cardSuit(c) }}</span>
              <span class="card-rank-bot">{{ cardRank(c) }}</span>
            </div>
          </div>
          <div v-if="youWinHandRank" class="youwin-rank">{{ youWinHandRank }}</div>
          <div class="youwin-chips">{{ mergedPlayers[mySeat]?.chips }}</div>
        </div>
      </div>
    </transition>

    <!-- ══ 牌局结束弹窗 ══ -->
    <transition name="modal">
      <div v-if="game.state.sessionStatus === 2" class="overlay session-end-overlay">
        <div class="session-end-card">
          <div class="session-end-title">牌局已结束</div>
          <div class="session-end-body">
            请在"生涯→战绩"
            <span class="session-end-link" @click="goSessionDetail">查看个人成绩</span>
          </div>
          <div class="session-end-actions">
            <button class="session-end-btn stay" @click="game.state.sessionStatus = 0">留在牌局</button>
            <button class="session-end-btn leave" @click="leaveToHome">返回大厅</button>
          </div>
        </div>
      </div>
    </transition>

    <!-- ══ 补码弹窗 ══ -->
    <transition name="modal">
      <div v-if="showRebuy" class="overlay" @click.self="myChips > 0 && (showRebuy = false)">
        <div class="sheet">
          <div class="sheet-title">补码</div>
          <div v-if="isMyTurn && myChips === 0" class="rebuy-urgent-tip">
            ⚠️ 您的筹码已用尽，补码后才能继续参与本手牌
          </div>
          <div class="wallet-balance-row">
            <span class="wallet-label">账户余额</span>
            <span class="wallet-val">🪙 {{ auth.user?.chips ?? '—' }}</span>
          </div>
          <div class="range-hint">补码范围 {{ tableMinBuyin }} ~ {{ tableMaxBuyin }}</div>
          <input type="range" :min="tableMinBuyin" :max="tableMaxBuyin" :step="tableMinBuyin"
            v-model.number="rebuyAmount" :style="rangeStyle(rebuyAmount, tableMinBuyin, tableMaxBuyin)" />
          <div class="range-val">{{ rebuyAmount }}</div>
          <button class="btn btn-primary" style="width:100%;margin-top:16px" @click="doRebuy">确认补码</button>
          <button v-if="myChips > 0" class="btn" style="width:100%;margin-top:8px;background:rgba(255,255,255,.1);color:#fff;border:none" @click="showRebuy = false">取消</button>
        </div>
      </div>
    </transition>

    <!-- ══ 入座买入弹窗 ══ -->
    <transition name="modal">
      <div v-if="showSeatBuyin" class="overlay" @click.self="showSeatBuyin = false">
        <div class="sheet">
          <div class="sheet-title">入座 — 第 {{ pendingSeatNo }} 座</div>
          <div class="wallet-balance-row">
            <span class="wallet-label">账户余额</span>
            <span class="wallet-val">🪙 {{ auth.user?.chips ?? '—' }}</span>
          </div>
          <div class="range-hint">买入范围 {{ tableMinBuyin }} ~ {{ tableMaxBuyin }}</div>
          <input type="range" :min="tableMinBuyin" :max="tableMaxBuyin" :step="tableMinBuyin"
            v-model.number="seatBuyinAmount" :style="rangeStyle(seatBuyinAmount, tableMinBuyin, tableMaxBuyin)" />
          <div class="range-val">{{ seatBuyinAmount }}</div>
          <button class="btn btn-primary" style="width:100%;margin-top:16px"
            :disabled="takingSeat" @click="confirmTakeSeat">
            {{ takingSeat ? '入座中...' : '确认入座' }}
          </button>
        </div>
      </div>
    </transition>

    <!-- ══ 邀请好友 ══ -->
    <transition name="modal">
      <div v-if="showInvite" class="overlay" @click.self="showInvite = false">
        <div class="sheet">
          <div class="sheet-title">邀请好友</div>
          <div class="invite-code-box">
            <span class="invite-code-label">房间号</span>
            <span class="invite-code-val">{{ tableNo || tableId }}</span>
            <button class="copy-btn" @click="copyCode">复制</button>
          </div>
          <div class="invite-meta">
            <span>{{ blindStr }}</span>
            <span>{{ seatedCount }}/{{ maxSeats }} 人</span>
          </div>
          <p class="invite-tip">将房间号发给好友，通过「加入牌局」输入即可加入</p>
        </div>
      </div>
    </transition>

    <!-- ══ 实时排名 ══ -->
    <transition name="modal">
      <div v-if="showRank" class="overlay" @click.self="showRank = false">
        <div class="sheet">
          <div class="sheet-title">实时排名</div>
          <div class="rank-meta">
            <span>{{ game.state.totalHands }} 手</span>
            <span>共 {{ totalBuyin }} 带入</span>
            <span v-if="sessionRemain !== null">剩余 {{ formatRemain(sessionRemain) }}</span>
          </div>
          <div v-if="rankList.length === 0" class="rank-empty">等待开局...</div>
          <div v-else class="rank-list">
            <div v-for="(p, i) in rankList" :key="p.user_id" class="rank-row">
              <span class="rank-no" :class="{ top: i < 3 }">{{ i+1 }}</span>
              <div class="rank-avatar">
                <img v-if="p.avatar" :src="p.avatar" style="width:100%;height:100%;border-radius:50%;object-fit:cover"/>
                <span v-else>{{ p.nickname?.[0] || '?' }}</span>
              </div>
              <div class="rank-info">
                <div class="rank-name">{{ p.nickname }}</div>
                <div class="rank-buyin">带入 {{ p.total_buyin }}</div>
              </div>
              <div class="rank-result" :class="p.result >= 0 ? 'green' : 'red'">
                {{ p.result >= 0 ? '+' : '' }}{{ p.result }}
              </div>
            </div>
          </div>
          <button class="btn btn-outline" style="width:100%;margin-top:16px"
            @click="showRank=false; router.push('/career')">📋 历史记录</button>
        </div>
      </div>
    </transition>

    <!-- ══ 牌型大小参考 ══ -->
    <transition name="modal">
      <div v-if="showHandsRef" class="overlay" @click.self="showHandsRef = false">
        <div class="sheet hands-sheet">
          <div class="sheet-title">🃏 牌型大小（由大到小）</div>
          <div class="hands-table">
            <div class="hand-row" v-for="(h, i) in handRankList" :key="i">
              <span class="hand-rank">{{ i+1 }}</span>
              <span class="hand-name">{{ h.name }}</span>
              <span class="hand-example">{{ h.example }}</span>
            </div>
          </div>
          <div class="hands-note">💡 不比花色大小，相同牌型比点数，再比踢脚牌</div>
          <button class="btn btn-outline" style="width:100%;margin-top:14px" @click="showHandsRef=false">关闭</button>
        </div>
      </div>
    </transition>

    <!-- ══ 添加机器人 ══ -->
    <transition name="modal">
      <div v-if="showBotSheet" class="overlay" @click.self="showBotSheet = false">
        <div class="sheet">
          <div class="sheet-title">🤖 添加机器人</div>
          <p class="bot-hint">机器人将自动入座并使用AI策略对局</p>
          <div class="bot-count-row">
            <button v-for="n in [1,2,3,4,5]" :key="n"
              class="bot-count-btn" :class="{ active: botCount === n }" @click="botCount = n">{{ n }}</button>
          </div>
          <button class="btn btn-primary" style="width:100%;margin-top:16px"
            :disabled="addingBots" @click="doAddBots">
            {{ addingBots ? '添加中...' : `添加 ${botCount} 个机器人` }}
          </button>
        </div>
      </div>
    </transition>

  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useGameStore } from '@/stores/game'
import { takeSeat as apiTakeSeat, startSession, endSession, leaveSeat, buyIn, getTableRank, getTableInfo, addBots } from '@/api'

const route  = useRoute()
const router = useRouter()
const auth   = useAuthStore()
const game   = useGameStore()

const tableId   = computed(() => Number(route.params.id))
const creatorId = ref(0)
const isCreator = computed(() =>
  creatorId.value > 0 && auth.user?.id > 0 && creatorId.value === auth.user.id
)

// ── 状态 ──────────────────────────────────────────────────────
const mySeat        = ref(-1)
const seated        = ref(false)
const errorMsg      = ref('')
let   errorTimer    = null
const dialogMsg     = ref('')
const dialogVisible = ref(false)
const dialogConfirm = ref(false)    // true = 确认弹窗（有取消按钮）
let   dialogResolve = null
function showDialog(msg) {
  dialogMsg.value = msg
  dialogConfirm.value = false
  dialogVisible.value = true
}
function showConfirm(msg) {
  return new Promise(resolve => {
    dialogMsg.value = msg
    dialogConfirm.value = true
    dialogVisible.value = true
    dialogResolve = resolve
  })
}
function closeDialog(result = false) {
  dialogVisible.value = false
  if (dialogResolve) { dialogResolve(result); dialogResolve = null }
}
const maxSeats      = ref(9)
const tableMinBuyin = ref(200)
const tableMaxBuyin = ref(2000)
const tableNo       = ref('')
const tableDuration = ref(0)
const sessionRemain = ref(null)

const restSeats = ref({})

const showMenu        = ref(false)
const showInvite      = ref(false)
const showRebuy       = ref(false)
const showTools       = ref(false)
const showSeatBuyin   = ref(false)
const showRank        = ref(false)
const showChat        = ref(false)
const showRaiseSlider = ref(false)
const showBotSheet    = ref(false)

const winFloats       = ref({}) // { seatNo: amount } — cleared after animation
const showYouWin      = ref(false)
const youWinAmount    = ref(0)
const youWinCards     = ref([])
const youWinHandRank  = ref('')
const potBubbleEl     = ref(null)
const seatEls         = reactive({}) // seatNo → DOM element

function flyChipsFromPot(winnerSeatNos, savedPotRect) {
  // Use pre-saved rect (pot element may have disappeared by the time this runs)
  const potRect = savedPotRect || potBubbleEl.value?.getBoundingClientRect()
  if (!potRect) return
  const fromX = potRect.left + potRect.width / 2
  const fromY = potRect.top  + potRect.height / 2

  winnerSeatNos.forEach((seatNo, si) => {
    const toEl = seatEls[seatNo]
    if (!toEl) return
    const toRect = toEl.getBoundingClientRect()
    const toX = toRect.left + toRect.width  / 2
    const toY = toRect.top  + toRect.height / 2

    // Spawn several chip tokens staggered
    const count = 6
    for (let i = 0; i < count; i++) {
      setTimeout(() => {
        const chip = document.createElement('div')
        chip.className = 'flying-chip'
        chip.textContent = '🪙'
        chip.style.cssText = `
          position: fixed;
          left: ${fromX}px; top: ${fromY}px;
          font-size: 18px; pointer-events: none; z-index: 600;
          transform: translate(-50%, -50%);
          transition: left 0.65s cubic-bezier(.4,0,.2,1),
                      top  0.65s cubic-bezier(.4,0,.2,1),
                      opacity 0.65s ease,
                      transform 0.65s ease;
          opacity: 1;
        `
        document.body.appendChild(chip)
        // Slight random spread from pot center
        const jx = (Math.random() - 0.5) * 24
        const jy = (Math.random() - 0.5) * 24
        requestAnimationFrame(() => {
          requestAnimationFrame(() => {
            chip.style.left    = `${toX + jx}px`
            chip.style.top     = `${toY + jy}px`
            chip.style.opacity = '0'
            chip.style.transform = 'translate(-50%,-50%) scale(0.4)'
          })
        })
        setTimeout(() => chip.remove(), 750)
      }, si * 80 + i * 55)
    }
  })
}

const pendingSeatNo   = ref(0)
const seatBuyinAmount = ref(200)
const takingSeat      = ref(false)
const starting        = ref(false)
const botCount        = ref(1)
const addingBots      = ref(false)
const rebuyAmount     = ref(200)
const raiseAmount     = ref(0)
const chatInput       = ref('')
const chatMsgsEl      = ref(null)
const timerPct        = ref(1)
const newCardIndexes  = ref(new Set())
let   prevCardCount   = 0

let timerInterval   = null
let sessionInterval = null
let rankPollTimer   = null

const menuItems = [
  { key: 'standup', icon: '🪑', label: '站起/旁观' },
  { key: 'invite',  icon: '📲', label: '邀请好友' },
  { key: 'rank',    icon: '📊', label: '实时排名' },
  { key: 'hands',   icon: '🃏', label: '牌型大小' },
  { key: 'history', icon: '📋', label: '历史记录' },
  { key: 'leave',   icon: '🚪', label: '退出牌局' },
  { key: 'disband', icon: '💔', label: '解散牌局' },
]
const showHandsRef = ref(false)

// ── 座位布局（椭圆排列）─────────────────────────────────────
const seatPositions = computed(() => {
  const all = [
    { no: 1, left: '50%', top: '74%' },
    { no: 2, left: '82%', top: '64%' },
    { no: 3, left: '88%', top: '46%' },
    { no: 4, left: '82%', top: '26%' },
    { no: 5, left: '50%', top: '13%' },
    { no: 6, left: '18%', top: '26%' },
    { no: 7, left: '12%', top: '46%' },
    { no: 8, left: '18%', top: '64%' },
    { no: 9, left: '50%', top: '52%' },
  ]
  return all.slice(0, maxSeats.value).map(p => ({
    ...p,
    style: { position: 'absolute', left: p.left, top: p.top, transform: 'translate(-50%,-50%)' }
  }))
})

// Always merge restSeats (from HTTP) with game.state.players (from WS).
// game.state.players wins for any seat present in both (has live chip/bet/status).
const mergedPlayers = computed(() => {
  return { ...restSeats.value, ...game.state.players }
})
const seatedCount = computed(() => Object.keys(mergedPlayers.value).length)

// ── 计算属性 ──────────────────────────────────────────────────
const isMyTurn = computed(() => game.state.currentSeat === mySeat.value && game.state.stage > 0)
// 玩家座位的水平位置，用于对齐行动按钮
const mySeatLeft = computed(() => {
  const pos = seatPositions.value.find(s => s.no === mySeat.value)
  return pos ? pos.left : '50%'
})
const myChips  = computed(() => mergedPlayers.value[mySeat.value]?.chips || 0)

const callAmount = computed(() => {
  const maxBet = Math.max(0, ...Object.values(game.state.players).map(p => p.bet || 0))
  return Math.max(0, maxBet - (game.state.players[mySeat.value]?.bet || 0))
})
const canCheck     = computed(() => callAmount.value === 0)
const minRaise     = computed(() => Math.max(callAmount.value * 2, game.state.bigBlind || 2))
const pot2Raise    = computed(() => (game.state.pot || 0) * 2 + callAmount.value)
const pot1Raise    = computed(() => (game.state.pot || 0) + callAmount.value)
const halfPotRaise = computed(() => Math.floor((game.state.pot || 0) / 2) + callAmount.value)

const thirdPotRaise   = computed(() => Math.floor((game.state.pot || 0) / 3) + callAmount.value)
const twoThirdPotRaise = computed(() => Math.floor((game.state.pot || 0) * 2 / 3) + callAmount.value)

const raiseLabel = computed(() =>
  showRaiseSlider.value ? `🪙 确认 ${raiseAmount.value}` : (canCheck.value ? '🪙 下注' : '🪙 加注')
)

// 计时剩余秒数
const timerSec = computed(() => Math.ceil(timerPct.value * 30))

// 预选下一手操作
const preAction = ref('')
function togglePreAction(action) {
  preAction.value = preAction.value === action ? '' : action
}
const blindStr = computed(() => {
  const sb = game.state.smallBlind
  const bb = game.state.bigBlind
  return sb && bb ? `${sb}/${bb}` : '—'
})
const totalBuyin = computed(() => game.rankList.reduce((s, p) => s + (p.total_buyin || 0), 0))
const rankList   = computed(() => game.rankList?.length ? game.rankList : [])

// 摊牌前实时牌型提示（从 showdownCards 拿自己的，否则取 game.myHandRank 如果后端有推）
const myLiveHandRank = computed(() => {
  if (game.state.stage < 2) return ''
  const sd = game.showdownCards?.find(p => p.seatNo === mySeat.value)
  if (sd?.handRank) return handRankName(sd.handRank)
  if (game.myHandRank) return handRankName(game.myHandRank)
  return ''
})

// ── 场次倒计时 ──────────────────────────────────────────────
function startSessionTimer(startedAtStr, durationHours) {
  clearInterval(sessionInterval)
  if (!startedAtStr || !durationHours) return
  const endMs = new Date(startedAtStr).getTime() + durationHours * 3600000
  const tick  = () => {
    const rem = Math.max(0, Math.floor((endMs - Date.now()) / 1000))
    sessionRemain.value = rem
    if (rem <= 0) clearInterval(sessionInterval)
  }
  tick()
  sessionInterval = setInterval(tick, 1000)
}

function formatRemain(secs) {
  if (secs === null) return ''
  const h = Math.floor(secs / 3600)
  const m = Math.floor((secs % 3600) / 60)
  const s = secs % 60
  return h > 0
    ? `${h}:${String(m).padStart(2,'0')}:${String(s).padStart(2,'0')}`
    : `${String(m).padStart(2,'0')}:${String(s).padStart(2,'0')}`
}

// ── 方法 ──────────────────────────────────────────────────────
function leaveToHome() {
  game.disconnect()
  router.push('/home')
}

function goSessionDetail() {
  if (game.state.sessionId) {
    router.push(`/sessions/${game.state.sessionId}`)
  }
}

function showdownSeat(seatNo) {
  return game.showdownCards?.find(p => p.seatNo === seatNo) || null
}

function shouldShowBackCards(seatNo) {
  const p = mergedPlayers.value[seatNo]
  return p && p.status !== 2 && game.state.gameId && seatNo !== mySeat.value
}

// ── 音效（Web Audio API，无需外部文件）─────────────────────────
let _audioCtx = null
function getAudioCtx() {
  if (!_audioCtx) _audioCtx = new (window.AudioContext || window.webkitAudioContext)()
  return _audioCtx
}
function playSound(type) {
  try {
    const ctx = getAudioCtx()
    const o = ctx.createOscillator()
    const g = ctx.createGain()
    o.connect(g); g.connect(ctx.destination)
    if (type === 'win') {
      o.frequency.setValueAtTime(523, ctx.currentTime)
      o.frequency.setValueAtTime(659, ctx.currentTime + 0.1)
      o.frequency.setValueAtTime(784, ctx.currentTime + 0.2)
      g.gain.setValueAtTime(0.15, ctx.currentTime)
      g.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.5)
      o.start(); o.stop(ctx.currentTime + 0.5)
    } else if (type === 'chip') {
      o.type = 'triangle'
      o.frequency.setValueAtTime(800, ctx.currentTime)
      o.frequency.exponentialRampToValueAtTime(400, ctx.currentTime + 0.08)
      g.gain.setValueAtTime(0.1, ctx.currentTime)
      g.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.08)
      o.start(); o.stop(ctx.currentTime + 0.1)
    } else if (type === 'fold') {
      o.type = 'sawtooth'
      o.frequency.setValueAtTime(300, ctx.currentTime)
      o.frequency.exponentialRampToValueAtTime(150, ctx.currentTime + 0.15)
      g.gain.setValueAtTime(0.08, ctx.currentTime)
      g.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.15)
      o.start(); o.stop(ctx.currentTime + 0.15)
    } else if (type === 'check' || type === 'call') {
      o.type = 'sine'
      o.frequency.setValueAtTime(440, ctx.currentTime)
      g.gain.setValueAtTime(0.08, ctx.currentTime)
      g.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.12)
      o.start(); o.stop(ctx.currentTime + 0.12)
    } else if (type === 'deal') {
      o.type = 'triangle'
      o.frequency.setValueAtTime(1200, ctx.currentTime)
      o.frequency.exponentialRampToValueAtTime(600, ctx.currentTime + 0.05)
      g.gain.setValueAtTime(0.06, ctx.currentTime)
      g.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.07)
      o.start(); o.stop(ctx.currentTime + 0.08)
    }
  } catch {}
}

function onRaiseClick() {
  if (!showRaiseSlider.value) { raiseAmount.value = minRaise.value; showRaiseSlider.value = true }
  else { act(canCheck.value ? 'bet' : 'raise', raiseAmount.value); showRaiseSlider.value = false }
}

function act(action, amount = 0) {
  game.sendAction(game.state.gameId, action, amount)
  showRaiseSlider.value = false
  // 操作音效
  if (action === 'fold') playSound('fold')
  else if (action === 'check') playSound('check')
  else if (action === 'call') playSound('call')
  else if (action === 'bet' || action === 'raise') playSound('chip')
  else if (action === 'allin') playSound('win')
}

function sendEmoji(emoji) {
  if (game.connected) game.sendChat(game.state.sessionId, emoji, 2)
}

function sendChatMsg() {
  if (!chatInput.value.trim() || !game.connected) return
  game.sendChat(game.state.sessionId, chatInput.value.trim())
  chatInput.value = ''
  nextTick(() => { if (chatMsgsEl.value) chatMsgsEl.value.scrollTop = chatMsgsEl.value.scrollHeight })
}

const rankMap = { A:'A',K:'K',Q:'Q',J:'J',T:'10' }
const suitMap = { h:'♥',d:'♦',c:'♣',s:'♠' }
function cardRank(c) { return rankMap[c?.[0]] || c?.[0] || '' }
function cardSuit(c) { return suitMap[c?.[c.length-1]] || '' }
function cardColor(c) { return (c?.[c.length-1] === 'h' || c?.[c.length-1] === 'd') ? 'red' : 'black' }

// 阶段名称
const STAGE_NAMES = { 1: '翻牌前', 2: '翻牌圈', 3: '转牌圈', 4: '河牌圈', 5: '摊牌' }
const stageLabel = computed(() => STAGE_NAMES[game.state.stage] || '')

// 中文牌型名称（后端rank: 1=高牌, 2=一对, ..., 9=同花顺, 10=皇家同花顺）
const HAND_RANK_NAMES = {
  1: '高牌', 2: '一对', 3: '两对', 4: '三条',
  5: '顺子', 6: '同花', 7: '葫芦', 8: '四条', 9: '同花顺', 10: '皇家同花顺'
}
function handRankName(rank) { return HAND_RANK_NAMES[rank] || '' }

const handRankList = [
  { name: '皇家同花顺', example: '♠A ♠K ♠Q ♠J ♠10' },
  { name: '同花顺',   example: '♥9 ♥8 ♥7 ♥6 ♥5' },
  { name: '四条',     example: '9-9-9-9-K' },
  { name: '葫芦',     example: 'Q-Q-Q-8-8' },
  { name: '同花',     example: '♦A ♦J ♦9 ♦6 ♦4' },
  { name: '顺子',     example: '5-4-3-2-A' },
  { name: '三条',     example: '7-7-7-K-2' },
  { name: '两对',     example: 'J-J-4-4-A' },
  { name: '一对',     example: '10-10-K-Q-3' },
  { name: '高牌',     example: 'A-10-7-5-2' },
]

// 座位位置标识 (D/SB/BB/UTG/HJ/CO/BTN)
// SeatOrder from engine: [1,8,7,6,5,4,3,2,9] — same clockwise direction as seat actions
const SEAT_ORDER_CLOCKWISE = [1, 8, 7, 6, 5, 4, 3, 2, 9]
function seatPosition(seatNo) {
  if (game.state.stage === 0) return ''
  const dealer = game.state.dealerSeat
  if (seatNo < 0 || dealer < 0) return ''
  // Build active seat list in clockwise order
  const activeSeatNos = SEAT_ORDER_CLOCKWISE.filter(s => mergedPlayers.value[s])
  const n = activeSeatNos.length
  if (n < 2) return ''
  const dealerIdx = activeSeatNos.indexOf(dealer)
  if (dealerIdx < 0) return ''
  // Position labels based on number of players
  const posFromDealer = (offset) => activeSeatNos[(dealerIdx + offset) % n]
  if (seatNo === dealer) return n === 2 ? 'BTN' : 'D'
  if (seatNo === posFromDealer(1)) return 'SB'
  if (seatNo === posFromDealer(2)) return 'BB'
  if (n <= 3) return ''
  if (seatNo === posFromDealer(3)) return n <= 4 ? 'CO' : (n <= 5 ? 'HJ' : 'UTG')
  if (n === 5) {
    if (seatNo === posFromDealer(4)) return 'CO'
  }
  if (n === 6) {
    if (seatNo === posFromDealer(4)) return 'HJ'
    if (seatNo === posFromDealer(5)) return 'CO'
  }
  if (n >= 7) {
    const posMap = {}
    posMap[posFromDealer(n-1)] = 'CO'
    posMap[posFromDealer(n-2)] = 'HJ'
    posMap[posFromDealer(3)] = 'UTG'
    if (n >= 8) posMap[posFromDealer(4)] = 'UTG+1'
    if (n >= 9) posMap[posFromDealer(5)] = 'MP'
    if (posMap[seatNo]) return posMap[seatNo]
  }
  return ''
}

function rangeStyle(val, min, max) {
  const pct = max === min ? 0 : ((val - min) / (max - min)) * 100
  return `background:linear-gradient(to right,#2ecc71 ${pct}%,rgba(255,255,255,.25) ${pct}%)`
}

async function handleSeatClick(seatNo) {
  if (seated.value) return
  pendingSeatNo.value = seatNo
  seatBuyinAmount.value = tableMinBuyin.value
  auth.fetchProfile()
  showSeatBuyin.value = true
}

async function confirmTakeSeat() {
  if (takingSeat.value) return
  takingSeat.value = true
  try {
    await apiTakeSeat(tableId.value, pendingSeatNo.value, seatBuyinAmount.value)
    mySeat.value  = pendingSeatNo.value
    seated.value  = true
    showSeatBuyin.value = false
    auth.fetchProfile()
    restSeats.value = {
      ...restSeats.value,
      [pendingSeatNo.value]: {
        seat: pendingSeatNo.value, user_id: auth.user?.id || 0,
        nickname: auth.user?.nickname || '我', avatar: auth.user?.avatar || '',
        chips: seatBuyinAmount.value,
      }
    }
  } catch (e) { showDialog(e.message || '入座失败') }
  finally { takingSeat.value = false }
}

async function tryStartSession() {
  if (starting.value) return
  starting.value = true
  try {
    const data = await startSession(tableId.value)
    if (data?.session_id) {
      game.state.sessionId = data.session_id
      startSessionTimer(new Date().toISOString(), tableDuration.value)
    }
  } catch (e) { if (!e.message?.includes('已开局')) showDialog(e.message || '开局失败') }
  finally { starting.value = false }
}

async function doAddBots() {
  if (addingBots.value) return
  addingBots.value = true
  try {
    const data = await addBots(tableId.value, botCount.value)
    showBotSheet.value = false
    await loadTableInfo()
    showDialog(`已添加 ${data?.count || botCount.value} 个机器人`)
  } catch (e) { showDialog(e.message || '添加失败') }
  finally { addingBots.value = false }
}

async function handleMenu(key) {
  showMenu.value = false
  if (key === 'invite')  { showInvite.value = true; return }
  if (key === 'history') { router.push('/career'); return }
  if (key === 'rank')    { showRank.value = true; await loadRank(); return }
  if (key === 'hands')   { showHandsRef.value = true; return }
  if (key === 'standup') {
    try {
      await leaveSeat(tableId.value)
      const copy = { ...restSeats.value }
      delete copy[mySeat.value]
      restSeats.value = copy
      mySeat.value = -1; seated.value = false
      showToast('已站起，筹码已结算回账户')
    } catch (e) { showToast(e.message || '操作失败') }
    return
  }
  if (key === 'leave') {
    try { await leaveSeat(tableId.value) } catch {}
    game.disconnect(); router.push('/home'); return
  }
  if (key === 'disband') {
    if (!await showConfirm('确认解散牌局？')) return
    try { if (game.state.sessionId) await endSession(game.state.sessionId, 2) } catch {}
    game.disconnect(); router.push('/home')
  }
}

function openRebuySheet() {
  auth.fetchProfile()
  showRebuy.value = true
}

async function doRebuy() {
  try {
    const sid = game.state.sessionId
    if (!sid) { showDialog('场次未开始'); return }
    await buyIn(sid, rebuyAmount.value)
    showRebuy.value = false
    await auth.fetchProfile()
    showToast(`补码 ${rebuyAmount.value} 成功`)
  } catch (e) { showDialog(e.message || '操作失败') }
}

function copyCode() {
  navigator.clipboard?.writeText(tableNo.value || String(tableId.value))
  showToast('房间号已复制，去发给好友吧 🎉')
}

async function loadRank() {
  try {
    const data = await getTableRank(tableId.value)
    if (data?.players) game.rankList.splice(0, Infinity, ...data.players)
  } catch {}
}

async function loadTableInfo() {
  try {
    const data = await getTableInfo(tableId.value)
    if (!data) return
    tableNo.value       = data.table_no   || ''
    creatorId.value     = data.creator_id || 0
    maxSeats.value      = data.max_seats  || 9
    tableMinBuyin.value = data.min_buyin  || 200
    tableMaxBuyin.value = data.max_buyin  || 2000
    tableDuration.value = data.duration   || 0
    seatBuyinAmount.value = data.min_buyin || 200
    rebuyAmount.value   = data.min_buyin  || 200

    if (!game.state.smallBlind && data.small_blind) game.state.smallBlind = data.small_blind
    if (!game.state.bigBlind   && data.big_blind)   game.state.bigBlind   = data.big_blind

    const seats = {}
    for (const s of (data.seats || [])) {
      seats[s.seat_no] = { seat: s.seat_no, user_id: s.user_id, nickname: s.nickname, avatar: s.avatar, chips: s.chips, status: 1 }
    }
    restSeats.value = seats

    const myInfo = auth.user
    if (myInfo) {
      for (const s of (data.seats || [])) {
        if (s.user_id === myInfo.id) { mySeat.value = s.seat_no; seated.value = true; break }
      }
    }

    if (data.session_id) {
      game.state.sessionId = data.session_id
    }
    if (data.session_status != null && data.session_status === 2) {
      game.state.sessionStatus = 2
    }
    if (data.session_id && data.started_at && tableDuration.value > 0 && data.session_status !== 2) {
      startSessionTimer(data.started_at, tableDuration.value)
    }
  } catch {}
}

function startTimer(deadline) {
  clearInterval(timerInterval)
  const total = 30000
  const end   = deadline || (Date.now() + total)
  timerInterval = setInterval(() => {
    timerPct.value = Math.max(0, (end - Date.now()) / total)
    if (timerPct.value <= 0) clearInterval(timerInterval)
  }, 100)
}

function showToast(msg) {
  errorMsg.value = msg
  clearTimeout(errorTimer)
  errorTimer = setTimeout(() => { errorMsg.value = '' }, 3000)
}

watch(() => game.lastError, msg => {
  if (!msg) return
  showToast(msg)
  game.lastError = ''
})

// 翻牌动画：检测新增的公共牌
watch(() => game.state.communityCards, (cards) => {
  const cur = cards?.length || 0
  if (cur > prevCardCount) {
    const indexes = new Set()
    for (let i = prevCardCount; i < cur; i++) indexes.add(i)
    newCardIndexes.value = indexes
    setTimeout(() => { newCardIndexes.value = new Set() }, 700)
  }
  prevCardCount = cur
}, { deep: true })

// 下注筹码飞入动画：记录新增下注
const animatingBets = ref(new Set())
watch(() => game.state.players, (players, old) => {
  for (const [seatNo, p] of Object.entries(players || {})) {
    const prev = old?.[seatNo]?.bet || 0
    if (p.bet > prev) {
      animatingBets.value.add(Number(seatNo))
      setTimeout(() => {
        animatingBets.value.delete(Number(seatNo))
      }, 500)
    }
  }
}, { deep: true })

// 赢筹码浮动动画 + 飞筹码 + YOU WIN 弹窗
watch(() => game.lastResult, result => {
  if (!result?.players) return

  // 提前记录底池位置（hand_result 后 pot 变 0，pot-bar 消失，ref 变 null）
  const savedPotRect = potBubbleEl.value?.getBoundingClientRect() ?? null

  const floats = {}
  const winnerSeats = []
  for (const p of result.players) {
    if (p.is_winner) {
      floats[p.seat_no] = p.win_amount || 0
      winnerSeats.push(p.seat_no)
    }
  }
  const myResult = result.players?.find(p => p.seat_no === mySeat.value)
  if (myResult?.is_winner) {
    youWinAmount.value   = myResult.win_amount || 0
    youWinCards.value    = game.lastDealCards?.length ? [...game.lastDealCards] : []
    const sd = game.showdownCards?.find(p => p.seatNo === mySeat.value)
    youWinHandRank.value = sd ? handRankName(sd.handRank) : ''
    game.holdNewHand()
  }

  setTimeout(() => {
    winFloats.value = floats
    setTimeout(() => { winFloats.value = {} }, 3000)

    // 筹码飞入动画
    nextTick(() => flyChipsFromPot(winnerSeats, savedPotRect))

    // 播放赢钱音效
    if (winnerSeats.length > 0) playSound('win')

    // YOU WIN 弹窗
    if (myResult?.is_winner) {
      showYouWin.value = true
      setTimeout(() => {
        showYouWin.value = false
        game.flushPendingHand()
      }, 2500)
    } else {
      // 非赢家不显示弹窗，直接放行下一手
      game.flushPendingHand()
    }
  }, 500)
})

watch(() => game.state.actionDeadline, dl => { if (dl > 0) startTimer(dl) })
watch(() => game.state.currentSeat, seat => {
  if (seat === mySeat.value && game.state.stage > 0) {
    startTimer(game.state.actionDeadline)
    // 筹码为0时自动弹出补码
    if (myChips.value === 0) {
      showRebuy.value = true
      return
    }
    // 执行预选操作
    if (preAction.value) {
      const pa = preAction.value
      preAction.value = ''
      nextTick(() => {
        if (pa === 'fold') { act('fold'); playSound('fold') }
        else if (pa === 'check' && canCheck.value) { act('check'); playSound('check') }
        else if (pa === 'call') { act(canCheck.value ? 'check' : 'call'); playSound('call') }
      })
    }
  }
})
watch(() => game.state.players, players => {
  if (Object.keys(players).length > 0) {
    for (const [seatNo, p] of Object.entries(players)) {
      if (!p.avatar && restSeats.value[seatNo]?.avatar) p.avatar = restSeats.value[seatNo].avatar
    }
  }
}, { deep: true })

onMounted(async () => {
  await auth.fetchProfile()
  await loadTableInfo()
  game.connect(tableId.value)
  game.startPing()
  await loadRank()
  rankPollTimer = setInterval(loadRank, 15000)
})

onUnmounted(() => {
  game.stopPing()
  clearInterval(timerInterval)
  clearInterval(sessionInterval)
  clearInterval(rankPollTimer)
})
</script>

<style scoped>
/* ── 全局 ── */
* { box-sizing: border-box; }

/* ── Toast ── */
.game-toast {
  position: fixed; top: 64px; left: 50%; transform: translateX(-50%);
  background: rgba(10,10,10,.82); color: #fff;
  padding: 9px 22px; border-radius: 24px; font-size: 13px;
  letter-spacing: .3px; z-index: 300; pointer-events: none;
  backdrop-filter: blur(6px);
  border: 1px solid rgba(255,255,255,.12);
  white-space: nowrap;
}
.toast-slide-enter-active { animation: toastIn .25s ease; }
.toast-slide-leave-active { animation: toastIn .2s ease reverse; }
@keyframes toastIn {
  from { opacity: 0; transform: translateX(-50%) translateY(-8px); }
  to   { opacity: 1; transform: translateX(-50%) translateY(0); }
}

/* ── 自定义弹窗 ── */
.game-dialog-mask {
  position: fixed; inset: 0; z-index: 400;
  background: rgba(0,0,0,.55);
  display: flex; align-items: center; justify-content: center;
  backdrop-filter: blur(3px);
}
.game-dialog {
  background: linear-gradient(160deg, #1e3a28 0%, #142b1e 100%);
  border: 1px solid rgba(14,196,176,.35);
  border-radius: 20px;
  padding: 28px 24px 20px;
  min-width: 260px; max-width: 80vw;
  box-shadow: 0 8px 32px rgba(0,0,0,.6), 0 0 0 1px rgba(255,255,255,.06) inset;
  display: flex; flex-direction: column; align-items: center; gap: 20px;
}
.game-dialog-body {
  color: #e8f5e9; font-size: 15px; line-height: 1.6;
  text-align: center; letter-spacing: .3px;
}
.game-dialog-actions {
  display: flex; gap: 12px; width: 100%; justify-content: center;
}
.game-dialog-btn {
  flex: 1; max-width: 130px;
  border: none; border-radius: 24px;
  padding: 11px 0; font-size: 15px; font-weight: 600;
  cursor: pointer; letter-spacing: .5px;
  transition: opacity .15s, transform .1s;
}
.game-dialog-btn.confirm {
  background: linear-gradient(135deg, #0ec4b0, #09a090);
  color: #fff;
  box-shadow: 0 4px 14px rgba(14,196,176,.4);
}
.game-dialog-btn.cancel {
  background: rgba(255,255,255,.08);
  color: rgba(255,255,255,.7);
  border: 1px solid rgba(255,255,255,.15);
}
.game-dialog-btn:active { opacity: .8; transform: scale(.96); }
.dialog-fade-enter-active { animation: dialogIn .2s ease; }
.dialog-fade-leave-active { animation: dialogIn .15s ease reverse; }
@keyframes dialogIn {
  from { opacity: 0; transform: scale(.9); }
  to   { opacity: 1; transform: scale(1); }
}
.disconnected-bar {
  position: fixed; top: 48px; left: 0; right: 0;
  background: rgba(0,0,0,.6); color: #f5c842;
  text-align: center; font-size: 12px; padding: 3px 0; z-index: 200;
}
.fade-enter-active, .fade-leave-active { transition: opacity .3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.game-wrap {
  position: fixed; inset: 0;
  background: #1a6b3c;
  overflow: hidden; display: flex; flex-direction: column;
  font-family: -apple-system, 'PingFang SC', sans-serif;
}

/* 背景纹理 */
.bubbles-bg { display: none; }
@keyframes bubbleRise { to {} }

/* ── 顶部工具栏 ── */
.top-bar {
  flex-shrink: 0; height: 48px;
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 10px;
  background: rgba(0,0,0,.35);
  border-bottom: 1px solid rgba(255,255,255,.1);
  z-index: 50;
}
.icon-btn {
  width: 34px; height: 34px; border-radius: 17px;
  background: rgba(255,255,255,.12); border: none;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer;
}
.icon-btn:active { opacity: .7; }
.top-right-btns { display: flex; gap: 8px; }

.top-center {
  flex: 1; display: flex; align-items: center; justify-content: center;
  gap: 10px; font-size: 13px; color: rgba(255,255,255,.9);
}
.top-blind { font-weight: 700; color: #f5c842; }
.top-room  { color: rgba(255,255,255,.7); font-size: 12px; }
.top-timer { font-weight: 700; color: #fff; letter-spacing: .5px; }
.top-timer.urgent { color: #ff5252; animation: pulse 1s infinite; }
@keyframes pulse { 0%,100%{opacity:1} 50%{opacity:.4} }

/* ── 侧边菜单 ── */
.side-menu {
  position: absolute; top: 0; left: 0; bottom: 0; z-index: 200;
  width: 160px;
  background: rgba(2,10,35,.97);
  border-right: 1px solid rgba(50,120,220,.25);
  padding: 56px 0 20px;
  box-shadow: 4px 0 30px rgba(0,0,0,.7);
}
.menu-item {
  display: flex; align-items: center; gap: 12px;
  padding: 13px 20px; color: #fff; font-size: 14px; cursor: pointer;
}
.menu-item:active { background: rgba(255,255,255,.08); }
.menu-icon { font-size: 18px; }
.menu-mask { position: absolute; inset: 0; z-index: 190; background: rgba(0,0,0,.35); }

/* ── 牌桌场景 ── */
.table-scene {
  flex: 1; position: relative;
  padding-bottom: 110px; /* 为行动按钮留出空间，防止遮挡手牌 */
}

/* 椭圆桌面 — 绿色绒布 */
.table-oval {
  position: absolute;
  left: 50%; top: 50%;
  transform: translate(-50%, -50%);
  width: 76%; height: 58%;
  border-radius: 50%;
  background: radial-gradient(ellipse at 50% 40%, #2e8b57 0%, #1d6b42 60%, #165233 100%);
  border: 10px solid #8B6914;
  box-shadow:
    0 0 0 3px #a07820,
    inset 0 0 60px rgba(0,0,0,.25),
    0 12px 40px rgba(0,0,0,.6);
  display: flex; align-items: center; justify-content: center;
}
@keyframes tableGlow { to {} }

/* 游戏中央区域（底池+公共牌+信息） */
.table-center-game {
  display: flex; flex-direction: column; align-items: center; gap: 8px;
  position: absolute; top: 22%; left: 50%; transform: translateX(-50%);
  pointer-events: none;
}

/* 底池横条 */
.pot-bar {
  display: flex; justify-content: center;
}
.pot-bar-inner {
  display: inline-flex; align-items: center; gap: 6px;
  background: rgba(0,0,0,.45);
  border-radius: 16px; padding: 4px 16px;
  border: 1px solid rgba(255,255,255,.2);
  backdrop-filter: blur(4px);
}
.pot-bar-label {
  font-size: 11px; color: rgba(255,255,255,.65); letter-spacing: 1px;
}
.pot-bar-val {
  font-size: 16px; font-weight: 800; color: #fff;
}

/* 公共牌 */
.community-area {
  display: flex; gap: 8px;
}

/* ── 牌通用 ── */
.playing-card {
  width: 46px; height: 66px; border-radius: 7px;
  background: #fff;
  display: flex; flex-direction: column;
  align-items: flex-start; justify-content: space-between;
  padding: 3px 5px;
  font-weight: 800; line-height: 1;
  box-shadow: 0 4px 12px rgba(0,0,0,.6), inset 0 1px 0 rgba(255,255,255,.9);
  position: relative; overflow: hidden;
}
.playing-card.red  { color: #cc0000; }
.playing-card.black { color: #111; }

/* 左上角数字 */
.card-rank-top {
  font-size: 15px; font-weight: 900; line-height: 1;
  align-self: flex-start;
}
/* 中央大花色 */
.card-suit-mid {
  font-size: 26px; line-height: 1;
  position: absolute; top: 50%; left: 50%;
  transform: translate(-50%, -50%);
}
/* 右下角倒置数字 */
.card-rank-bot {
  font-size: 15px; font-weight: 900; line-height: 1;
  align-self: flex-end; transform: rotate(180deg);
}

/* 自己的手牌 — 较大 */
.playing-card.hole-mine {
  width: 60px; height: 84px;
}
.playing-card.hole-mine .card-rank-top,
.playing-card.hole-mine .card-rank-bot { font-size: 20px; }
.playing-card.hole-mine .card-suit-mid { font-size: 34px; }

/* ── 底池 ── */
.pot-area {
  position: absolute; bottom: 26%; left: 50%; transform: translateX(-50%);
  display: flex; flex-direction: column; align-items: center; gap: 5px;
  pointer-events: none;
}
.pot-total {
  background: rgba(0,0,0,.55);
  color: #fff; font-size: 18px; font-weight: 800;
  border-radius: 20px; padding: 3px 18px;
  letter-spacing: .5px; min-width: 80px; text-align: center;
  box-shadow: 0 2px 8px rgba(0,0,0,.4);
}
.pot-chips-row {
  display: flex; gap: 8px; align-items: center;
}
.pot-chip-bubble {
  display: flex; align-items: center; gap: 4px;
  border-radius: 12px; padding: 2px 10px;
  font-size: 12px; font-weight: 700;
}
.pot-chip-blue  { background: rgba(30,80,200,.7); color: #fff; }
.pot-chip-teal  { background: rgba(0,150,136,.7); color: #fff; }
.pot-chip-dot   { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.blue-dot       { background: #5b9dff; border: 2px solid #fff; }
.teal-dot       { background: #4dd0c4; border: 2px solid #fff; }

/* ── 翻牌动画 ── */
@keyframes cardDeal {
  0%   { transform: rotateY(90deg) scale(.6) translateY(-20px); opacity: 0; }
  50%  { transform: rotateY(45deg) scale(.9) translateY(-5px);  opacity: .7; }
  100% { transform: rotateY(0deg)  scale(1)  translateY(0);     opacity: 1; }
}
.card-flip {
  animation: cardDeal .5s cubic-bezier(.34,1.56,.64,1) both;
}
.community-area .playing-card.card-flip:nth-child(2) { animation-delay: .10s; }
.community-area .playing-card.card-flip:nth-child(3) { animation-delay: .20s; }

/* 发牌旋转指示器 */

/* 我的手牌实时牌型 */
.my-hand-rank-badge {
  font-size: 12px; font-weight: 700; color: #fff;
  background: rgba(0,0,0,.8);
  border-radius: 12px; padding: 3px 12px; margin-top: 4px;
  letter-spacing: .5px; text-align: center;
  box-shadow: 0 2px 8px rgba(0,0,0,.5);
  animation: badgePop .3s cubic-bezier(.34,1.56,.64,1) both;
}
@keyframes badgePop {
  from { transform: scale(.6); opacity: 0; }
  to   { transform: scale(1);  opacity: 1; }
}

/* ── 开局前中央内容 ── */
.table-center-pre {
  display: flex; flex-direction: column; align-items: center; gap: 12px;
}

/* 总底池标签 */
.pre-pot-label {
  background: rgba(0,0,0,.35);
  color: rgba(255,255,255,.85);
  font-size: 12px; font-weight: 600; letter-spacing: 2px;
  padding: 4px 20px; border-radius: 14px;
  border: 1px solid rgba(255,255,255,.2);
}

/* 开始游戏大按钮 */
.pre-start-btn {
  height: 52px; padding: 0 56px; border-radius: 26px;
  background: rgba(255,255,255,.18);
  border: 1.5px solid rgba(255,255,255,.45);
  color: #fff; font-size: 18px; font-weight: 700; cursor: pointer;
  letter-spacing: 2px;
  box-shadow: 0 4px 20px rgba(0,0,0,.3);
  transition: all .2s;
  backdrop-filter: blur(4px);
}
.pre-start-btn:hover   { background: rgba(255,255,255,.25); }
.pre-start-btn:active  { transform: scale(.97); }
.pre-start-btn:disabled { opacity: .45; cursor: default; }

.pre-wait-hint {
  font-size: 14px; color: rgba(255,255,255,.6); letter-spacing: .5px;
}

/* 信息行 */
.pre-meta-row {
  display: flex; align-items: center; gap: 16px;
  font-size: 12px; color: rgba(255,255,255,.55);
}
.pre-meta-item {
  display: flex; align-items: center; gap: 4px;
}

.pre-room-name {
  font-size: 12px; color: rgba(255,255,255,.45); letter-spacing: .5px;
}

.pre-bot-btn {
  height: 32px; padding: 0 16px; border-radius: 16px;
  background: rgba(255,255,255,.12); border: 1px solid rgba(255,255,255,.25);
  color: rgba(255,255,255,.8); font-size: 12px; cursor: pointer;
  transition: background .2s;
}
.pre-bot-btn:active { background: rgba(255,255,255,.22); }

/* 兼容旧 class（可能还在用） */
.start-btn { display: none; }
.wait-hint { display: none; }
.bot-btn   { display: none; }
.table-meta-bar   { display: none; }
.invite-center-btn { display: none; }

/* 庄家按钮 */
.dealer-chip {
  position: absolute; right: 52%; bottom: 35%;
  width: 22px; height: 22px; border-radius: 50%;
  background: #fff; color: #1a1a1a;
  font-size: 10px; font-weight: 900;
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 2px 6px rgba(0,0,0,.4);
}

/* ── 座位 ── */
.seat-pos { z-index: 10; }

.seat-empty {
  display: flex; flex-direction: column; align-items: center; cursor: pointer;
}
.seat-empty-circle {
  width: 56px; height: 56px; border-radius: 50%;
  border: 2px dashed rgba(255,255,255,.28);
  background: rgba(0,0,0,.28);
  display: flex; align-items: center; justify-content: center;
  font-size: 22px; color: rgba(255,255,255,.4);
  transition: all .2s;
  box-shadow: inset 0 0 12px rgba(0,0,0,.3);
}
.seat-empty:active .seat-empty-circle {
  background: rgba(255,255,255,.12);
  border-color: rgba(255,255,255,.5);
}

.seat-player {
  display: flex; flex-direction: column; align-items: center; gap: 2px;
  position: relative;
}
.seat-player.is-fold { opacity: .5; }

/* 行动标签 */
.action-badge {
  font-size: 13px; font-weight: 700; border-radius: 12px;
  padding: 3px 11px; margin-bottom: 3px; white-space: nowrap;
  text-align: center; letter-spacing: .3px;
  box-shadow: 0 2px 8px rgba(0,0,0,.4);
}
.act-fold    { background: rgba(80,80,80,.9);  color: #ccc; }
.act-check   { background: rgba(30,120,200,.9); color: #fff; }
.act-call    { background: rgba(30,150,80,.9);  color: #fff; }
.act-bet     { background: rgba(0,160,140,.9);  color: #fff; }
.act-raise   { background: rgba(200,120,0,.9);  color: #fff; }
.act-allin   { background: linear-gradient(135deg,#e53935,#c62828); color: #fff; }
.act-blind   { background: rgba(60,60,60,.8);   color: #aaa; }
.act-ante    { background: rgba(60,60,60,.8);   color: #aaa; }

.player-name {
  font-size: 11px; color: #fff; font-weight: 600;
  white-space: nowrap; max-width: 80px;
  overflow: hidden; text-overflow: ellipsis;
  text-shadow: 0 1px 4px rgba(0,0,0,.9);
  display: flex; align-items: center; gap: 3px;
}

/* 庄家标识 */
.dealer-dot {
  width: 16px; height: 16px; border-radius: 50%;
  background: #fff; color: #1a1a1a;
  font-size: 9px; font-weight: 900;
  display: inline-flex; align-items: center; justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 1px 4px rgba(0,0,0,.5);
}

.player-avatar-wrap {
  position: relative; width: 56px; height: 56px; border-radius: 50%;
  overflow: hidden;
  border: 3px solid rgba(255,255,255,.3);
  background: rgba(20,20,20,.6);
  box-shadow: 0 3px 12px rgba(0,0,0,.7);
  transition: border-color .3s, box-shadow .3s;
}
.seat-player.is-turn .player-avatar-wrap {
  border-color: #f5c842;
  box-shadow: 0 0 0 3px rgba(245,200,66,.4), 0 0 20px rgba(245,200,66,.6);
  animation: turnPulse 1s ease-in-out infinite;
}
.seat-player.is-me .player-avatar-wrap {
  border-color: #4dd0e1;
  box-shadow: 0 0 0 2px rgba(77,208,225,.4), 0 3px 12px rgba(0,0,0,.7);
}
@keyframes turnPulse {
  0%,100% { box-shadow: 0 0 0 3px rgba(245,200,66,.4), 0 0 20px rgba(245,200,66,.6); }
  50%     { box-shadow: 0 0 0 5px rgba(245,200,66,.2), 0 0 32px rgba(245,200,66,.8); }
}

.player-avatar {
  width: 100%; height: 100%; object-fit: cover;
  display: block;
}
.player-avatar-txt {
  position: absolute; inset: 0;
  display: flex; align-items: center; justify-content: center;
  font-size: 20px; font-weight: 700; color: #fff;
}

.timer-ring {
  position: absolute; inset: -3px; width: calc(100% + 6px); height: calc(100% + 6px);
  pointer-events: none; border-radius: 50%; overflow: visible;
}

.fold-mask {
  position: absolute; inset: 0;
  background: rgba(0,0,0,.6);
  display: flex; align-items: center; justify-content: center;
  font-size: 11px; font-weight: 700; color: #ccc;
  letter-spacing: .5px; border-radius: 50%;
}

.player-chips {
  background: rgba(0,0,0,.75);
  color: #fff; font-weight: 700;
  font-size: 12px; border-radius: 12px; padding: 3px 10px;
  white-space: nowrap; min-width: 44px; text-align: center;
  box-shadow: 0 2px 6px rgba(0,0,0,.6);
  border: 1px solid rgba(255,255,255,.15);
}

.bet-chip {
  position: absolute; top: -22px;
  background: linear-gradient(135deg, #f5c842, #e6a800);
  color: #1a1a1a;
  border-radius: 12px; padding: 3px 10px;
  font-size: 13px; font-weight: 800;
  box-shadow: 0 2px 10px rgba(245,200,66,.6), 0 1px 4px rgba(0,0,0,.4);
  white-space: nowrap;
}
.bet-chip-fly {
  animation: betFly .4s cubic-bezier(.34,1.56,.64,1) both;
}
@keyframes betFly {
  0%   { transform: scale(.5) translateY(10px); opacity: 0; }
  100% { transform: scale(1)  translateY(0);    opacity: 1; }
}

.hole-cards { display: flex; gap: 4px; margin-top: 4px; }
.hole-back  { display: flex; gap: 4px; margin-top: 4px; }
.showdown-hand-label {
  font-size: 11px; color: #fff; font-weight: 700;
  background: rgba(0,0,0,.8); border-radius: 12px;
  padding: 3px 10px; margin-top: 4px; text-align: center;
  white-space: nowrap;
}
.showdown-hand-label.winner-label {
  background: rgba(39,174,96,.85); color: #fff; font-size: 12px;
}
.playing-card.winner-card {
  box-shadow: 0 0 10px 4px #69f0ae, 0 0 24px 8px rgba(105,240,174,.4), 0 2px 8px rgba(0,0,0,.5);
  animation: winnerGlow 1.2s ease-in-out infinite;
}
@keyframes winnerGlow {
  0%,100% { box-shadow: 0 0 10px 4px #69f0ae, 0 0 24px 8px rgba(105,240,174,.4), 0 2px 8px rgba(0,0,0,.5); }
  50%     { box-shadow: 0 0 16px 6px #69f0ae, 0 0 40px 14px rgba(105,240,174,.6), 0 2px 8px rgba(0,0,0,.5); }
}
.card-back {
  width: 30px; height: 44px; border-radius: 5px;
  background: linear-gradient(135deg, #c0392b 0%, #96281b 50%, #c0392b 100%);
  border: 2px solid rgba(255,255,255,.3);
  box-shadow: 0 2px 8px rgba(0,0,0,.6);
  position: relative; overflow: hidden;
}
.card-back::after {
  content: '';
  position: absolute; inset: 3px;
  border: 1px solid rgba(255,255,255,.2);
  border-radius: 2px;
  background: repeating-linear-gradient(
    45deg,
    transparent,
    transparent 3px,
    rgba(255,255,255,.05) 3px,
    rgba(255,255,255,.05) 6px
  );
}

/* ── 底部操作区 ── */
.bottom-bar {
  flex-shrink: 0;
  display: flex; align-items: flex-end; justify-content: space-between;
  padding: 6px 10px 14px;
  background: rgba(0,0,0,.4);
  border-top: 1px solid rgba(255,255,255,.1);
  z-index: 50;
  --bottom-bar-h: 90px;
}
.bottom-left { display: flex; flex-direction: column; gap: 8px; align-items: center; }
.bottom-right { display: flex; flex-direction: column; gap: 8px; }
.tools-expanded { display: flex; flex-direction: column; gap: 8px; }
.tools-slide-enter-active, .tools-slide-leave-active { transition: opacity .2s, transform .2s; }
.tools-slide-enter-from, .tools-slide-leave-to { opacity: 0; transform: translateY(8px); }
.act-icon-btn {
  width: 40px; height: 40px; border-radius: 50%;
  background: rgba(255,255,255,.15); border: none; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  font-size: 17px;
}
.act-icon-btn:active { opacity: .7; }

.action-area { display: flex; flex-direction: column; gap: 8px; align-items: center; }
.action-float {
  position: absolute;
  bottom: 10px;
  transform: translateX(-50%);
  display: flex; flex-direction: column; gap: 8px; align-items: center;
  z-index: 55;
  pointer-events: auto;
}

.raise-row {
  display: flex; align-items: center; gap: 6px;
  background: rgba(0,0,0,.55); border-radius: 20px; padding: 6px 12px;
  max-width: 340px; flex-wrap: wrap; justify-content: center;
}
.preset {
  background: rgba(255,255,255,.15); border: none; color: #fff;
  border-radius: 12px; height: 26px; padding: 0 10px;
  font-size: 11px; cursor: pointer;
}
.preset:active { background: rgba(46,204,113,.4); }
.raise-slider { flex: 1; min-width: 80px; }
.raise-num { color: #f5c842; font-size: 14px; font-weight: 700; min-width: 44px; text-align: center; }
.close-raise { background: none; border: none; color: rgba(255,255,255,.4); font-size: 14px; cursor: pointer; }

.action-btns-row { display: none; } /* replaced by circle-btns-row */
@keyframes callBtnPulse { to {} }

/* ── 聊天 ── */
.chat-panel {
  position: absolute; bottom: 70px; left: 10px; right: 10px; z-index: 60;
  background: rgba(0,0,0,.75); border-radius: 14px;
  display: flex; flex-direction: column; max-height: 200px;
}
.chat-msgs {
  flex: 1; overflow-y: auto; padding: 10px 12px;
  display: flex; flex-direction: column; gap: 6px;
}
.chat-msg { display: flex; flex-wrap: wrap; gap: 4px; align-items: baseline; }
.chat-msg.mine { flex-direction: row-reverse; }
.msg-nick { font-size: 10px; color: rgba(255,255,255,.5); }
.msg-body {
  font-size: 12px; color: #fff;
  background: rgba(255,255,255,.12); border-radius: 10px; padding: 3px 10px;
}
.chat-msg.mine .msg-body { background: rgba(39,174,96,.5); }
.chat-input-row { display: flex; border-top: 1px solid rgba(255,255,255,.1); }
.chat-input {
  flex: 1; background: transparent; border: none; outline: none;
  color: #fff; font-size: 13px; padding: 8px 12px;
}
.chat-input::placeholder { color: rgba(255,255,255,.3); }
.chat-send {
  background: #27ae60; border: none; color: #fff;
  padding: 8px 14px; font-size: 13px; cursor: pointer; border-radius: 0 0 14px 0;
}

/* ── 弹窗 ── */
.result-overlay {
  position: absolute; inset: 0; z-index: 300;
  background: rgba(0,0,0,.65);
  display: flex; align-items: center; justify-content: center;
}
.result-card {
  background: var(--surface); border-radius: 16px;
  padding: 20px; min-width: 260px; max-width: 90%;
}
.result-title { font-size: 16px; font-weight: 700; text-align: center; margin-bottom: 12px; }
.pot-result { margin-bottom: 12px; padding-bottom: 12px; border-bottom: 1px solid var(--border); }
.pot-result:last-child { border-bottom: none; }
.pot-type { font-size: 12px; color: var(--text-3); margin-bottom: 4px; }
.pot-amount { font-size: 20px; font-weight: 700; margin-bottom: 6px; }
.winner-row { display: flex; gap: 8px; align-items: center; font-size: 14px; margin-top: 4px; }
.winner-nick { font-weight: 500; }
.win-amount  { color: #27ae60; font-weight: 700; }
.hand-desc   { font-size: 11px; color: var(--text-3); }

.overlay {
  position: fixed; inset: 0; z-index: 500;
  background: rgba(0,0,0,.5); display: flex; align-items: flex-end;
}
.session-end-overlay { align-items: center; justify-content: center; }
.session-end-card {
  background: #fff; border-radius: 16px;
  width: min(360px, 90vw); overflow: hidden;
  box-shadow: 0 8px 32px rgba(0,0,0,.28);
}
.session-end-title {
  text-align: center; font-size: 17px; font-weight: 700;
  padding: 22px 20px 10px; color: #1a1a1a;
}
.session-end-body {
  text-align: center; font-size: 14px; color: #555;
  padding: 0 20px 22px; line-height: 1.8;
}
.session-end-link {
  color: #0EC4B0; font-weight: 600; cursor: pointer; text-decoration: underline;
}
.session-end-actions {
  display: flex; border-top: 1px solid #eee;
}
.session-end-btn {
  flex: 1; padding: 16px 0; font-size: 16px; font-weight: 500;
  background: none; border: none; cursor: pointer;
}
.session-end-btn.stay  { color: #1a1a1a; border-right: 1px solid #eee; }
.session-end-btn.leave { color: #0EC4B0; }
.sheet {
  width: 100%; background: var(--surface);
  border-radius: 20px 20px 0 0;
  padding: 20px 16px calc(28px + env(safe-area-inset-bottom));
  animation: slideUp .3s ease;
}
.sheet-title { font-size: 16px; font-weight: 600; text-align: center; margin-bottom: 16px; }
.wallet-balance-row {
  display: flex; justify-content: space-between; align-items: center;
  background: var(--bg); border-radius: 10px; padding: 10px 14px;
  margin-bottom: 12px;
}
.wallet-label { font-size: 13px; color: var(--text-3); }
.wallet-val   { font-size: 16px; font-weight: 700; color: var(--text-1); }

.range-hint { font-size: 13px; color: var(--text-3); margin-bottom: 10px; }
.range-val  { text-align: center; font-size: 26px; font-weight: 700; margin-top: 8px; }

.invite-code-box {
  display: flex; align-items: center; gap: 12px;
  background: var(--bg); border-radius: 12px; padding: 14px 16px; margin-bottom: 10px;
}
.invite-code-label { font-size: 13px; color: var(--text-3); }
.invite-code-val   { font-size: 22px; font-weight: 700; flex: 1; }
.copy-btn {
  background: #27ae60; color: #fff; border: none;
  border-radius: 14px; height: 30px; padding: 0 14px; font-size: 13px; cursor: pointer;
}
.invite-meta { display: flex; gap: 12px; font-size: 12px; color: var(--text-3); margin-bottom: 10px; }
.invite-tip  { font-size: 12px; color: var(--text-3); text-align: center; line-height: 1.6; }

.rank-meta  { display: flex; gap: 16px; font-size: 12px; color: var(--text-3); margin-bottom: 12px; flex-wrap: wrap; }
.rank-empty { text-align: center; padding: 24px; color: var(--text-3); font-size: 14px; }
.rank-list  { display: flex; flex-direction: column; gap: 10px; max-height: 50vh; overflow-y: auto; }
.rank-row   { display: flex; align-items: center; gap: 10px; }
.rank-no {
  width: 22px; height: 22px; border-radius: 50%; background: var(--bg);
  display: flex; align-items: center; justify-content: center;
  font-size: 11px; font-weight: 600; color: var(--text-3); flex-shrink: 0;
}
.rank-no.top { background: #f5c842; color: #1a1a1a; }
.rank-avatar {
  width: 36px; height: 36px; border-radius: 50%; background: #27ae60;
  display: flex; align-items: center; justify-content: center;
  font-size: 14px; font-weight: 700; color: #fff; flex-shrink: 0; overflow: hidden;
}
.rank-info  { flex: 1; }
.rank-name  { font-size: 14px; font-weight: 500; }
.rank-buyin { font-size: 11px; color: var(--text-3); }
.rank-result { font-size: 16px; font-weight: 700; }
.green { color: #27ae60; }
.red   { color: var(--red); }

/* ── 阶段标签（椭圆外，顶部居中） ── */
.stage-label-outer {
  position: absolute; top: 22%; left: 50%; transform: translateX(-50%);
  background: rgba(0,100,80,.9);
  color: #fff;
  font-size: 12px; font-weight: 700; letter-spacing: 2px;
  padding: 5px 20px; border-radius: 14px;
  border: 1px solid rgba(255,255,255,.25);
  box-shadow: 0 2px 10px rgba(0,0,0,.5);
  pointer-events: none; z-index: 15;
  animation: stageIn .4s cubic-bezier(.34,1.56,.64,1) both;
}
/* 保留旧 class 名（防止其他引用） */
.stage-label { display: none; }
@keyframes stageIn {
  from { transform: translateX(-50%) scale(.7); opacity: 0; }
  to   { transform: translateX(-50%) scale(1);  opacity: 1; }
}

/* ── 位置标识 ── */
.pos-badge {
  display: inline-block; font-size: 9px; font-weight: 700;
  padding: 1px 4px; border-radius: 4px; margin-left: 3px;
  vertical-align: middle; line-height: 1.4;
}
.pos-d  { background: #f5c842; color: #1a1a1a; }
.pos-sb { background: #3498db; color: #fff; }
.pos-bb { background: #e74c3c; color: #fff; }

/* ── 全押快捷 ── */
.preset.allin { background: rgba(231,76,60,.5); color: #fff; }
.preset.allin:active { background: rgba(231,76,60,.8); }

/* ── 牌型参考 ── */
.hands-sheet { max-height: 85vh; overflow-y: auto; }
.hands-table { display: flex; flex-direction: column; gap: 6px; margin-bottom: 12px; }
.hand-row {
  display: flex; align-items: center; gap: 8px;
  padding: 6px 10px; border-radius: 8px;
  background: var(--bg);
}
.hand-rank {
  width: 20px; height: 20px; border-radius: 50%;
  background: #f5c842; color: #1a1a1a;
  font-size: 10px; font-weight: 800;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.hand-name  { font-size: 13px; font-weight: 600; flex: 1; }
.hand-example { font-size: 10px; color: var(--text-3); font-family: monospace; }
.hands-note { font-size: 11px; color: var(--text-3); text-align: center; line-height: 1.6; }

.bot-hint { font-size: 12px; color: var(--text-3); text-align: center; margin: 0 0 14px; }
.bot-count-row { display: flex; gap: 10px; justify-content: center; }
.bot-count-btn {
  width: 44px; height: 44px; border-radius: 50%;
  background: var(--bg); border: 2px solid var(--border);
  font-size: 18px; font-weight: 700; cursor: pointer; color: var(--text-1); transition: all .15s;
}
.bot-count-btn.active { background: #27ae60; border-color: #27ae60; color: #fff; }

.btn-outline {
  background: transparent; border: 1px solid var(--border); color: var(--text-1);
  border-radius: 12px; height: 40px; font-size: 14px; cursor: pointer;
}
.btn-outline:active { background: var(--bg); }


/* ── 动画 ── */
.menu-slide-enter-active, .menu-slide-leave-active { transition: transform .25s ease; }
.menu-slide-enter-from,   .menu-slide-leave-to     { transform: translateX(-100%); }
.fadeIn-enter-active { animation: fadeIn .3s; }
.fadeIn-leave-active { animation: fadeIn .2s reverse; }
.chat-slide-enter-active, .chat-slide-leave-active { transition: all .2s ease; }
.chat-slide-enter-from,   .chat-slide-leave-to     { opacity: 0; transform: translateY(10px); }
.modal-enter-active { animation: slideUp .3s ease; }
.modal-leave-active { animation: slideUp .25s ease reverse; }

@keyframes slideUp { from { transform: translateY(100%); } to { transform: translateY(0); } }
@keyframes fadeIn  { from { opacity: 0; } to { opacity: 1; } }

/* ── 补码提示（筹码为0时替代行动按钮） ── */
.rebuy-prompt {
  display: flex; flex-direction: column; align-items: center; gap: 10px;
  padding: 12px 0;
}
.rebuy-prompt-text {
  font-size: 14px; color: rgba(255,255,255,.85);
  text-shadow: 0 1px 3px rgba(0,0,0,.6);
}
.rebuy-prompt-btn {
  height: 48px; padding: 0 40px; border-radius: 24px;
  background: linear-gradient(135deg, #f5c842, #e6a800);
  border: none; color: #1a1a1a;
  font-size: 16px; font-weight: 800; cursor: pointer;
  box-shadow: 0 4px 16px rgba(245,200,66,.6);
  transition: transform .1s;
}
.rebuy-prompt-btn:active { transform: scale(.96); }

/* 补码弹窗紧急提示 */
.rebuy-urgent-tip {
  background: rgba(255,100,50,.15);
  border: 1px solid rgba(255,120,60,.4);
  border-radius: 10px; padding: 8px 12px;
  font-size: 13px; color: #ffab76; text-align: center;
  margin-bottom: 4px;
}

/* ── 圆形行动按钮 ── */
.circle-btns-row {
  display: flex; align-items: flex-end; justify-content: center;
  gap: 24px; padding-bottom: 4px;
}
.circle-btn-wrap { display: flex; flex-direction: column; align-items: center; }
.circle-main-wrap { margin: 0 8px; }

.circle-btn {
  border: none; cursor: pointer;
  border-radius: 50%;
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  position: relative; overflow: hidden;
  transition: transform .1s, box-shadow .1s;
  color: #fff; font-weight: 700;
}
.circle-btn:active { transform: scale(.93); }

.circle-fold {
  width: 80px; height: 80px;
  background: radial-gradient(circle at 40% 35%, #e57373, #c62828);
  box-shadow: 0 6px 20px rgba(198,40,40,.6), 0 2px 0 rgba(0,0,0,.3);
  font-size: 16px; letter-spacing: .5px;
}
.circle-call {
  width: 100px; height: 100px;
  background: radial-gradient(circle at 40% 35%, #42a5f5, #1565c0);
  box-shadow: 0 8px 28px rgba(21,101,192,.7), 0 2px 0 rgba(0,0,0,.3);
  font-size: 17px; gap: 2px;
}
.circle-raise {
  width: 80px; height: 80px;
  background: radial-gradient(circle at 40% 35%, #66bb6a, #2e7d32);
  box-shadow: 0 6px 20px rgba(46,125,50,.6), 0 2px 0 rgba(0,0,0,.3);
  font-size: 15px;
}
.circle-raise.raise-confirm {
  background: radial-gradient(circle at 40% 35%, #ffd54f, #f57f17);
  box-shadow: 0 6px 20px rgba(245,127,23,.7), 0 2px 0 rgba(0,0,0,.3);
  animation: confirmPulse 0.6s ease-in-out infinite alternate;
}
@keyframes confirmPulse {
  from { transform: scale(1); }
  to   { transform: scale(1.06); }
}

/* 计时扇形覆盖在让牌按钮上 */
.circle-timer-svg {
  position: absolute; inset: 0; width: 100%; height: 100%;
  pointer-events: none;
}
.circle-btn-label { position: relative; z-index: 1; font-size: inherit; }
.circle-btn-sub   { position: relative; z-index: 1; font-size: 12px; opacity: .85; }
.circle-timer-sec {
  position: relative; z-index: 1; font-size: 10px; opacity: .6; margin-top: 1px;
}

/* 计时紧迫（< 6s）红色脉冲 */
.timer-urgent {
  animation: urgentPulse 0.5s ease-in-out infinite alternate;
}
@keyframes urgentPulse {
  from { box-shadow: 0 0 0 0 rgba(255,50,50,.4); }
  to   { box-shadow: 0 0 0 10px rgba(255,50,50,0); background: rgba(180,30,30,.9); }
}

/* 边池标签 */
.side-pot-tag {
  font-size: 10px; color: rgba(255,255,255,.7);
  background: rgba(0,0,0,.4); border-radius: 8px;
  padding: 2px 8px; margin-left: 4px;
}

/* ALL IN 标签脉冲 */
.allin-pulse {
  animation: allinPulse 1s ease-in-out infinite alternate;
}
@keyframes allinPulse {
  from { opacity: 1; transform: scale(1); }
  to   { opacity: .8; transform: scale(1.08); }
}

/* 预选操作行（WePoker风格） */
.pre-action-row {
  display: flex; gap: 8px; justify-content: center;
  margin-top: 6px;
}
.pre-act-btn {
  font-size: 11px; padding: 4px 10px;
  border-radius: 14px; border: 1px solid rgba(255,255,255,.3);
  background: rgba(0,0,0,.3); color: rgba(255,255,255,.7);
  cursor: pointer; transition: all .15s;
}
.pre-act-btn.active {
  background: rgba(46,204,113,.35); border-color: #2ecc71;
  color: #fff; font-weight: 700;
}

/* 旧气泡样式已替换为 .pot-bar，保留空规则防报错 */
.pot-bubbles-area { display: none; }
.pot-bubble, .pot-bubble-main, .pot-bubble-label, .pot-bubble-val { display: none; }

/* ── 桌面中心信息行（跟在公共牌下面） ── */
.table-info-row {
  display: flex; gap: 12px;
  font-size: 11px; color: rgba(255,255,255,.4);
  white-space: nowrap; pointer-events: none;
}
.tinfo-item { display: flex; align-items: center; gap: 3px; }

/* ── YOU WIN 弹窗 ── */
.youwin-overlay {
  position: fixed; inset: 0; z-index: 500;
  background: rgba(0,0,0,.15);
  display: flex; align-items: center; justify-content: center;
  pointer-events: auto;
}
.youwin-content {
  display: flex; flex-direction: column; align-items: center; gap: 8px;
  padding: 0 20px;
}
.youwin-title {
  font-size: 52px; font-weight: 900; font-style: italic;
  color: #f5c842;
  text-shadow: 0 0 30px rgba(245,200,66,.8), 2px 3px 0 #8B6914, -1px -1px 0 #a07820;
  letter-spacing: 2px; line-height: 1;
  animation: youwinPop .5s cubic-bezier(.34,1.56,.64,1) both;
}
.youwin-amount {
  font-size: 28px; font-weight: 800; color: #f5c842;
  text-shadow: 0 0 12px rgba(245,200,66,.6);
  animation: youwinPop .5s .1s cubic-bezier(.34,1.56,.64,1) both;
}
.youwin-nick {
  font-size: 16px; color: rgba(255,255,255,.9); font-weight: 600;
}
.youwin-avatar-wrap {
  width: 88px; height: 88px; border-radius: 50%; overflow: hidden;
  border: 3px solid #f5c842;
  box-shadow: 0 0 20px rgba(245,200,66,.5);
}
.youwin-avatar { width: 100%; height: 100%; object-fit: cover; }
.youwin-avatar-txt {
  width: 100%; height: 100%;
  display: flex; align-items: center; justify-content: center;
  font-size: 32px; font-weight: 700; color: #fff;
  background: rgba(0,0,0,.4);
}
.youwin-cards {
  display: flex; gap: 8px; margin-top: 4px;
}
.youwin-rank {
  font-size: 14px; color: #fff; font-weight: 600;
  background: rgba(0,0,0,.6); border-radius: 12px;
  padding: 3px 14px;
}
.youwin-chips {
  background: rgba(0,0,0,.8); color: #fff; font-weight: 700;
  font-size: 16px; border-radius: 16px; padding: 5px 20px;
  border: 1px solid rgba(255,255,255,.2);
}
@keyframes youwinPop {
  from { transform: scale(.4) translateY(20px); opacity: 0; }
  to   { transform: scale(1)  translateY(0);    opacity: 1; }
}
.youwin-enter-active { animation: fadeIn .3s ease both; }
.youwin-leave-active { transition: opacity .4s; }
.youwin-leave-to    { opacity: 0; }

/* ── 赢筹码浮动动画 ─────────────────────────────────────────── */
.chip-win-float {
  position: absolute;
  top: -36px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 20px;
  font-weight: 900;
  color: #ffe84d;
  text-shadow: 0 0 8px rgba(255,200,0,.9), 0 2px 4px rgba(0,0,0,.8);
  pointer-events: none;
  white-space: nowrap;
  animation: chipFloatUp 1.8s ease-out forwards;
  z-index: 200;
}
@keyframes chipFloatUp {
  0%   { opacity: 1; transform: translateX(-50%) translateY(0) scale(1.2); }
  30%  { opacity: 1; transform: translateX(-50%) translateY(-20px) scale(1.3); }
  100% { opacity: 0; transform: translateX(-50%) translateY(-60px) scale(1); }
}
.chip-float-enter-active { animation: chipFloatUp 1.8s ease-out forwards; }
.chip-float-leave-active { opacity: 0; }
</style>
