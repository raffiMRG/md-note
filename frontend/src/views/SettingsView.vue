<template>
  <div class="settings">
    <h1>Pengaturan</h1>

    <section class="section">
      <h2><i class="fa-solid fa-shield-halved"></i> Whitelist CORS Origin</h2>
      <p class="hint">
        Daftar origin (protokol + domain + port) yang diizinkan mengakses API ini dari browser.
        Origin dari file konfigurasi server (<code>.env</code>) selalu aktif dan tidak bisa dihapus dari sini.
      </p>

      <div class="origin-list">
        <div v-for="o in origins" :key="o.id" class="origin-row">
          <span class="origin-url"><i class="fa-solid fa-circle-check"></i> {{ o.origin }}</span>
          <span class="origin-date">{{ formatDate(o.created_at) }}</span>
          <button class="btn danger icon-btn" title="Hapus" @click="onDelete(o)">
            <i class="fa-solid fa-trash"></i>
          </button>
        </div>
        <p v-if="origins.length === 0" class="empty">Belum ada origin tambahan.</p>
      </div>

      <form class="add-form" @submit.prevent="onAdd">
        <input
          v-model="newOrigin"
          placeholder="Contoh: http://192.168.0.143:1234"
          :class="{ invalid: addError }"
        />
        <button type="submit" class="btn" :disabled="adding">
          <i :class="adding ? 'fa-solid fa-spinner fa-spin' : 'fa-solid fa-plus'"></i>
          Tambah
        </button>
      </form>
      <p v-if="addError" class="error">{{ addError }}</p>
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { listCORSOrigins, addCORSOrigin, deleteCORSOrigin } from '../api/cors'

const origins = ref([])
const newOrigin = ref('')
const adding = ref(false)
const addError = ref('')

async function load() {
  const { data } = await listCORSOrigins()
  origins.value = data.origins ?? []
}

onMounted(load)

async function onAdd() {
  const val = newOrigin.value.trim().replace(/\/$/, '')
  if (!val) return
  if (!val.startsWith('http://') && !val.startsWith('https://')) {
    addError.value = 'Origin harus dimulai dengan http:// atau https://'
    return
  }
  adding.value = true
  addError.value = ''
  try {
    await addCORSOrigin(val)
    newOrigin.value = ''
    await load()
  } catch (e) {
    addError.value = e.response?.data?.error || 'Gagal menambah origin'
  } finally {
    adding.value = false
  }
}

async function onDelete(o) {
  if (!confirm(`Hapus origin "${o.origin}"?`)) return
  await deleteCORSOrigin(o.id)
  await load()
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('id-ID', { dateStyle: 'medium' })
}
</script>

<style scoped>
.settings {
  max-width: 720px;
}

h2 {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 6px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.hint {
  font-size: 13px;
  color: var(--text-muted);
  margin: 0 0 16px;
  line-height: 1.5;
}

.hint code {
  background: var(--accent-bg);
  color: var(--accent);
  padding: 1px 5px;
  border-radius: 3px;
  font-size: 12px;
}

.origin-list {
  border: 1px solid var(--border);
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 12px;
}

.origin-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-bottom: 1px solid var(--border);
}

.origin-row:last-child {
  border-bottom: none;
}

.origin-url {
  flex: 1;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 7px;
  color: var(--text);
}

.origin-url i {
  color: #22c55e;
  font-size: 12px;
}

.origin-date {
  font-size: 12px;
  color: var(--text-muted);
  white-space: nowrap;
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

.add-form {
  display: flex;
  gap: 8px;
}

.add-form input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: 6px;
  font-size: 13px;
  font-family: 'Courier New', monospace;
}

.add-form input.invalid {
  border-color: var(--danger);
}

.error {
  color: var(--danger);
  font-size: 13px;
  margin: 6px 0 0;
}

.section {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 20px;
  margin-top: 20px;
}
</style>
