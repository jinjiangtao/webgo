<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal">
      <div class="modal-header">
        <h2>{{ isEdit ? '编辑任务' : '创建任务' }}</h2>
        <button class="close-btn" @click="$emit('close')">✕</button>
      </div>

      <form class="modal-body" @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>任务名称 *</label>
          <input v-model="form.name" type="text" placeholder="例如：每日阅读" required>
        </div>

        <div class="form-group">
          <label>任务描述</label>
          <textarea v-model="form.description" rows="3" placeholder="描述一下这个打卡任务的目标..."></textarea>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>打卡周期 *</label>
            <select v-model="form.cycle_type" required>
              <option value="daily">每日</option>
              <option value="weekly">每周</option>
              <option value="monthly">每月</option>
              <option value="custom">自定义</option>
            </select>
          </div>

          <div class="form-group">
            <label>主题颜色</label>
            <div class="color-picker">
              <div
                v-for="color in colors"
                :key="color"
                class="color-option"
                :class="{ active: form.color === color }"
                :style="{ background: color }"
                @click="form.color = color"
              ></div>
              <input v-model="form.color" type="color" class="color-input" title="自定义颜色">
            </div>
          </div>
        </div>

        <div class="form-group">
          <label>倒计时时长 *</label>
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
          <p class="hint">总时长: {{ formatDuration(totalSeconds) }}</p>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>开始时间</label>
            <input v-model="form.start_time" type="time" placeholder="可选">
          </div>
          <div class="form-group">
            <label>结束时间</label>
            <input v-model="form.end_time" type="time" placeholder="可选">
          </div>
        </div>

        <div class="form-group checkbox-group">
          <label>
            <input v-model="form.is_active" type="checkbox">
            <span>启用任务</span>
          </label>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-cancel" @click="$emit('close')">取消</button>
          <button type="submit" class="btn btn-submit">{{ isEdit ? '保存' : '创建' }}</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { formatDuration } from '../utils'

const props = defineProps({
  task: { type: Object, default: null }
})

const emit = defineEmits(['close', 'submit'])

const colors = [
  '#6366f1', '#8b5cf6', '#ec4899', '#ef4444',
  '#f59e0b', '#10b981', '#14b8a6', '#3b82f6'
]

const isEdit = computed(() => !!props.task)

const form = ref({
  name: '',
  description: '',
  cycle_type: 'daily',
  color: '#6366f1',
  hours: 0,
  minutes: 30,
  seconds: 0,
  start_time: '',
  end_time: '',
  is_active: true
})

const totalSeconds = computed(() => {
  return (form.value.hours || 0) * 3600 + (form.value.minutes || 0) * 60 + (form.value.seconds || 0)
})

watch(() => props.task, (task) => {
  if (task) {
    form.value = {
      name: task.name || '',
      description: task.description || '',
      cycle_type: task.cycle_type || 'daily',
      color: task.color || '#6366f1',
      hours: Math.floor((task.countdown_seconds || 0) / 3600),
      minutes: Math.floor(((task.countdown_seconds || 0) % 3600) / 60),
      seconds: (task.countdown_seconds || 0) % 60,
      start_time: task.start_time || '',
      end_time: task.end_time || '',
      is_active: task.is_active !== false
    }
  }
}, { immediate: true })

function handleSubmit() {
  if (totalSeconds.value <= 0) {
    alert('请设置有效的倒计时时长')
    return
  }

  const data = {
    name: form.value.name,
    description: form.value.description,
    cycle_type: form.value.cycle_type,
    color: form.value.color,
    countdown_seconds: totalSeconds.value,
    start_time: form.value.start_time,
    end_time: form.value.end_time,
    is_active: form.value.is_active
  }

  emit('submit', data)
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
  max-width: 520px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
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
  font-size: 20px;
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
  transform: rotate(90deg);
}

.modal-body {
  padding: 24px 28px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
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

.form-group input[type="text"],
.form-group input[type="time"],
.form-group input[type="number"],
.form-group select,
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
.form-group select:focus,
.form-group textarea:focus {
  border-color: #6366f1;
  background: white;
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.color-picker {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
}

.color-option {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  cursor: pointer;
  border: 3px solid transparent;
  transition: all 0.2s;
}

.color-option:hover {
  transform: scale(1.1);
}

.color-option.active {
  border-color: #1e293b;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.color-input {
  width: 32px;
  height: 32px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  padding: 2px;
  cursor: pointer;
  background: white;
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
  border-color: #6366f1;
  background: white;
}

.duration-item span {
  font-size: 13px;
  color: #64748b;
  font-weight: 500;
}

.hint {
  font-size: 12px;
  color: #64748b;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 500;
  cursor: pointer;
}

.checkbox-group input[type="checkbox"] {
  width: 20px;
  height: 20px;
  accent-color: #6366f1;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding-top: 8px;
  border-top: 1px solid #e2e8f0;
  margin-top: 4px;
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
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
}

.btn-submit:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
}
</style>
