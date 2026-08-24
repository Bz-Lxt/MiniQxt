<template>
  <div class="grid gap-4 lg:grid-cols-[280px_1fr]">
    <aside class="card p-4">
      <h2 class="font-display text-xl">标记</h2>
      <ul class="mt-3 space-y-2 text-sm">
        <li v-for="f in flags" :key="f.id" class="cursor-pointer rounded-lg bg-black/5 p-2 hover:bg-brass/10" @click="sid = f.session_id; loadTl()">
          <div class="font-semibold">{{ f.rule_code }} {{ f.title }}</div>
          <div class="text-xs text-inktext/50">会话 {{ f.session_id }} · {{ fmt(f.created_at) }}</div>
        </li>
      </ul>
    </aside>
    <section class="card p-5">
      <div class="flex items-center justify-between">
        <h1 class="font-display text-2xl">答题时序回放</h1>
        <button v-if="sid" class="btn-ghost !text-inktext !border-black/10" @click="run">重跑规则</button>
      </div>
      <div v-if="tl" class="mt-6">
        <div class="flex flex-wrap gap-2">
          <span v-for="f in tl.flags" :key="f.id" class="rounded-full bg-danger/10 px-3 py-1 text-xs text-danger">{{ f.rule_code }} {{ f.title }}</span>
        </div>
        <ol class="relative mt-6 border-l border-brass/40 pl-6">
          <li v-for="e in tl.traces" :key="e.seq" class="mb-5">
            <div class="absolute -left-1.5 mt-1 h-3 w-3 rounded-full bg-brass"></div>
            <div class="text-xs text-inktext/40">{{ fmt(e.occurred_at) }} · 第 {{ e.display_no }} 题 · seq {{ e.seq }}</div>
            <div class="text-sm">由「{{ e.from_label || '空' }}」改为「{{ e.to_label || e.answer_text || '空' }}」</div>
          </li>
        </ol>
      </div>
      <p v-else class="mt-8 text-inktext/40">选择左侧标记查看证据链。流水只存 option_id，展示时用 shuffle_seed 还原考生卷面位置。</p>
    </section>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { req, fmt } from '../api'
import { toast } from '../bus'
const flags = ref([]); const tl = ref(null); const sid = ref(0)
const route = useRoute()
async function load() { flags.value = await req('/api/v1/audit-flags') }
async function loadTl() { if (!sid.value) return; tl.value = await req('/api/v1/exam-sessions/' + sid.value + '/timeline') }
async function run() {
  await req('/api/v1/exam-sessions/' + sid.value + '/run-audit', { method: 'POST' })
  toast('规则已重跑'); load(); loadTl()
}
onMounted(async () => {
  await load()
  if (route.params.id) { sid.value = Number(route.params.id); loadTl() }
})
</script>
