<template>
  <div class="note-editor">
    <h1>{{ isEdit ? 'Edit Tulisan' : 'Tulisan Baru' }}</h1>
    <p v-if="loadingNote">Memuat...</p>
    <template v-else>
      <input v-model="title" class="title-input" placeholder="Judul tulisan" />
      <TagInput v-model="selectedTagIds" :available-tags="notesStore.tags" @create-tag="onCreateTag" />
      <MdEditor v-model="content" language="en-US" style="height: 480px" />
      <div class="actions">
        <button type="button" class="btn" :disabled="saving" @click="onSave">
          <i :class="saving ? 'fa-solid fa-spinner fa-spin' : 'fa-solid fa-floppy-disk'"></i>
          {{ saving ? 'Menyimpan...' : 'Simpan' }}
        </button>
        <router-link class="btn secondary" :to="cancelTarget">
          <i class="fa-solid fa-xmark"></i> Batal
        </router-link>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { useNotesStore } from '../stores/notes'
import TagInput from '../components/TagInput.vue'

const route = useRoute()
const router = useRouter()
const notesStore = useNotesStore()

const isEdit = computed(() => route.name === 'note-edit')
const title = ref('')
const content = ref('')
const selectedTagIds = ref([])
const saving = ref(false)
const loadingNote = ref(isEdit.value)

const cancelTarget = computed(() => (isEdit.value ? `/notes/${route.params.id}` : '/'))

async function onCreateTag(name) {
  const tag = await notesStore.createTag(name)
  selectedTagIds.value.push(tag.id)
}

async function onSave() {
  saving.value = true
  try {
    const payload = { title: title.value, content: content.value, tag_ids: selectedTagIds.value }
    if (isEdit.value) {
      await notesStore.updateNote(route.params.id, payload)
      router.push(`/notes/${route.params.id}`)
    } else {
      const note = await notesStore.createNote(payload)
      router.push(`/notes/${note.id}`)
    }
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await notesStore.fetchTags()
  if (isEdit.value) {
    const note = await notesStore.fetchNote(route.params.id)
    title.value = note.title
    content.value = note.content
    selectedTagIds.value = note.tags.map((tag) => tag.id)
    loadingNote.value = false
  }
})
</script>

<style scoped>
.title-input {
  width: 100%;
  padding: 10px 12px;
  font-size: 20px;
  border: 1px solid var(--border);
  border-radius: 6px;
  margin: 12px 0;
}
.actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
}
</style>
