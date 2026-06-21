export function formatDuration(seconds) {
  if (!seconds || seconds <= 0) return '00:00:00'
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

export function formatDate(date) {
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export function formatTime(date) {
  const d = new Date(date)
  const h = String(d.getHours()).padStart(2, '0')
  const m = String(d.getMinutes()).padStart(2, '0')
  const s = String(d.getSeconds()).padStart(2, '0')
  return `${h}:${m}:${s}`
}

export function formatDateTime(date) {
  return `${formatDate(date)} ${formatTime(date)}`
}

export function getStatusText(status) {
  const map = {
    'on_time': '准时',
    'checked_in': '已打卡',
    'late': '迟到',
    'makeup': '补卡',
    'absent': '缺勤'
  }
  return map[status] || status
}

export function getStatusColor(status) {
  const map = {
    'on_time': '#10b981',
    'checked_in': '#10b981',
    'late': '#f59e0b',
    'makeup': '#8b5cf6',
    'absent': '#ef4444'
  }
  return map[status] || '#64748b'
}

export function getCycleText(type) {
  const map = {
    'daily': '每日',
    'weekly': '每周',
    'monthly': '每月',
    'custom': '自定义'
  }
  return map[type] || type
}
