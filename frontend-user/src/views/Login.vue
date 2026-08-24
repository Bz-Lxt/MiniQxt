<template>
  <div class="flex min-h-screen items-center justify-center p-4">
    <div class="card w-full max-w-md p-8">
      <div class="font-display text-4xl">企学通</div>
      <p class="mt-2 text-sm text-inktext/60">学员工作台 · 学习 / 考试 / 证书</p>
      <form class="mt-8 space-y-4" @submit.prevent="go">
        <div>
          <label class="label">账号 *</label>
          <input v-model="u" class="field" />
          <p v-if="eu" class="text-xs text-danger">{{ eu }}</p>
        </div>
        <div>
          <label class="label">密码 *</label>
          <input v-model="p" type="password" class="field" />
          <p v-if="ep" class="text-xs text-danger">{{ ep }}</p>
        </div>
        <button class="btn-brass w-full">进入学习</button>
      </form>
      <p class="mt-6 text-xs text-inktext/50">emp.li@hqtech / Emp@123</p>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'
import { toast } from '../bus'
const u = ref('emp.li@hqtech'); const p = ref('Emp@123')
const eu = ref(''); const ep = ref('')
const auth = useAuth(); const router = useRouter()
async function go() {
  eu.value = u.value ? '' : '请填写账号'; ep.value = p.value ? '' : '请填写密码'
  if (eu.value || ep.value) { toast('请完成必填项', 'err'); return }
  try { await auth.login(u.value, p.value); router.push('/') } catch (e) { toast(e.message, 'err') }
}
</script>
