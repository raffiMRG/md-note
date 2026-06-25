<template>
  <MdPreview :model-value="content" language="en-US" />
</template>

<script setup>
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'

defineProps({ content: { type: String, default: '' } })
</script>
