<script setup>
import { ref, watch } from 'vue'
import { api } from '../api'

const props = defineProps({
  show: { type: Boolean, default: false }
})
const emit = defineEmits(['close', 'changed'])

const types = ref([])
const newName = ref('')
const newDesc = ref('')
const error = ref('')

async function load() {
  types.value = await api.getDeviceTypes()
}
watch(() => props.show, (v) => {
  if (v) { error.value = ''; load() }
})

async function add() {
  if (!newName.value.trim()) return
  try {
    await api.createDeviceType({ name: newName.value.trim(), description: newDesc.value.trim() })
    newName.value = ''
    newDesc.value = ''
    await load()
    emit('changed')
  } catch (e) {
    error.value = e.message
  }
}

async function remove(id) {
  try {
    await api.deleteDeviceType(id)
    await load()
    emit('changed')
  } catch (e) {
    error.value = e.message
  }
}
</script>

<template>
  <div v-if="show" class="modal-overlay" @click.self="emit('close')">
    <div class="modal">
      <div class="modal-head">
        <h2>设备类型管理</h2>
        <button class="modal-close" @click="emit('close')">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      <div class="modal-body">
        <div class="add-form">
          <input v-model="newName" type="text" placeholder="类型名称" @keyup.enter="add" />
          <input v-model="newDesc" type="text" placeholder="描述（可选）" />
          <button class="btn-add" @click="add">添加</button>
        </div>
        <div v-if="error" class="error">{{ error }}</div>
        <div class="type-list">
          <div v-for="t in types" :key="t.id" class="type-item">
            <div class="type-info">
              <span class="type-name">{{ t.name }}</span>
              <span class="type-desc">{{ t.description || '—' }}</span>
            </div>
            <button class="btn-del" @click="remove(t.id)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
            </button>
          </div>
          <div v-if="types.length === 0" class="empty">暂无设备类型</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  z-index: 200;
  display: flex; align-items: center; justify-content: center;
  padding: 20px; animation: fade-in 0.2s;
}
.modal {
  background: var(--bg-surface);
  border: 1px solid var(--border-bright);
  border-radius: var(--radius-lg);
  width: 100%; max-width: 500px;
  animation: scale-in 0.25s;
  box-shadow: 0 24px 60px rgba(0,0,0,0.5);
}
.modal-head {
  display: flex; justify-content: space-between; align-items: center;
  padding: 18px 22px; border-bottom: 1px solid var(--border);
}
.modal-head h2 { font-size: 16px; font-weight: 600; }
.modal-close {
  width: 32px; height: 32px; display: flex; align-items: center; justify-content: center;
  border-radius: 6px; color: var(--text-muted);
}
.modal-close:hover { background: var(--bg-elevated); color: var(--text-primary); }
.modal-body { padding: 22px; display: flex; flex-direction: column; gap: 14px; }
.add-form { display: flex; gap: 8px; }
.add-form input {
  flex: 1; background: var(--bg-base); border: 1px solid var(--border);
  border-radius: var(--radius); padding: 9px 12px; font-size: 13px; outline: none;
}
.add-form input:focus { border-color: var(--accent-dim); }
.btn-add {
  padding: 9px 16px; background: var(--accent); color: #04141a;
  border-radius: var(--radius); font-weight: 600; font-size: 13px; white-space: nowrap;
}
.btn-add:hover { background: #67e8f9; }
.error {
  font-size: 12px; color: var(--abnormal);
  background: var(--abnormal-bg); padding: 8px 12px; border-radius: var(--radius);
}
.type-list { display: flex; flex-direction: column; gap: 6px; max-height: 300px; overflow-y: auto; }
.type-item {
  display: flex; justify-content: space-between; align-items: center;
  padding: 10px 14px; background: var(--bg-base); border: 1px solid var(--border);
  border-radius: var(--radius); transition: border-color 0.15s;
}
.type-item:hover { border-color: var(--border-bright); }
.type-info { display: flex; flex-direction: column; gap: 2px; }
.type-name { font-size: 14px; font-weight: 500; }
.type-desc { font-size: 11px; color: var(--text-muted); }
.btn-del {
  width: 28px; height: 28px; display: flex; align-items: center; justify-content: center;
  border-radius: 5px; color: var(--text-muted);
}
.btn-del:hover { background: var(--abnormal-bg); color: var(--abnormal); }
.empty { text-align: center; color: var(--text-muted); padding: 30px; font-size: 13px; }
</style>
