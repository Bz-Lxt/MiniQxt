<template>
  <div class="card p-6">
    <h1 class="font-display text-2xl">批改队列</h1>
    <p class="text-sm text-inktext/50">Channel 只是调度加速层，submission.status 才是持久化真相。</p>
    <dl class="mt-6 grid gap-4 md:grid-cols-3">
      <div v-for="k in keys" :key="k" class="rounded-xl bg-black/5 p-4">
        <dt class="text-xs uppercase tracking-widest text-inktext/40">{{ k }}</dt>
        <dd class="mt-2 font-display text-3xl">{{ m[k] ?? '—' }}</dd>
      </div>
    </dl>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { req } from '../api'
const m = ref({})
const keys = ['queue_depth', 'workers', 'enqueued', 'graded_total', 'failed', 'recovered', 'trace_rows', 'sql_inserts', 'write_amplification']
onMounted(async () => { m.value = await req('/api/v1/grader/metrics') })
</script>
