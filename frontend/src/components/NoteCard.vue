<template>
  <router-link class="note-card" :to="`/notes/${note.id}`">
    <h3>{{ note.title }}</h3>
    <p class="excerpt">{{ excerpt }}</p>
    <div class="tags">
      <span v-for="tag in note.tags" :key="tag.id" class="tag-chip">{{ tag.name }}</span>
    </div>
    <p class="meta">Diedit oleh {{ editorName }} &middot; {{ formattedDate }}</p>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ note: { type: Object, required: true } })

const excerpt = computed(() => {
  const plain = props.note.content.replace(/[#*_`>\-]/g, '').trim()
  return plain.length > 140 ? `${plain.slice(0, 140)}...` : plain
})

const editorName = computed(
  () => props.note.updated_by_user?.username || props.note.created_by_user?.username || 'tidak diketahui',
)

const formattedDate = computed(() => new Date(props.note.updated_at).toLocaleString())
</script>

<style scoped>
.note-card {
  display: block;
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 16px;
  color: inherit;
}
.note-card:hover {
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.06);
}
.excerpt {
  color: var(--text-muted);
  font-size: 14px;
  margin: 8px 0;
}
.tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-bottom: 8px;
}
.tag-chip {
  font-size: 12px;
  background: var(--accent-bg);
  color: var(--accent);
  padding: 2px 8px;
  border-radius: 999px;
}
.meta {
  font-size: 12px;
  color: var(--text-muted);
  margin: 0;
}
</style>
