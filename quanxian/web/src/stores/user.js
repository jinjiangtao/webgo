import { defineStore } from 'pinia'
import { login, logout, getUserInfo } from '@/api/auth'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo') || 'null'),
    menus: JSON.parse(localStorage.getItem('menus') || '[]'),
    buttons: JSON.parse(localStorage.getItem('buttons') || '[]')
  }),

  actions: {
    async login(loginForm) {
      const res = await login(loginForm)
      this.token = res.data.token
      this.userInfo = res.data.user
      localStorage.setItem('token', this.token)
      localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
      return res
    },

    async getUserInfo() {
      const res = await getUserInfo()
      this.userInfo = res.data.user
      this.menus = res.data.menus
      this.buttons = res.data.buttons
      localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
      localStorage.setItem('menus', JSON.stringify(this.menus))
      localStorage.setItem('buttons', JSON.stringify(this.buttons))
      return res
    },

    async logout() {
      try {
        await logout()
      } catch (e) {
        console.error(e)
      }
      this.token = ''
      this.userInfo = null
      this.menus = []
      this.buttons = []
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      localStorage.removeItem('menus')
      localStorage.removeItem('buttons')
    },

    hasPermission(code) {
      if (!this.buttons || this.buttons.length === 0) return false
      return this.buttons.includes(code)
    }
  }
})
