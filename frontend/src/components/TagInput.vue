<template>
  <div class="tag-input">
    <label v-for="tag in availableTags" :key="tag.id" class="tag-option" :class="{ selected: modelValue.includes(tag.id) }">
      <input type="checkbox" :checked="modelValue.includes(tag.id)" @change="toggle(tag.id)" />
      {{ tag.name }}
    </label>
    <div class="new-tag">
      <input v-model="newTagName" placeholder="Tag baru" @keyup.enter="onCreate" />
      <button type="button" class="btn secondary" @click="onCreate"><i class="fa-solid fa-plus"></i> Tambah</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  modelValue: { type: Array, required: true },
  availableTags: { type: Array, default: () => [] },
})
const emit = defineEmits(['update:modelValue', 'create-tag'])

const newTagName = ref('')

function toggle(id) {
  const next = props.modelValue.includes(id)
    ? props.modelValue.filter((existingId) => existingId !== id)
    : [...props.modelValue, id]
  emit('update:modelValue', next)
}

function onCreate() {
  const name = newTagName.value.trim()
  if (!name) return
  emit('create-tag', name)
  newTagName.value = ''
}
</script>

<style scoped>
.tag-input {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  margin: 12px 0;
}
.tag-option {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: 1px solid var(--border);
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 13px;
  cursor: pointer;
}
.tag-option.selected {
  background: var(--accent-bg);
  border-color: var(--accent);
  color: var(--accent);
}
.new-tag {
  display: flex;
  gap: 6px;
}
.new-tag input {
  padding: 4px 8px;
  border: 1px solid var(--border);
  border-radius: 6px;
  font-size: 13px;
  width: 120px;
}
</style>
