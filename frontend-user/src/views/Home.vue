<template>
  <div class="min-h-screen">
    <header class="flex items-center justify-between border-b border-white/10 px-6 py-4">
      <div class="font-display text-2xl text-brass">企学通</div>
      <div class="flex items-center gap-4 text-sm">
        <router-link to="/certs" class="hover:text-brass">我的证书</router-link>
        <button class="text-paper/50" @click="out">退出</button>
      </div>
    </header>
    <main class="w-full p-6">
      <h1 class="font-display text-3xl">你好，{{ auth.user?.display_name }}</h1>
      <section class="mt-8">
        <h2 class="text-sm tracking-widest text-paper/40">待学课程</h2>
        <div class="mt-3 grid gap-4 md:grid-cols-2">
          <article v-for="c in courses" :key="c.id" class="card p-5">
            <div class="font-semibold">{{ c.title }}</div>
            <p class="mt-1 text-sm text-inktext/60">{{ c.summary }}</p>
            <div class="mt-3 h-2 rounded-full bg-black/10">
              <div class="h-2 rounded-full bg-brass" :style="{width: (c.percent||0)+'%'}"></div>
            </div>
            <router-link class="btn-brass mt-4 inline-flex" :to="'/learn/'+c.id">进入工作台</router-link>
          </article>
        </div>
      </section>
      <section class="mt-10">
        <h2 class="text-sm tracking-widest text-paper/40">待考认证</h2>
        <div class="mt-3 grid gap-4 md:grid-cols-2">
          <article v-for="e in exams" :key="e.id" class="card p-5">
            <div class="font-semibold">{{ e.title }}</div>
            <div class="mt-1 text-xs text-inktext/50">{{ fmt(e.start_at) }} — {{ fmt(e.end_at) }}</div>
            <router-link class="btn-brass mt-4 inline-flex" :to="'/exam/'+e.id">进入考场</router-link>
          </article>
        </div>
      </section>
    </main>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'
import { req, fmt } from '../api'
const auth = useAuth(); const router = useRouter()
const courses = ref([]); const exams = ref([])
onMounted(async () => {
  courses.value = await req('/api/v1/my/courses')
  exams.value = await req('/api/v1/my/exams')
})
function out() { auth.logout(); router.push('/login') }
</script>
