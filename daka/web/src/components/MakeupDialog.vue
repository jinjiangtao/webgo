<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal">
      <div class="modal-header">
        <h2>📋 补卡申请</h2>
        <button class="close-btn" @click="$emit('close')">✕</button>
      </div>

      <form class="modal-body" @submit.prevent="handleSubmit">
        <div class="task-info">
          <div class="task-color" :style="{ background: task?.color || '#6366f1' }"></div>
          <span class="task-name">{{ task?.name }}</span>
        </div>

        <div class="form-group">
          <label>补卡日期 *</label>
          <input v-model="form.date" type="date" :max="today" required>
        </div>

        <div class="form-group">
          <label>已进行时长</label>
          <div class="duration-inputs">
            <div class="duration-item">
              <input v-model.number="form.hours" type="number" min="0" max="23">
              <span>时</span>
            </div>
            <div class="duration-item">
              <input v-model.number="form.minutes" type="number" min="0" max="59">
              <span>分</span>
            </div>
            <div class="duration-item">
              <input v-model.number="form.seconds" type="number" min="0" max="59">
              <span>秒</span>
            </div>
          </div>
        </div>

        <div class="form-group">
          <label>补卡说明</label>
          <textarea v-model="form.note" rows="3" placeholder="可选：说明补卡原因..."></textarea>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-cancel" @click="$emit('close')">取消</button>
          <button type="submit" class="btn btn-submit">确认补卡</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  task: { type: Object, default: null },
  defaultDuration: { type: Number, default: 0 }
})

const emit = defineEmits(['close', 'submit'])

const today = new Date().toISOString().slice(0, 10)

const form = ref({
  date: today,
  hours: 0,
  minutes: 0,
  seconds: 0,
  note: ''
})

watch(() => props.defaultDuration, (dur) => {
  if (dur) {
    form.value.hours = Math.floor(dur / 3600)
    form.value.minutes = Math.floor((dur % 3600) / 60)
    form.value.seconds = dur % 60
  }
}, { immediate: true })

const totalSeconds = computed(() => {
  return (form.value.hours || 0) * 3600 + (form.value.minutes || 0) * 60 + (form.value.seconds || 0)
})

function handleSubmit() {
  emit('submit', {
    taskId: props.task?.id,
    date: form.value.date,
    duration: totalSeconds.value,
    note: form.value.note
  })
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
  backdrop-filter: blur(4px);
}

.modal {
  background: white;
  border-radius: 24px;
  width: 100%;
  max-width: 460px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.modal-header {
  padding: 24px 28px;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  font-size: 18px;
  font-weight: 700;
}

.close-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #f8fafc;
  font-size: 18px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #f1f5f9;
}

.modal-body {
  padding: 24px 28px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.task-info {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: rgba(139, 92, 246, 0.08);
  border-radius: 12px;
}

.task-color {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.task-name {
  font-weight: 600;
  color: #7c3aed;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-weight: 600;
  font-size: 14px;
  color: #334155;
}

.form-group input,
.form-group textarea {
  padding: 12px 16px;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  font-size: 14px;
  outline: none;
  transition: all 0.2s;
  background: #f8fafc;
}

.form-group input:focus,
.form-group textarea:focus {
  border-color: #8b5cf6;
  background: white;
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.duration-inputs {
  display: flex;
  gap: 12px;
}

.duration-item {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
}

.duration-item input {
  width: 100%;
  padding: 10px 12px;
  border: 2px solid #e2e8f0;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  text-align: center;
  outline: none;
  background: #f8fafc;
}

.duration-item input:focus {
  border-color: #8b5cf6;
  background: white;
}

.duration-item span {
  font-size: 13px;
  color: #64748b;
  font-weight: 500;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding-top: 8px;
  border-top: 1px solid #e2e8f0;
}

.btn {
  padding: 12px 24px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-cancel {
  background: #f1f5f9;
  color: #475569;
}

.btn-cancel:hover {
  background: #e2e8f0;
}

.btn-submit {
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
  color: white;
}

.btn-submit:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
}
</style>
