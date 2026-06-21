import { defineStore } from 'pinia'
import { taskApi, recordApi, statsApi } from '../api'

export const useStore = defineStore('main', {
  state: () => ({
    tasks: [],
    loading: false,
    currentYear: new Date().getFullYear(),
    currentMonth: new Date().getMonth() + 1
  }),
  actions: {
    async fetchTasks() {
      this.loading = true
      try {
        const res = await taskApi.list()
        this.tasks = res.data.data || []
      } finally {
        this.loading = false
      }
    },
    async createTask(data) {
      const res = await taskApi.create(data)
      await this.fetchTasks()
      return res.data
    },
    async updateTask(id, data) {
      const res = await taskApi.update(id, data)
      await this.fetchTasks()
      return res.data
    },
    async deleteTask(id) {
      await taskApi.remove(id)
      await this.fetchTasks()
    },
    async checkIn(data) {
      const res = await recordApi.checkIn(data)
      await this.fetchTasks()
      return res.data
    },
    async fetchRecords(params) {
      const res = await recordApi.list(params)
      return res.data.data
    },
    async fetchCalendar(params) {
      const res = await statsApi.calendar(params)
      return res.data
    },
    async fetchMonthly(params) {
      const res = await statsApi.monthly(params)
      return res.data.data
    },
    async fetchTaskStats(taskId, params) {
      const res = await statsApi.task(taskId, params)
      return res.data.data
    }
  }
})
