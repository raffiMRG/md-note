import { defineStore } from 'pinia'
import { loginRequest, registerRequest } from '../api/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || sessionStorage.getItem('token') || null,
    user: JSON.parse(localStorage.getItem('user') || sessionStorage.getItem('user') || 'null'),
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin',
  },
  actions: {
    async login(credentials, rememberMe) {
      const { data } = await loginRequest(credentials)
      this.setSession(data.token, data.user, rememberMe)
    },
    async register(payload) {
      const { data } = await registerRequest(payload)
      this.setSession(data.token, data.user, false)
    },
    setSession(token, user, rememberMe) {
      this.token = token
      this.user = user
      const store = rememberMe ? localStorage : sessionStorage
      const other = rememberMe ? sessionStorage : localStorage
      store.setItem('token', token)
      store.setItem('user', JSON.stringify(user))
      other.removeItem('token')
      other.removeItem('user')
    },
    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      sessionStorage.removeItem('token')
      sessionStorage.removeItem('user')
    },
  },
})
