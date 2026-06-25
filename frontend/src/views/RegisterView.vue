<template>
  <div class="auth-page">
    <h1>md-note</h1>
    <h2>Daftar</h2>
    <form @submit.prevent="onSubmit">
      <label>
        Username
        <input v-model="username" type="text" required minlength="3" maxlength="50" />
      </label>
      <label>
        Email
        <input v-model="email" type="email" required />
      </label>
      <label>
        Password
        <input v-model="password" type="password" required minlength="6" />
      </label>
      <p v-if="error" class="error">{{ error }}</p>
      <button type="submit" class="btn" :disabled="loading">{{ loading ? 'Mendaftar...' : 'Daftar' }}</button>
    </form>
    <p>Sudah punya akun? <router-link to="/login">Login</router-link></p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const authStore = useAuthStore()
const router = useRouter()

async function onSubmit() {
  error.value = ''
  loading.value = true
  try {
    await authStore.register({ username: username.value, email: email.value, password: password.value })
    router.push('/')
  } catch (err) {
    error.value = err.response?.data?.error || 'Registrasi gagal'
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
