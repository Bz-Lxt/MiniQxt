<template>
  <div class="card p-6">
    <h1 class="font-display text-2xl">课程</h1>
    <div class="mt-4 grid gap-4 md:grid-cols-2">
      <article v-for="c in list" :key="c.id" class="rounded-xl bg-black/5 p-4">
        <div class="font-semibold">{{ c.title }}</div>
        <p class="mt-1 text-sm text-inktext/60">{{ c.summary }}</p>
        <ul class="mt-3 text-xs text-inktext/50">
          <li v-for="ch in c.chapters" :key="ch.id">{{ ch.title }} · {{ ch.duration_sec }}s</li>
        </ul>
      </article>
    </div>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { req } from '../api'
const list = ref([])
onMounted(async () => { list.value = await req('/api/v1/courses') })
</script>
