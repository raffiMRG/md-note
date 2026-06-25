<template>
  <div id="app-shell">
    <NavBar v-if="!route.meta.hideNavBar" />
    <main :class="['container', { 'container--notepad': route.meta.notepadLayout }]">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import NavBar from './components/NavBar.vue'

const route = useRoute()
</script>

<style scoped>
.container {
  max-width: 960px;
  margin: 0 auto;
  padding: 24px 16px;
}

.container--notepad {
  max-width: 100%;
  padding: 40px 24px 60px;
  background: linear-gradient(135deg, #e4770a 0%, #f5b800 100%);
  min-height: calc(100vh - 57px);
}

:global([data-theme='dark']) .container--notepad {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
}
</style>
