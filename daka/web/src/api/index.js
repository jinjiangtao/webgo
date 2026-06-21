import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

api.interceptors.response.use(
  response => response,
  error => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const taskApi = {
  list: () => api.get('/tasks'),
  get: id => api.get(`/tasks/${id}`),
  create: data => api.post('/tasks', data),
  update: (id, data) => api.put(`/tasks/${id}`, data),
  remove: id => api.delete(`/tasks/${id}`)
}

export const recordApi = {
  list: params => api.get('/records', { params }),
  get: id => api.get(`/records/${id}`),
  remove: id => api.delete(`/records/${id}`),
  checkIn: data => api.post('/check-in', data),
  markAbsent: data => api.post('/mark-absent', data)
}

export const statsApi = {
  calendar: params => api.get('/stats/calendar', { params }),
  monthly: params => api.get('/stats/monthly', { params }),
  task: (taskId, params) => api.get(`/stats/task/${taskId}`, { params })
}

export default api
