import { defineStore } from 'pinia'
import { req } from '../api'

export const useAuth = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('qxt_token') || '',
    user: JSON.parse(localStorage.getItem('qxt_user') || 'null'),
  }),
  actions: {
    async login(username, password) {
      const data = await req('/api/v1/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) })
      this.token = data.token
      this.user = data.user
      localStorage.setItem('qxt_token', data.token)
      localStorage.setItem('qxt_user', JSON.stringify(data.user))
    },
    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('qxt_token')
      localStorage.removeItem('qxt_user')
    },
  },
})
