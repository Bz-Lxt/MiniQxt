<template>
  <router-view />
  <div class="fixed right-4 top-4 z-[80] flex w-80 flex-col gap-2">
    <div v-for="t in ui.toasts" :key="t.id" class="flex justify-between rounded-xl px-4 py-3 text-sm"
         :class="t.type==='err' ? 'bg-danger text-white' : 'bg-navy text-paper border border-brass/40'">
      <span>{{ t.message }}</span>
      <button @click="dismiss(t.id)">×</button>
    </div>
  </div>
  <div v-if="ui.modal" class="fixed inset-0 z-[70] flex items-center justify-center bg-black/60 p-4">
    <div class="card w-full max-w-md p-6">
      <h3 class="font-display text-xl">{{ ui.modal.title }}</h3>
      <p class="mt-2 text-sm text-inktext/70">{{ ui.modal.body }}</p>
      <div class="mt-6 flex justify-end gap-2">
        <button class="btn-ghost !text-inktext !border-black/10" @click="closeModal(false)">取消</button>
        <button class="btn-brass" @click="closeModal(true)">确认</button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ui, dismiss, closeModal } from './bus'
</script>
