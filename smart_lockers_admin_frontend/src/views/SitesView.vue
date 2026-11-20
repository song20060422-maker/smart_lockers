<template>
  <div>
    <div class="page-head">
      <h2>网点管理</h2>
      <button v-if="auth.can('site:write')" class="btn primary" @click="showCreate = true">+ 新建网点</button>
    </div>
    <div class="card">
      <table>
        <thead><tr><th>ID</th><th>编码</th><th>名称</th><th>城市</th><th>地址</th><th>负责人</th><th>柜门</th><th>故障</th><th>状态</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="s in rows" :key="s.id">
            <td>{{ s.id }}</td><td>{{ s.code }}</td><td>{{ s.name }}</td>
            <td>{{ s.city }}</td><td>{{ s.address }}</td><td>{{ s.manager_phone }}</td>
            <td>{{ s.door_total }}</td><td>{{ s.door_fault }}</td>
            <td><span class="tag" :class="s.status === 1 ? 'ok' : 'err'">{{ s.status === 1 ? '启用' : '停用' }}</span></td>
            <td>
              <button v-if="auth.can('site:write')" class="btn" @click="toggleStatus(s)">
                {{ s.status === 1 ? '停用' : '启用' }}
              </button>
              <span v-else class="muted">—</span>
            </td>
          </tr>
          <tr v-if="!rows.length"><td colspan="10" class="muted">暂无数据</td></tr>
        </tbody>
      </table>
    </div>

    <div v-if="showCreate" class="modal">
      <form class="card" @submit.prevent="onCreate">
        <h3>新建网点</h3>
        <label>编码<input v-model="form.code" required /></label>
        <label>名称<input v-model="form.name" required /></label>
        <label>城市<input v-model="form.city" /></label>
        <label>地址<input v-model="form.address" /></label>
        <label>负责人电话<input v-model="form.manager_phone" /></label>
        <div class="actions">
          <button class="btn primary" type="submit">提交</button>
          <button class="btn" type="button" @click="showCreate = false">取消</button>
        </div>
      </form>
    </div>
    <div v-if="message" class="toast">{{ message }}</div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'
import type { Site } from '@/types/api'

const auth = useAuthStore()
const { message, notify } = useNotify()
const rows = ref<Site[]>([])
const showCreate = ref(false)
const form = ref({ code: '', name: '', city: '', address: '', manager_phone: '' })
async function load() {
  try { rows.value = (await api.listSites()).data || [] } catch (e: any) { notify(e.message) }
}
async function toggleStatus(s: Site) {
  try { await api.updateSiteStatus(s.id, s.status === 1 ? 0 : 1); notify('操作成功'); await load() }
  catch (e: any) { notify(e.message) }
}
async function onCreate() {
  try { await api.createSite({ ...form.value }); notify('创建成功'); showCreate.value = false; form.value = { code: '', name: '', city: '', address: '', manager_phone: '' }; await load() }
  catch (e: any) { notify(e.message) }
}
onMounted(load)
</script>
<style scoped>
.muted { color: var(--muted); }
.modal { position: fixed; inset: 0; background: #0006; display: flex; align-items: center; justify-content: center; z-index: 50; }
.modal form { width: 380px; display: flex; flex-direction: column; gap: 12px; }
.modal label { display: flex; flex-direction: column; gap: 5px; font-size: 13px; color: var(--muted); }
.actions { display: flex; gap: 10px; justify-content: flex-end; }
</style>
