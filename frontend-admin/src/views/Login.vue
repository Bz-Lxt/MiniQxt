<template>
  <div class="flex min-h-screen items-center justify-center bg-[radial-gradient(circle_at_20%_20%,#1b2a4a,transparent_40%),radial-gradient(circle_at_80%_0%,#3a2a12,transparent_35%)] p-4">
    <div class="card w-full max-w-md p-8">
      <div class="font-display text-4xl text-ink">企学通</div>
      <p class="mt-2 text-sm text-inktext/60">管理台账 · 培训 / 考试 / 认证</p>
      <form class="mt-8 space-y-4" @submit.prevent="onSubmit">
        <div>
          <label class="label">账号 *</label>
          <input v-model="form.username" class="field" />
          <p v-if="err.username" class="mt-1 text-xs text-danger">{{ err.username }}</p>
        </div>
        <div>
          <label class="label">密码 *</label>
          <input v-model="form.password" type="password" class="field" />
          <p v-if="err.password" class="mt-1 text-xs text-danger">{{ err.password }}</p>
        </div>
        <button class="btn-brass w-full" :disabled="loading">{{ loading ? '登录中…' : '进入台账' }}</button>
      </form>
      <p class="mt-6 text-xs text-inktext/50">讲师 teach.zhou@hqtech / Teach@123 · 管理员 admin.hq@hqtech / Tenant@123</p>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'
import { toast } from '../bus'
const router = useRouter()
const auth = useAuth()
const form = reactive({ username: 'teach.zhou@hqtech', password: 'Teach@123' })
const err = reactive({ username: '', password: '' })
const loading = ref(false)
async function onSubmit() {
  err.username = form.username ? '' : '请填写账号'
  err.password = form.password ? '' : '请填写密码'
  if (err.username || err.password) { toast('请先完成必填项', 'err'); return }
  loading.value = true
  try { await auth.login(form.username, form.password); router.push('/') }
  catch (e) { toast(e.message, 'err') }
  finally { loading.value = false }
}
</script>
