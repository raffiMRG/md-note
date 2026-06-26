<template>
  <div class="admin-users">
    <h1><i class="fa-solid fa-users-gear"></i> Manajemen Akun</h1>

    <section class="section">
      <p v-if="loading" class="hint">Memuat...</p>
      <p v-else-if="error" class="error">{{ error }}</p>
      <div v-else class="user-table">
        <div class="user-row header-row">
          <span>Username</span>
          <span>Email</span>
          <span>Role</span>
          <span>Dibuat</span>
          <span></span>
        </div>
        <div v-for="u in users" :key="u.id" class="user-row">
          <span class="username">{{ u.username }}</span>
          <span class="email">{{ u.email }}</span>
          <span>
            <select
              :value="u.role"
              :disabled="u.id === currentUserID"
              class="role-select"
              @change="onRoleChange(u, $event.target.value)"
            >
              <option value="penulis">penulis</option>
              <option value="admin">admin</option>
            </select>
          </span>
          <span class="date">{{ formatDate(u.created_at) }}</span>
          <span>
            <button
              class="btn danger icon-btn"
              title="Hapus"
              :disabled="u.id === currentUserID"
              @click="onDelete(u)"
            >
              <i class="fa-solid fa-trash"></i>
            </button>
          </span>
        </div>
        <p v-if="users.length === 0" class="empty">Belum ada user.</p>
      </div>
      <p v-if="actionError" class="error">{{ actionError }}</p>
      <p v-if="actionSuccess" class="success">{{ actionSuccess }}</p>
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { listUsers, updateUserRole, deleteUser } from '../api/users'

const authStore = useAuthStore()
const currentUserID = authStore.user?.id

const users = ref([])
const loading = ref(false)
const error = ref('')
const actionError = ref('')
const actionSuccess = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    const { data } = await listUsers()
    users.value = data.users ?? []
  } catch {
    error.value = 'Gagal memuat daftar user'
  } finally {
    loading.value = false
  }
}

onMounted(load)

async function onRoleChange(user, newRole) {
  actionError.value = ''
  actionSuccess.value = ''
  try {
    await updateUserRole(user.id, newRole)
    user.role = newRole
    actionSuccess.value = `Role ${user.username} diubah menjadi ${newRole}`
  } catch {
    actionError.value = 'Gagal mengubah role'
    await load()
  }
}

async function onDelete(user) {
  if (!confirm(`Hapus akun "${user.username}" (${user.email})? Tindakan ini tidak bisa dibatalkan.`)) return
  actionError.value = ''
  actionSuccess.value = ''
  try {
    await deleteUser(user.id)
    users.value = users.value.filter((u) => u.id !== user.id)
    actionSuccess.value = `Akun ${user.username} berhasil dihapus`
  } catch (e) {
    actionError.value = e.response?.data?.error || 'Gagal menghapus user'
  }
}

function formatDate(iso) {
  if (!iso) return '-'
  return new Date(iso).toLocaleDateString('id-ID', { dateStyle: 'medium' })
}
</script>

<style scoped>
.admin-users {
  max-width: 860px;
}

h1 {
  display: flex;
  align-items: center;
  gap: 10px;
}

.section {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 20px;
  margin-top: 20px;
}

.user-table {
  border: 1px solid var(--border);
  border-radius: 8px;
  overflow: hidden;
}

.user-row {
  display: grid;
  grid-template-columns: 1fr 2fr 110px 120px 48px;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-bottom: 1px solid var(--border);
  font-size: 14px;
}

.user-row:last-child {
  border-bottom: none;
}

.header-row {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  background: var(--accent-bg);
}

.username {
  font-weight: 500;
}

.email {
  color: var(--text-muted);
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.date {
  font-size: 12px;
  color: var(--text-muted);
}

.role-select {
  padding: 4px 8px;
  border: 1px solid var(--border);
  border-radius: 5px;
  font-size: 13px;
  background: var(--surface);
  color: var(--text);
  cursor: pointer;
  width: 100%;
}

.role-select:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.icon-btn {
  padding: 5px 8px;
  font-size: 12px;
}

.empty {
  padding: 16px 14px;
  color: var(--text-muted);
  font-size: 13px;
  margin: 0;
}

.hint {
  font-size: 13px;
  color: var(--text-muted);
}

.error {
  color: var(--danger);
  font-size: 13px;
  margin: 10px 0 0;
}

.success {
  color: #22c55e;
  font-size: 13px;
  margin: 10px 0 0;
}
</style>
