<template>
  <div class="tag-manager">
    <h1>Kelola Tag</h1>
    <form class="new-tag-form" @submit.prevent="onCreate">
      <input v-model="newName" placeholder="Nama tag baru" />
      <button type="submit" class="btn" :disabled="creating"><i class="fa-solid fa-plus"></i> Tambah</button>
    </form>
    <p v-if="error" class="error">{{ error }}</p>
    <table v-if="notesStore.tags.length" class="tag-table">
      <thead>
        <tr>
          <th>Nama</th>
          <th>Slug</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="tag in notesStore.tags" :key="tag.id">
          <td>
            <input
              v-if="editingId === tag.id"
              v-model="editingName"
              class="edit-input"
              @keyup.enter="onSaveEdit(tag.id)"
            />
            <span v-else>{{ tag.name }}</span>
          </td>
          <td class="slug">{{ tag.slug }}</td>
          <td class="row-actions">
            <template v-if="editingId === tag.id">
              <button type="button" class="btn" @click="onSaveEdit(tag.id)"><i class="fa-solid fa-check"></i> Simpan</button>
              <button type="button" class="btn secondary" @click="editingId = null"><i class="fa-solid fa-xmark"></i> Batal</button>
            </template>
            <template v-else>
              <button type="button" class="btn secondary" @click="onEdit(tag)"><i class="fa-solid fa-pen"></i> Edit</button>
              <button type="button" class="btn danger" @click="onDelete(tag)"><i class="fa-solid fa-trash"></i> Hapus</button>
            </template>
          </td>
        </tr>
      </tbody>
    </table>
    <p v-else>Belum ada tag.</p>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useNotesStore } from '../stores/notes'

const notesStore = useNotesStore()
const newName = ref('')
const creating = ref(false)
const error = ref('')
const editingId = ref(null)
const editingName = ref('')

onMounted(() => notesStore.fetchTags())

async function onCreate() {
  const name = newName.value.trim()
  if (!name) return
  creating.value = true
  error.value = ''
  try {
    await notesStore.createTag(name)
    newName.value = ''
  } catch (e) {
    error.value = e.response?.data?.error || 'Gagal membuat tag'
  } finally {
    creating.value = false
  }
}

function onEdit(tag) {
  editingId.value = tag.id
  editingName.value = tag.name
}

async function onSaveEdit(id) {
  const name = editingName.value.trim()
  if (!name) return
  error.value = ''
  try {
    await notesStore.updateTag(id, name)
    editingId.value = null
  } catch (e) {
    error.value = e.response?.data?.error || 'Gagal menyimpan tag'
  }
}

async function onDelete(tag) {
  if (!confirm(`Hapus tag "${tag.name}"? Tag ini akan dilepas dari semua tulisan.`)) return
  error.value = ''
  try {
    await notesStore.deleteTag(tag.id)
  } catch (e) {
    error.value = e.response?.data?.error || 'Gagal menghapus tag'
  }
}
</script>

<style scoped>
.new-tag-form {
  display: flex;
  gap: 8px;
  margin: 16px 0;
}
.new-tag-form input {
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: 6px;
  flex: 1;
  max-width: 280px;
}
.tag-table {
  width: 100%;
  border-collapse: collapse;
}
.tag-table th,
.tag-table td {
  text-align: left;
  padding: 8px 12px;
  border-bottom: 1px solid var(--border);
}
.tag-table .slug {
  color: var(--text-muted);
  font-size: 13px;
}
.edit-input {
  padding: 4px 8px;
  border: 1px solid var(--border);
  border-radius: 6px;
}
.row-actions {
  display: flex;
  gap: 8px;
  white-space: nowrap;
}
.error {
  color: var(--danger);
}
</style>
