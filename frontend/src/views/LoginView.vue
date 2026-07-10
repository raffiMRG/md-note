<template>
  <div class="auth-page">
    <h1>md-note</h1>
    <h2>Login</h2>
    <form @submit.prevent="onSubmit">
      <label>
        Username
        <input v-model="username" type="text" required />
      </label>
      <label>
        Password
        <input v-model="password" type="password" required />
      </label>
      <label class="remember-me">
        <input v-model="rememberMe" type="checkbox" />
        Ingat saya
      </label>
      <p v-if="error" class="error">{{ error }}</p>
      <button type="submit" class="btn" :disabled="loading">{{ loading ? 'Masuk...' : 'Login' }}</button>
    </form>
    <p>Belum punya akun? <router-link to="/register">Daftar</router-link></p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const username = ref('')
const password = ref('')
const rememberMe = ref(false)
const error = ref('')
const loading = ref(false)

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

async function onSubmit() {
  error.value = ''
  loading.value = true
  try {
    await authStore.login(
      { username: username.value, password: password.value, remember_me: rememberMe.value },
      rememberMe.value,
    )
    router.push(route.query.redirect || '/')
  } catch (err) {
    error.value = err.response?.data?.error || 'Login gagal'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  max-width: 360px;
  margin: 80px auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 14px;
}
input {
  padding: 8px 10px;
  border: 1px solid var(--border);
  border-radius: 6px;
}
.remember-me {
  flex-direction: row;
  align-items: center;
  gap: 8px;
}
.remember-me input {
  padding: 0;
  width: auto;
}
.error {
  color: var(--danger);
  font-size: 14px;
}
</style>
