<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1>{{ headerTitle }}</h1>
      <router-link class="btn" to="/notes/new"><i class="fa-solid fa-plus"></i> Tulisan Baru</router-link>
    </div>
    <TagFilterBar v-if="!route.query.q" :tags="notesStore.tags" :active="route.query.tag || null" @select="onSelectTag" />
    <p v-if="loading">Memuat...</p>
    <p v-else-if="notesStore.notes.length === 0">Belum ada tulisan.</p>
    <div v-else class="note-grid">
      <NoteCard v-for="note in notesStore.notes" :key="note.id" :note="note" />
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNotesStore } from '../stores/notes'
import NoteCard from '../components/NoteCard.vue'
import TagFilterBar from '../components/TagFilterBar.vue'

const notesStore = useNotesStore()
const route = useRoute()
const router = useRouter()
const loading = ref(false)

const headerTitle = computed(() => (route.query.q ? `Hasil cari: "${route.query.q}"` : 'Tulisan'))

async function load() {
  loading.value = true
  try {
    if (route.query.q) {
      await notesStore.searchNotes(route.query.q)
    } else {
      await notesStore.fetchNotes({ tag: route.query.tag })
    }
  } finally {
    loading.value = false
  }
}

function onSelectTag(slug) {
  router.push({ path: '/', query: slug ? { tag: slug } : {} })
}

onMounted(async () => {
  await notesStore.fetchTags()
  await load()
})

watch(() => [route.query.q, route.query.tag], load)
</script>

<style scoped>
.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}
.note-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
  margin-top: 16px;
}
</style>
