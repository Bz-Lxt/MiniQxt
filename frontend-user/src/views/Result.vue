<template>
  <div class="flex min-h-screen items-center justify-center p-4">
    <div class="card w-full max-w-lg p-8 text-center">
      <div class="text-xs tracking-[0.3em] text-inktext/40">GRADER QUEUE</div>
      <h1 class="mt-2 font-display text-3xl">{{ title }}</h1>
      <p class="mt-3 text-sm text-inktext/60">{{ desc }}</p>
      <div class="mt-6 font-display text-5xl text-brass">{{ score }}</div>
      <div class="mt-2 text-xs">队列位置 {{ pos }} · 完整性 {{ integrity }}</div>
      <router-link to="/" class="btn-brass mt-8 inline-flex">返回工作台</router-link>
    </div>
  </div>
</template>
<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { req } from '../api'
const route = useRoute()
const sub = ref(null)
const pos = ref(0)
let t = 0
const title = computed(() => {
  const s = sub.value?.status
  if (s === 'queued' || s === 'grading') return '正在批改客观题'
  if (s === 'pending_manual') return '等待讲师阅卷'
  if (s === 'graded') return sub.value.pass ? '已通过' : '未通过'
  return '交卷状态'
})
const desc = computed(() => {
  const s = sub.value?.status
  if (s === 'queued') return `前面还有 ${pos.value} 份试卷，Channel Worker 正在消化峰值。`
  if (s === 'grading') return '评分器正在对照 option_id 判定。'
  if (s === 'pending_manual') return '客观题已出分，简答题进入人工队列。'
  return '成绩已写入台账。'
})
const score = computed(() => sub.value?.total_score ?? '…')
const integrity = computed(() => sub.value?.integrity ?? '—')
async function poll() {
  const r = await req('/api/v1/submissions/' + route.params.id)
  sub.value = r.submission
  pos.value = r.queue_position
}
onMounted(() => { poll(); t = setInterval(poll, 1500) })
onUnmounted(() => clearInterval(t))
</script>
