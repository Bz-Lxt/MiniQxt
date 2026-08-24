<template>
  <div class="card p-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="font-display text-2xl">题库</h1>
      <button class="btn-brass" @click="open = true">新建题目</button>
    </div>
    <div class="mt-4 flex flex-wrap gap-2">
      <select v-model="typ" class="field !w-40" @change="load">
        <option value="">全部题型</option>
        <option value="single">单选</option>
        <option value="multi">多选</option>
        <option value="judge">判断</option>
        <option value="essay">简答</option>
      </select>
      <input v-model="q" class="field !max-w-xs" placeholder="搜索题干" @keyup.enter="load" />
    </div>
    <div class="mt-4 overflow-x-auto">
      <table class="w-full min-w-[720px] text-left text-sm">
        <thead class="text-inktext/50"><tr><th class="py-2">题干</th><th>类型</th><th>难度</th><th>分</th><th></th></tr></thead>
        <tbody>
          <tr v-for="item in items" :key="item.id" class="border-t border-black/5">
            <td class="py-2 pr-4">{{ item.stem }}</td>
            <td>{{ item.type }}</td>
            <td>{{ item.difficulty }}</td>
            <td>{{ item.score }}</td>
            <td><button class="text-danger" @click="del(item)">删除</button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="open" class="fixed inset-0 z-40 flex items-center justify-center bg-black/50 p-4">
      <form class="card w-full max-w-xl p-6" @submit.prevent="save">
        <h3 class="font-display text-xl">新题</h3>
        <label class="label mt-4">题干 *</label>
        <textarea v-model="form.stem" class="field min-h-24" />
        <p v-if="ferr.stem" class="text-xs text-danger">{{ ferr.stem }}</p>
        <div class="mt-3 grid grid-cols-3 gap-2">
          <select v-model="form.type" class="field">
            <option value="single">单选</option><option value="multi">多选</option>
            <option value="judge">判断</option><option value="essay">简答</option>
          </select>
          <input v-model.number="form.difficulty" type="number" min="1" max="5" class="field" />
          <input v-model.number="form.score" type="number" min="1" class="field" />
        </div>
        <div v-if="form.type !== 'essay'" class="mt-3 space-y-2">
          <div v-for="(o,i) in form.options" :key="i" class="flex gap-2">
            <input v-model="o.label" class="field" :placeholder="'选项 '+(i+1)" />
            <label class="flex items-center gap-1 text-xs"><input v-model="o.is_correct" type="checkbox" />正确</label>
          </div>
          <button type="button" class="btn-ghost !text-inktext !border-black/10" @click="form.options.push({label:'',is_correct:false})">加选项</button>
        </div>
        <div class="mt-5 flex justify-end gap-2">
          <button type="button" class="btn-ghost !text-inktext !border-black/10" @click="open=false">取消</button>
          <button class="btn-brass">保存</button>
        </div>
      </form>
    </div>
  </div>
</template>
<script setup>
import { onMounted, reactive, ref } from 'vue'
import { req } from '../api'
import { confirm, toast } from '../bus'
const items = ref([]); const typ = ref(''); const q = ref(''); const open = ref(false)
const form = reactive({ stem: '', type: 'single', difficulty: 1, score: 5, options: [{label:'',is_correct:true},{label:'',is_correct:false}] })
const ferr = reactive({ stem: '' })
async function load() {
  const page = await req(`/api/v1/questions?type=${typ.value}&q=${encodeURIComponent(q.value)}&page_size=50`)
  items.value = page.items || []
}
async function save() {
  ferr.stem = form.stem.trim() ? '' : '题干必填'
  if (ferr.stem) { toast('请填写题干', 'err'); return }
  await req('/api/v1/questions', { method: 'POST', body: JSON.stringify(form) })
  open.value = false; toast('题目已入库'); load()
}
async function del(item) {
  if (!await confirm('删除题目', '删除后组卷将无法再引用该题。')) return
  await req('/api/v1/questions/' + item.id, { method: 'DELETE' })
  toast('已删除'); load()
}
onMounted(load)
</script>
