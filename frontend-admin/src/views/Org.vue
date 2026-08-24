<template>
  <div class="card p-6">
    <h1 class="font-display text-2xl">组织</h1>
    <div class="mt-6 grid gap-6 lg:grid-cols-2">
      <section>
        <h2 class="text-sm font-semibold">部门</h2>
        <form class="mt-3 flex gap-2" @submit.prevent="addDept">
          <input v-model="deptName" class="field" placeholder="新部门" />
          <button class="btn-brass">添加</button>
        </form>
        <p v-if="e.dept" class="mt-1 text-xs text-danger">{{ e.dept }}</p>
        <ul class="mt-3 space-y-2 text-sm">
          <li v-for="d in depts" :key="d.id" class="flex justify-between rounded-lg bg-black/5 px-3 py-2">{{ d.name }}</li>
        </ul>
      </section>
      <section>
        <h2 class="text-sm font-semibold">岗位</h2>
        <form class="mt-3 flex gap-2" @submit.prevent="addPos">
          <input v-model="posName" class="field" placeholder="岗位名" />
          <input v-model.number="posLevel" type="number" min="1" class="field !w-24" />
          <button class="btn-brass">添加</button>
        </form>
        <ul class="mt-3 space-y-2 text-sm">
          <li v-for="p in positions" :key="p.id" class="flex justify-between rounded-lg bg-black/5 px-3 py-2">{{ p.name }} · L{{ p.level }}</li>
        </ul>
      </section>
    </div>
    <section class="mt-8">
      <h2 class="text-sm font-semibold">员工</h2>
      <div class="mt-3 overflow-x-auto">
        <table class="w-full min-w-[640px] text-left text-sm">
          <thead class="text-inktext/50"><tr><th class="py-2">姓名</th><th>账号</th><th>角色</th></tr></thead>
          <tbody>
            <tr v-for="u in emps" :key="u.id" class="border-t border-black/5">
              <td class="py-2">{{ u.display_name }}</td><td>{{ u.username }}</td><td>{{ u.role }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { req } from '../api'
import { toast } from '../bus'
const depts = ref([]); const positions = ref([]); const emps = ref([])
const deptName = ref(''); const posName = ref(''); const posLevel = ref(1)
const e = ref({ dept: '' })
async function load() {
  depts.value = await req('/api/v1/departments')
  positions.value = await req('/api/v1/positions')
  const page = await req('/api/v1/employees?page_size=50')
  emps.value = page.items || []
}
async function addDept() {
  e.value.dept = deptName.value ? '' : '部门名称必填'
  if (e.value.dept) { toast('请填写部门', 'err'); return }
  await req('/api/v1/departments', { method: 'POST', body: JSON.stringify({ name: deptName.value }) })
  deptName.value = ''; toast('已添加部门'); load()
}
async function addPos() {
  if (!posName.value) { toast('岗位名称必填', 'err'); return }
  await req('/api/v1/positions', { method: 'POST', body: JSON.stringify({ name: posName.value, level: posLevel.value }) })
  posName.value = ''; toast('已添加岗位'); load()
}
onMounted(load)
</script>
