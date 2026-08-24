<template>
  <div class="flex min-h-screen">
    <aside class="hidden w-64 shrink-0 flex-col border-r border-white/10 bg-navy md:flex">
      <div class="px-6 py-7">
        <div class="font-display text-2xl text-brass">企学通</div>
        <div class="mt-1 text-xs tracking-[0.2em] text-paper/40">MIDNIGHT LEDGER</div>
      </div>
      <nav class="flex-1 space-y-1 px-3">
        <router-link v-for="l in links" :key="l.to" :to="l.to"
          class="block rounded-lg px-3 py-2 text-sm text-paper/70 hover:bg-white/5 hover:text-brass"
          active-class="!bg-brass/15 !text-brass">{{ l.label }}</router-link>
      </nav>
      <div class="border-t border-white/10 p-4 text-xs text-paper/50">
        <div>{{ auth.user?.display_name }}</div>
        <div class="mt-1">{{ auth.user?.role }}</div>
        <button class="mt-3 text-brass hover:underline" @click="out">退出</button>
      </div>
    </aside>
    <div class="flex min-w-0 flex-1 flex-col">
      <header class="flex items-center justify-between border-b border-white/10 px-4 py-3 md:hidden">
        <div class="font-display text-lg text-brass">企学通</div>
        <select class="field !bg-navy !text-paper !border-white/10" :value="$route.path" @change="$router.push($event.target.value)">
          <option v-for="l in links" :key="l.to" :value="l.to">{{ l.label }}</option>
        </select>
      </header>
      <main class="w-full flex-1 p-4 md:p-8">
        <router-view />
      </main>
    </div>
  </div>
</template>
<script setup>
import { useAuth } from '../stores/auth'
import { useRouter } from 'vue-router'
const auth = useAuth()
const router = useRouter()
const links = [
  { to: '/', label: '总览' },
  { to: '/org', label: '组织' },
  { to: '/questions', label: '题库' },
  { to: '/papers', label: '出题看板' },
  { to: '/exams', label: '考试' },
  { to: '/grade', label: '人工阅卷' },
  { to: '/audit', label: '时序审计' },
  { to: '/courses', label: '课程' },
  { to: '/certs', label: '认证' },
  { to: '/queue', label: '批改队列' },
]
function out() { auth.logout(); router.push('/login') }
</script>
