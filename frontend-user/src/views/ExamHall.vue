<template>
  <div class="min-h-screen bg-[#070b14] text-paper">
    <header class="flex items-center justify-between border-b border-white/10 px-6 py-3">
      <div>
        <div class="text-xs tracking-[0.3em] text-brass">IMMERSIVE EXAM</div>
        <div class="font-display text-xl">{{ payload?.exam?.title }}</div>
      </div>
      <div class="text-right">
        <div class="font-display text-4xl tabular-nums text-brass">{{ clock }}</div>
        <div class="text-xs text-paper/40">服务端权威剩余时间</div>
      </div>
    </header>
    <main class="mx-auto w-full p-4 md:p-8">
      <article v-for="q in paper" :key="q.question_id" class="mb-6 rounded-2xl border border-white/10 bg-white/[0.03] p-5">
        <div class="text-xs text-brass">第 {{ q.display_no }} 题 · {{ q.type }} · {{ q.score }} 分</div>
        <h3 class="mt-2 text-lg">{{ q.stem }}</h3>
        <div v-if="q.type !== 'essay'" class="mt-4 space-y-2">
          <label v-for="(o, i) in q.options" :key="o.id"
                 class="flex cursor-pointer items-center gap-3 rounded-xl border border-white/10 px-3 py-2 hover:border-brass/50"
                 :class="selected(q, o.id) && 'border-brass bg-brass/10'">
            <input :type="q.type==='multi' ? 'checkbox' : 'radio'" :name="'q'+q.question_id" :checked="selected(q,o.id)" @change="pick(q,o)" />
            <span>{{ letters[i] }}. {{ o.label }}</span>
          </label>
        </div>
        <textarea v-else v-model="essays[q.question_id]" class="field mt-3 min-h-28" @change="essayTrace(q)" placeholder="简答" />
      </article>
      <div class="flex justify-end">
        <button class="btn-brass px-8 py-3" @click="submit(false)">交卷</button>
      </div>
    </main>

    <div v-if="warn" class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4">
      <div class="card w-full max-w-md p-6 text-center">
        <div class="font-display text-2xl text-danger">切屏警告</div>
        <p class="mt-3 text-sm text-inktext/70">系统已记录本次离开页面。累计 {{ blur }} 次后将被标记异常，达到强制阈值将自动交卷。</p>
        <button class="btn-brass mt-6" @click="warn=false">回到考场</button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { req, API } from '../api'
import { confirm, toast } from '../bus'

const letters = 'ABCD'
const route = useRoute(); const router = useRouter()
const payload = ref(null)
const paper = computed(() => payload.value?.paper || [])
const answers = reactive({})
const essays = reactive({})
const remain = ref(0)
const warn = ref(false)
const blur = ref(0)
const traces = []
let seq = 0
let hb = 0
let tick = 0
let flush = 0

onMounted(async () => {
  payload.value = await req('/api/v1/exams/' + route.params.id + '/start')
  remain.value = payload.value.remaining_sec
  document.addEventListener('visibilitychange', onVis)
  window.addEventListener('blur', onBlur)
  window.addEventListener('beforeunload', flushBeacon)
  hb = setInterval(() => req('/api/v1/exam-sessions/' + sid() + '/heartbeat', { method: 'POST' }).catch(() => {}), 15000)
  tick = setInterval(async () => {
    remain.value = Math.max(0, remain.value - 1)
    if (remain.value <= 0) submit(true)
    if (remain.value % 20 === 0) {
      const s = await req('/api/v1/exam-sessions/' + sid())
      remain.value = s.remaining_sec
    }
  }, 1000)
  flush = setInterval(flushTraces, 3000)
})
onUnmounted(() => {
  document.removeEventListener('visibilitychange', onVis)
  window.removeEventListener('blur', onBlur)
  window.removeEventListener('beforeunload', flushBeacon)
  clearInterval(hb); clearInterval(tick); clearInterval(flush)
})
function sid() { return payload.value?.session?.id }
const clock = computed(() => {
  const s = remain.value
  const m = Math.floor(s / 60); const r = s % 60
  return `${String(m).padStart(2,'0')}:${String(r).padStart(2,'0')}`
})
function selected(q, oid) {
  const cur = answers[q.question_id] || []
  return cur.includes(oid)
}
function pick(q, o) {
  const prev = (answers[q.question_id] || [])[0] || null
  if (q.type === 'multi') {
    const cur = new Set(answers[q.question_id] || [])
    cur.has(o.id) ? cur.delete(o.id) : cur.add(o.id)
    answers[q.question_id] = [...cur]
  } else {
    answers[q.question_id] = [o.id]
  }
  traces.push({
    question_id: q.question_id, from_option_id: prev, to_option_id: o.id,
    occurred_at: new Date().toISOString(), seq: ++seq,
  })
}
function essayTrace(q) {
  traces.push({ question_id: q.question_id, answer_text: essays[q.question_id], occurred_at: new Date().toISOString(), seq: ++seq })
}
async function onVis() {
  if (document.hidden) await cheat('hidden')
}
async function onBlur() {
  if (document.hidden) return
  await cheat('blur')
}
async function cheat(type) {
  warn.value = true
  try {
    const r = await req('/api/v1/exam-sessions/' + sid() + '/anti-cheat', {
      method: 'POST', body: JSON.stringify({ event_type: type, occurred_at: new Date().toISOString(), duration_ms: 0 }),
    })
    blur.value = r.blur_count
    if (r.force_submit) submit(true)
  } catch {}
}
async function flushTraces() {
  if (!traces.length || !sid()) return
  const events = traces.splice(0, 80)
  await req('/api/v1/exam-sessions/' + sid() + '/traces', { method: 'POST', body: JSON.stringify({ events }) })
}
function flushBeacon() {
  if (!traces.length || !sid()) return
  const body = JSON.stringify({ events: traces.splice(0) })
  navigator.sendBeacon?.(API + '/api/v1/exam-sessions/' + sid() + '/traces', new Blob([body], { type: 'application/json' }))
}
async function submit(auto) {
  if (!auto && !await confirm('确认交卷', '交卷后客观题进入异步批改队列，请等待出分。')) return
  await flushTraces()
  const ans = paper.value.map((q) => ({
    question_id: q.question_id, option_ids: answers[q.question_id] || [], answer_text: essays[q.question_id] || '',
  }))
  const r = await req('/api/v1/exam-sessions/' + sid() + '/submit', { method: 'POST', body: JSON.stringify({ answers: ans }) })
  toast('已进入批改队列 #' + r.queue_position)
  router.push('/result/' + r.submission_id)
}
</script>
