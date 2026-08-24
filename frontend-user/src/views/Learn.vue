<template>
  <div class="min-h-screen bg-ink">
    <header class="flex items-center justify-between px-6 py-4 text-paper">
      <router-link to="/" class="text-brass">← 返回</router-link>
      <div class="font-display">{{ course?.title }}</div>
      <div class="text-sm">有效进度 {{ Math.round(chapter?.percent || 0) }}%</div>
    </header>
    <div class="w-full px-4 pb-10 md:px-10" @mouseleave="onLeave" @mouseenter="onEnter">
      <div class="relative overflow-hidden rounded-3xl bg-black shadow-2xl">
        <video ref="vid" class="w-full" :src="src" @timeupdate="onTime" @play="playing=true" @pause="playing=false"></video>
        <div class="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/80 to-transparent p-4">
          <input type="range" class="w-full" min="0" :max="duration" :value="pos" @input="seek" />
          <div class="mt-2 flex items-center justify-between text-xs text-paper/70">
            <button class="btn-brass" @click="toggle">{{ playing ? '暂停' : '播放' }}</button>
            <span>{{ Math.floor(pos) }} / {{ duration }}s · 拖动进度条不会增加有效学习时长</span>
          </div>
        </div>
        <div v-if="hint" class="absolute left-1/2 top-8 -translate-x-1/2 rounded-full bg-danger px-4 py-2 text-sm text-white">
          检测到鼠标离开，已自动暂停
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { req, API } from '../api'
const route = useRoute()
const course = ref(null)
const chapter = ref(null)
const vid = ref(null)
const playing = ref(false)
const pos = ref(0)
const duration = ref(30)
const hint = ref(false)
const lastTick = ref(0)
const leaveTimer = ref(0)
const src = computed(() => chapter.value ? `${API}/api/v1/media/${chapter.value.video_file}` : '')

onMounted(async () => {
  const list = await req('/api/v1/my/courses')
  course.value = list.find((c) => String(c.id) === String(route.params.id)) || list[0]
  chapter.value = course.value?.chapters?.[0]
  duration.value = chapter.value?.duration_sec || 30
  lastTick.value = performance.now()
})
function toggle() {
  const v = vid.value
  if (!v) return
  v.paused ? v.play() : v.pause()
}
function onTime() {
  const v = vid.value
  if (!v) return
  pos.value = v.currentTime
  if (!v.paused && chapter.value) {
    const now = performance.now()
    const delta = Math.min(5, (now - lastTick.value) / 1000)
    lastTick.value = now
    if (delta >= 1) report(delta)
  }
}
function seek(e) {
  const v = vid.value
  if (!v) return
  v.currentTime = Number(e.target.value)
  lastTick.value = performance.now()
}
function onLeave() {
  leaveTimer.value = setTimeout(() => {
    if (vid.value && !vid.value.paused) {
      vid.value.pause()
      hint.value = true
      setTimeout(() => hint.value = false, 2200)
    }
  }, 400)
}
function onEnter() { clearTimeout(leaveTimer.value) }
async function report(delta) {
  const p = await req('/api/v1/my/chapters/' + chapter.value.id + '/progress', {
    method: 'POST', body: JSON.stringify({ delta_seconds: Math.round(delta), position_sec: Math.floor(pos.value) }),
  })
  chapter.value.percent = p.percent
}
onUnmounted(() => clearTimeout(leaveTimer.value))
</script>
