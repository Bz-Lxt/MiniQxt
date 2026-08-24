<template>
  <div class="card p-6">
    <h1 class="font-display text-2xl">人工阅卷</h1>
    <p class="text-sm text-inktext/50">客观题已由 Grader 批完，简答题在此给分后才会汇总认证。</p>
    <div class="mt-4 space-y-3">
      <article v-for="s in list" :key="s.id" class="rounded-xl bg-black/5 p-4">
        <div class="flex justify-between text-sm">
          <span>提交 #{{ s.id }} · 考生 {{ s.user_id }} · 客观 {{ s.objective_score }}</span>
          <span>{{ s.status }}</span>
        </div>
        <div v-if="s.answers" class="mt-3 space-y-3">
          <div v-for="a in s.answers.filter(x=>x.is_essay)" :key="a.id">
            <div class="text-xs text-inktext/50">题 {{ a.question_id }}</div>
            <p class="mt-1 text-sm">{{ a.answer_text || '（空白）' }}</p>
            <input v-model.number="scores[a.question_id]" type="number" class="field mt-2 !w-32" placeholder="给分" />
          </div>
          <button class="btn-brass" @click="save(s)">写入主观分</button>
        </div>
        <button v-else class="mt-2 text-sm text-brass" @click="open(s)">展开</button>
      </article>
      <p v-if="!list.length" class="text-inktext/40">暂无待阅卷。</p>
    </div>
  </div>
</template>
<script setup>
import { onMounted, reactive, ref } from 'vue'
import { req } from '../api'
import { toast } from '../bus'
const list = ref([]); const scores = reactive({})
async function load() {
  const page = await req('/api/v1/submissions?status=pending_manual')
  list.value = page.items || []
}
async function open(s) {
  const full = await req('/api/v1/submissions/' + s.id)
  Object.assign(s, full.submission)
}
async function save(s) {
  const items = (s.answers || []).filter((x) => x.is_essay).map((a) => ({ question_id: a.question_id, score: Number(scores[a.question_id] || 0), comment: '' }))
  await req('/api/v1/submissions/' + s.id + '/manual-scores', { method: 'PUT', body: JSON.stringify({ items }) })
  toast('已给分'); load()
}
onMounted(load)
</script>
