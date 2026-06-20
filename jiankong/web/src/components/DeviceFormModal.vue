<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  show: { type: Boolean, default: false },
  device: { type: Object, default: null },
  deviceTypes: { type: Array, default: () => [] }
})
const emit = defineEmits(['close', 'submit'])

const form = ref({
  name: '',
  device_code: '',
  type_id: '',
  location: '',
  threshold_cpu: 90,
  threshold_temp: 75
})

watch(() => props.show, (val) => {
  if (val) {
    if (props.device) {
      form.value = {
        name: props.device.name,
        device_code: props.device.device_code,
        type_id: props.device.type_id,
        location: props.device.location,
        threshold_cpu: props.device.threshold_cpu,
        threshold_temp: props.device.threshold_temp
      }
    } else {
      form.value = {
        name: '',
        device_code: '',
        type_id: props.deviceTypes[0]?.id || '',
        location: '',
        threshold_cpu: 90,
        threshold_temp: 75
      }
    }
  }
})

function submit() {
  if (!form.value.name.trim() || !form.value.device_code.trim()) return
  emit('submit', { ...form.value })
}
</script>

<template>
  <div v-if="show" class="modal-overlay" @click.self="emit('close')">
    <div class="modal">
      <div class="modal-head">
        <h2>{{ device ? '编辑设备信息' : '录入新设备' }}</h2>
        <button class="modal-close" @click="emit('close')">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      <div class="modal-body">
        <div class="form-row">
          <div class="field">
            <label>设备名称 <span class="req">*</span></label>
            <input v-model="form.name" type="text" placeholder="如：核心数据库服务器" />
          </div>
          <div class="field">
            <label>设备编号 <span class="req">*</span></label>
            <input v-model="form.device_code" type="text" placeholder="如：SRV-001" />
          </div>
        </div>
        <div class="form-row">
          <div class="field">
            <label>设备类型</label>
            <select v-model="form.type_id">
              <option v-for="t in deviceTypes" :key="t.id" :value="t.id">{{ t.name }}</option>
            </select>
          </div>
          <div class="field">
            <label>部署位置</label>
            <input v-model="form.location" type="text" placeholder="如：机房A-01" />
          </div>
        </div>
        <div class="form-row">
          <div class="field">
            <label>CPU 告警阈值 (%)</label>
            <input v-model.number="form.threshold_cpu" type="number" min="1" max="100" />
          </div>
          <div class="field">
            <label>温度告警阈值 (°C)</label>
            <input v-model.number="form.threshold_temp" type="number" min="1" max="120" />
          </div>
        </div>
      </div>
      <div class="modal-foot">
        <button class="btn-cancel" @click="emit('close')">取消</button>
        <button class="btn-save" @click="submit">{{ device ? '保存修改' : '录入设备' }}</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  animation: fade-in 0.2s;
}
.modal {
  background: var(--bg-surface);
  border: 1px solid var(--border-bright);
  border-radius: var(--radius-lg);
  width: 100%;
  max-width: 540px;
  animation: scale-in 0.25s;
  box-shadow: 0 24px 60px rgba(0,0,0,0.5);
}
.modal-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 22px;
  border-bottom: 1px solid var(--border);
}
.modal-head h2 {
  font-size: 16px;
  font-weight: 600;
}
.modal-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  color: var(--text-muted);
}
.modal-close:hover { background: var(--bg-elevated); color: var(--text-primary); }
.modal-body {
  padding: 22px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.form-row {
  display: flex;
  gap: 16px;
}
.field {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.field label {
  font-size: 12px;
  color: var(--text-secondary);
}
.req { color: var(--abnormal); }
.field input, .field select {
  background: var(--bg-base);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 9px 12px;
  font-size: 13px;
  outline: none;
}
.field input:focus, .field select:focus { border-color: var(--accent-dim); }
.modal-foot {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 16px 22px;
  border-top: 1px solid var(--border);
}
.btn-cancel, .btn-save {
  padding: 9px 18px;
  border-radius: var(--radius);
  font-size: 13px;
  font-weight: 500;
}
.btn-cancel {
  background: var(--bg-elevated);
  color: var(--text-secondary);
  border: 1px solid var(--border);
}
.btn-cancel:hover { color: var(--text-primary); }
.btn-save {
  background: var(--accent);
  color: #04141a;
  font-weight: 600;
}
.btn-save:hover { background: #67e8f9; box-shadow: 0 0 14px var(--accent-glow); }
</style>
