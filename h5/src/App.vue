<template>
  <router-view v-slot="{ Component, route }">
    <transition :name="route.meta.landscape ? 'none' : 'slide'" mode="out-in">
      <component :is="Component" :key="route.path" />
    </transition>
  </router-view>
</template>

<style>
.slide-enter-active, .slide-leave-active { transition: transform .25s ease, opacity .25s; }
.slide-enter-from { transform: translateX(100%); opacity: 0; }
.slide-leave-to  { transform: translateX(-30%); opacity: 0; }
</style>
