const BASE = '/api'

async function request(path, options = {}) {
  const res = await fetch(BASE + path, {
    headers: { 'Content-Type': 'application/json' },
    ...options
  })
  const json = await res.json()
  if (json.code !== 200) {
    throw new Error(json.message || '请求失败')
  }
  return json.data
}

export const api = {
  getStats: () => request('/stats'),
  getDevices: (params = {}) => {
    const q = new URLSearchParams(params).toString()
    return request('/devices' + (q ? '?' + q : ''))
  },
  getDevice: (id) => request(`/devices/${id}`),
  createDevice: (data) => request('/devices', { method: 'POST', body: JSON.stringify(data) }),
  updateDevice: (id, data) => request(`/devices/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteDevice: (id) => request(`/devices/${id}`, { method: 'DELETE' }),
  updateDeviceData: (id, data) => request(`/devices/${id}/data`, { method: 'POST', body: JSON.stringify(data) }),
  getDeviceTypes: () => request('/device-types'),
  createDeviceType: (data) => request('/device-types', { method: 'POST', body: JSON.stringify(data) }),
  deleteDeviceType: (id) => request(`/device-types/${id}`, { method: 'DELETE' }),
  getAlarms: (status) => request('/alarms' + (status ? '?status=' + status : '')),
  acknowledgeAlarm: (id) => request(`/alarms/${id}/acknowledge`, { method: 'POST' }),
  getLogs: (params = {}) => {
    const q = new URLSearchParams(params).toString()
    return request('/logs' + (q ? '?' + q : ''))
  },
  clearLogs: () => request('/logs', { method: 'DELETE' })
}
