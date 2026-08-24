<template>
  <div>
    <h1 class="font-display text-3xl text-paper">夜间台账</h1>
    <p class="mt-1 text-sm text-paper/50">今天的认证、批改与异常都落在这一页。</p>
    <div class="mt-8 grid gap-4 md:grid-cols-3">
      <div v-for="c in cards" :key="c.t" class="card p-5">
        <div class="text-xs uppercase tracking-widest text-inktext/40">{{ c.t }}</div>
        <div class="mt-3 font-display text-3xl">{{ c.v }}</div>
        <p class="mt-2 text-sm text-inktext/60">{{ c.d }}</p>
      </div>
    </div>
  </div>
</template>
<script setup>
import { onMounted, reactive } from 'vue'
import { req } from '../api'
const cards = reactive([
  { t: '待阅卷', v: '—', d: '含简答题的试卷停在这里' },
  { t: '审计标记', v: '—', d: 'R1–R4 命中的会话' },
  { t: '队列深度', v: '—', d: 'Channel Worker 待批改' },
])
onMounted(async () => {
  try {
    const sub = await req('/api/v1/submissions?status=pending_manual&page_size=1')
    const flags = await req('/api/v1/audit-flags')
    const m = await req('/api/v1/grader/metrics')
    cards[0].v = sub.total ?? 0
    cards[1].v = (flags || []).length
    cards[2].v = m.queue_depth ?? 0
  } catch {}
})
</script>
