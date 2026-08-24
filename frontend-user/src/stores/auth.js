import { defineStore } from 'pinia'
import { req } from '../api'
export const useAuth = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('qxt_utoken') || '',
    user: JSON.parse(localStorage.getItem('qxt_uuser') || 'null'),
  }),
  actions: {
    async login(username, password) {
      const data = await req('/api/v1/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) })
      this.token = data.token; this.user = data.user
      localStorage.setItem('qxt_utoken', data.token)
      localStorage.setItem('qxt_uuser', JSON.stringify(data.user))
    },
    logout() {
      this.token = ''; this.user = null
      localStorage.removeItem('qxt_utoken'); localStorage.removeItem('qxt_uuser')
    },
  },
})
