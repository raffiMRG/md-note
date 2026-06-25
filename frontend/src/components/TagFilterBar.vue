<template>
  <div class="tag-filter-bar">
    <button type="button" :class="{ active: !active }" @click="$emit('select', null)">Semua</button>
    <button
      v-for="tag in tags"
      :key="tag.id"
      type="button"
      :class="{ active: active === tag.slug }"
      @click="$emit('select', tag.slug)"
    >
      {{ tag.name }}
    </button>
  </div>
</template>

<script setup>
defineProps({
  tags: { type: Array, default: () => [] },
  active: { type: String, default: null },
})
defineEmits(['select'])
</script>

<style scoped>
.tag-filter-bar {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin: 16px 0;
}
.tag-filter-bar button {
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text);
  border-radius: 999px;
  padding: 4px 12px;
  font-size: 13px;
}
.tag-filter-bar button.active {
  background: var(--accent);
  color: white;
  border-color: var(--accent);
}
</style>
