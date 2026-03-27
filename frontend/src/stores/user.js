import { defineStore } from 'pinia'
import { login, getUserInfo, logout } from '@/api/auth'
import { getMenus } from '@/api/menu'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: null,
    menus: []
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    permissions: (state) => state.menus.map(m => m.code) || []
  },

  actions: {
    async login(username, password) {
      const res = await login({ username, password })
      if (res.code === 200) {
        this.token = res.data.token
        this.userInfo = res.data.user_info
        localStorage.setItem('token', res.data.token)
        return true
      }
      return false
    },

    async fetchUserInfo() {
      const res = await getUserInfo()
      if (res.code === 200) {
        this.userInfo = res.data
      }
    },

    async fetchMenus() {
      const res = await getMenus()
      if (res.code === 200) {
        this.menus = res.data || []
      }
      return this.menus
    },

    async logout() {
      await logout()
      this.token = ''
      this.userInfo = null
      this.menus = []
      localStorage.removeItem('token')
    }
  }
})
