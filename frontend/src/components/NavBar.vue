<template>
  <header class="navbar">
    <router-link to="/" class="brand">md-note</router-link>
    <form class="search" @submit.prevent="onSearch">
      <input v-model="q" type="search" placeholder="Cari tulisan..." />
    </form>
    <div class="nav-actions">
      <button type="button" class="theme-toggle" :title="isDark ? 'Light mode' : 'Dark mode'" @click="toggle">
        <i :class="isDark ? 'fa-solid fa-sun' : 'fa-solid fa-moon'"></i>
      </button>
      <template v-if="authStore.isAuthenticated">
        <router-link class="btn secondary" to="/tags">
          <i class="fa-solid fa-tags"></i> Kelola Tag
        </router-link>
        <template v-if="authStore.isAdmin">
          <router-link class="btn secondary" to="/admin/users">
            <i class="fa-solid fa-users-gear"></i> Admin
          </router-link>
          <router-link class="btn secondary" to="/settings">
            <i class="fa-solid fa-gear"></i>
          </router-link>
        </template>
        <span class="username">
          <i class="fa-solid fa-user"></i> {{ authStore.user?.username }}
          <span v-if="authStore.isAdmin" class="role-badge">admin</span>
        </span>
        <button type="button" class="btn secondary" @click="onLogout">
          <i class="fa-solid fa-right-from-bracket"></i> Logout
        </button>
      </template>
      <template v-else>
        <router-link class="btn secondary" to="/login">Login</router-link>
        <router-link class="btn" to="/register">Daftar</router-link>
      </template>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useTheme } from '../composables/useTheme'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const q = ref(route.query.q || '')
const { isDark, toggle } = useTheme()

function onSearch() {
  if (!q.value.trim()) {
    router.push('/')
    return
  }
  router.push({ path: '/', query: { q: q.value.trim() } })
}

function onLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 24px;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
}
.brand {
  font-weight: 600;
  color: var(--text);
  white-space: nowrap;
}
.search {
  flex: 1;
  max-width: 400px;
}
.search input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: 6px;
}
.nav-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  white-space: nowrap;
}
.username {
  font-size: 14px;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 6px;
}

.role-badge {
  font-size: 11px;
  font-weight: 600;
  background: var(--accent);
  color: white;
  padding: 1px 6px;
  border-radius: 10px;
}
.theme-toggle {
  background: none;
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 5px 8px;
  font-size: 16px;
  line-height: 1;
  cursor: pointer;
  transition: background 0.15s;
}
.theme-toggle:hover {
  background: var(--accent-bg);
}
</style>
