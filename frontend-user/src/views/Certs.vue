<template>
  <div class="min-h-screen p-6">
    <router-link to="/" class="text-brass">← 返回</router-link>
    <h1 class="mt-4 font-display text-3xl text-paper">我的证书</h1>
    <div class="mt-6 grid gap-4 md:grid-cols-2">
      <article v-for="c in list" :key="c.id" class="relative overflow-hidden rounded-2xl border border-brass/40 bg-gradient-to-br from-[#1b2437] to-[#0b1220] p-6 text-paper">
        <div class="font-display text-4xl text-brass">{{ c.no }}</div>
        <div class="mt-3">{{ c.program?.name }}</div>
        <div class="mt-2 text-xs text-paper/50">签发 {{ fmt(c.issued_at) }} · 到期 {{ fmt(c.expire_at) }}</div>
        <div class="absolute -right-6 -top-6 font-display text-8xl text-brass/10">QXT</div>
      </article>
      <p v-if="!list.length" class="text-paper/40">通过考试并完成必修课后将自动签发。</p>
    </div>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { req, fmt } from '../api'
const list = ref([])
onMounted(async () => { list.value = await req('/api/v1/my/certificates') })
</script>
