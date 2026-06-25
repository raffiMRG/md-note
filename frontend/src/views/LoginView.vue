<template>
  <div class="auth-page">
    <h1>md-note</h1>
    <h2>Login</h2>
    <form @submit.prevent="onSubmit">
      <label>
        Email
        <input v-model="email" type="email" required />
      </label>
      <label>
        Password
        <input v-model="password" type="password" required />
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

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

async function onSubmit() {
  error.value = ''
  loading.value = true
  try {
    await authStore.login({ email: email.value, password: password.value })
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
.error {
  color: var(--danger);
  font-size: 14px;
}
</style>
