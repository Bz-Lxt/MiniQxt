<template>
  <div class="card p-6">
    <h1 class="font-display text-2xl">考试安排</h1>
    <div class="mt-4 overflow-x-auto">
      <table class="w-full min-w-[720px] text-left text-sm">
        <thead class="text-inktext/50"><tr><th class="py-2">名称</th><th>开始</th><th>结束</th><th>时长</th><th>及格</th></tr></thead>
        <tbody>
          <tr v-for="e in exams" :key="e.id" class="border-t border-black/5">
            <td class="py-2">{{ e.title }}</td>
            <td>{{ fmt(e.start_at) }}</td>
            <td>{{ fmt(e.end_at) }}</td>
            <td>{{ e.duration_sec / 60 }} 分钟</td>
            <td>{{ e.pass_score }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { req, fmt } from '../api'
const exams = ref([])
onMounted(async () => { exams.value = await req('/api/v1/exams') })
</script>
