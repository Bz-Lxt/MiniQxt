<template>
  <div class="card p-6">
    <h1 class="font-display text-2xl">岗位认证</h1>
    <div class="mt-4 space-y-3">
      <div v-for="p in programs" :key="p.id" class="flex items-center justify-between rounded-xl bg-black/5 p-4">
        <div>
          <div class="font-semibold">{{ p.name }}</div>
          <div class="text-xs text-inktext/50">及格 {{ p.min_score }} · 有效 {{ p.valid_days }} 天</div>
        </div>
        <button class="btn-brass" @click="evalP(p)">立即判定</button>
      </div>
    </div>
    <h2 class="mt-8 font-display text-xl">已签发证书</h2>
    <div class="mt-3 overflow-x-auto">
      <table class="w-full min-w-[640px] text-left text-sm">
        <thead class="text-inktext/50"><tr><th class="py-2">编号</th><th>学员</th><th>签发</th><th>到期</th></tr></thead>
        <tbody>
          <tr v-for="c in certs" :key="c.id" class="border-t border-black/5">
            <td class="py-2 font-display">{{ c.no }}</td>
            <td>{{ c.user_id }}</td>
            <td>{{ fmt(c.issued_at) }}</td>
            <td>{{ fmt(c.expire_at) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { req, fmt } from '../api'
import { toast } from '../bus'
const programs = ref([]); const certs = ref([])
async function load() {
  programs.value = await req('/api/v1/cert-programs')
  certs.value = await req('/api/v1/certificates')
}
async function evalP(p) {
  const r = await req('/api/v1/cert-programs/' + p.id + '/evaluate', { method: 'POST' })
  toast('新签发 ' + r.issued + ' 张'); load()
}
onMounted(load)
</script>
