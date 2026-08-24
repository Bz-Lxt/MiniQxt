<template>
  <div class="grid gap-4 lg:grid-cols-[340px_1fr]">
    <aside class="card p-4">
      <h2 class="font-display text-xl">题库抽屉</h2>
      <select v-model="typ" class="field mt-3" @change="loadQ">
        <option value="">全部</option>
        <option value="single">单选</option>
        <option value="multi">多选</option>
        <option value="judge">判断</option>
        <option value="essay">简答</option>
      </select>
      <div class="mt-3 max-h-[70vh] space-y-2 overflow-y-auto">
        <article v-for="q in bank" :key="q.id" draggable="true"
          class="cursor-grab rounded-xl border border-black/5 bg-white p-3 text-sm hover:border-brass"
          @dragstart="onDrag($event, q)">
          <div class="text-[10px] uppercase tracking-wider text-brass">{{ q.type }} · {{ q.score }}分</div>
          <div class="mt-1">{{ q.stem }}</div>
        </article>
      </div>
    </aside>
    <section class="card p-5">
      <div class="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 class="font-display text-2xl">拖拽智能出题看板</h1>
          <p class="text-sm text-inktext/50">把题目拖进网格。乱序只改变考生看到的位置，审计仍记原始 option_id。</p>
        </div>
        <select v-model.number="paperId" class="field !w-64" @change="loadPaper">
          <option :value="0">选择试卷</option>
          <option v-for="p in papers" :key="p.id" :value="p.id">{{ p.title }}</option>
        </select>
      </div>
      <div v-if="paperId" class="mt-4 flex flex-wrap items-center gap-4 text-sm">
        <label class="flex items-center gap-2"><input v-model="shuffleQ" type="checkbox" />题目乱序</label>
        <label class="flex items-center gap-2"><input v-model="shuffleO" type="checkbox" />选项乱序</label>
        <button class="btn-brass" @click="save">保存组卷</button>
      </div>
      <div class="mt-4 min-h-[360px] rounded-2xl border-2 border-dashed border-brass/40 bg-[#faf6ec] p-4"
           @dragover.prevent @drop="onDrop">
        <div v-if="!items.length" class="flex h-64 items-center justify-center text-inktext/40">从左侧拖入单选 / 多选 / 判断 / 简答</div>
        <div class="grid gap-3 md:grid-cols-2">
          <div v-for="(it, i) in items" :key="it.question_id + '-' + i"
               class="rounded-xl bg-white p-3 shadow-sm">
            <div class="flex items-start justify-between gap-2">
              <div>
                <div class="text-[10px] text-brass">#{{ i + 1 }} {{ it.question?.type || it.type }} · {{ it.score }}分</div>
                <div class="mt-1 text-sm">{{ it.question?.stem || it.stem }}</div>
              </div>
              <button class="text-danger" @click="items.splice(i,1)">移除</button>
            </div>
            <input v-model.number="it.score" type="number" class="field mt-2 !w-24" />
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { req } from '../api'
import { toast } from '../bus'
const bank = ref([]); const papers = ref([]); const items = ref([])
const typ = ref(''); const paperId = ref(0)
const shuffleQ = ref(true); const shuffleO = ref(true)
async function loadQ() {
  const page = await req(`/api/v1/questions?type=${typ.value}&page_size=100`)
  bank.value = page.items || []
}
async function loadPapers() { papers.value = await req('/api/v1/papers') }
async function loadPaper() {
  if (!paperId.value) return
  const p = await req('/api/v1/papers/' + paperId.value)
  shuffleQ.value = p.shuffle_questions; shuffleO.value = p.shuffle_options
  items.value = (p.items || []).map((it) => ({ ...it }))
}
function onDrag(e, q) { e.dataTransfer.setData('text/plain', JSON.stringify(q)) }
function onDrop(e) {
  const q = JSON.parse(e.dataTransfer.getData('text/plain'))
  if (items.value.some((x) => x.question_id === q.id)) { toast('该题已在卷中', 'err'); return }
  items.value.push({ question_id: q.id, score: q.score, group_name: q.type === 'essay' ? '主观题' : '客观题', question: q, stem: q.stem, type: q.type })
}
async function save() {
  if (!paperId.value) return
  await req('/api/v1/papers/' + paperId.value + '/items', {
    method: 'PUT',
    body: JSON.stringify({
      shuffle_questions: shuffleQ.value, shuffle_options: shuffleO.value,
      items: items.value.map((it) => ({ question_id: it.question_id, score: Number(it.score), group_name: it.group_name })),
    }),
  })
  toast('试卷已保存'); loadPaper()
}
onMounted(async () => { await loadQ(); await loadPapers(); if (papers.value[0]) { paperId.value = papers.value[0].id; loadPaper() } })
</script>
