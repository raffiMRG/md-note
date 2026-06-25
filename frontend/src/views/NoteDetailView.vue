<template>
  <div v-if="note" class="notepad-wrapper" :style="{ width: notepadWidth + 'px' }">
    <div class="resize-handle resize-handle--left" @mousedown="startResize($event, -1)" />
    <div class="resize-handle resize-handle--right" @mousedown="startResize($event, 1)" />
  <div class="win-window">

    <!-- ── Title bar ── -->
    <div class="titlebar">
      <div class="titlebar-left">
        <i class="fa-regular fa-file-lines titlebar-icon"></i>
        <span class="titlebar-title">{{ note.title }} — Notepad</span>
      </div>
      <div class="titlebar-controls">
        <button class="wbtn" title="Kembali ke beranda" @click="router.push('/')"><i class="fa-solid fa-minus"></i></button>
        <button class="wbtn" disabled><i class="fa-regular fa-square"></i></button>
        <button class="wbtn wbtn-close" @click="router.push('/')"><i class="fa-solid fa-xmark"></i></button>
      </div>
    </div>

    <!-- ── Menu bar ── -->
    <div class="menubar">
      <span class="menu-item">
        <u>F</u>ile
        <span class="dropdown">
          <span class="dd-item" @click="router.push('/')"><i class="fa-solid fa-house"></i> Ke Beranda</span>
        </span>
      </span>
      <span class="menu-item" :class="{ 'menu-item--disabled': !authStore.isAuthenticated }">
        <u>E</u>dit
        <span v-if="authStore.isAuthenticated" class="dropdown">
          <span class="dd-item" @click="router.push(`/notes/${note.id}/edit`)"><i class="fa-solid fa-pen"></i> Edit Tulisan</span>
          <span class="dd-sep" />
          <span class="dd-item dd-danger" @click="onDelete"><i class="fa-solid fa-trash"></i> Hapus Tulisan</span>
        </span>
      </span>
      <span class="menu-item menu-item--disabled"><u>F</u>ormat</span>
      <span class="menu-item menu-item--disabled"><u>V</u>iew</span>
      <span class="menu-item menu-item--disabled"><u>H</u>elp</span>
    </div>

    <!-- ── Editor area (scrollable content) ── -->
    <div class="editor-area">
      <!-- Tags shown as inline comment at top of file -->
      <div v-if="note.tags && note.tags.length" class="editor-tags">
        <!-- Tags: {{ note.tags.map((t) => t.name).join(', ') }} -->
      </div>
      <MarkdownPreview :content="note.content" />
    </div>

    <!-- ── Status bar ── -->
    <div class="statusbar">
      <span class="statusbar-tags">
        <span v-if="note.tags && note.tags.length">
          <span v-for="tag in note.tags" :key="tag.id" class="status-tag">{{ tag.name }}</span>
        </span>
        <span v-else class="statusbar-muted">Tanpa tag</span>
      </span>
      <span class="statusbar-right">
        Diedit oleh <strong>{{ editorName }}</strong> &middot; {{ formattedDate }}
      </span>
    </div>

  </div><!-- end win-window -->
  </div><!-- end notepad-wrapper -->
  <p v-else class="win-loading">Memuat...</p>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNotesStore } from '../stores/notes'
import { useAuthStore } from '../stores/auth'
import MarkdownPreview from '../components/MarkdownPreview.vue'

const route = useRoute()
const router = useRouter()
const notesStore = useNotesStore()
const authStore = useAuthStore()
const note = ref(null)

// ── Resizable width ──
const notepadWidth = ref(parseInt(localStorage.getItem('notepadWidth') || '920'))

let cleanupResize = null

function startResize(e, direction) {
  e.preventDefault()
  const startX = e.clientX
  const startWidth = notepadWidth.value

  document.body.style.cursor = 'e-resize'
  document.body.style.userSelect = 'none'

  function onMove(e) {
    const delta = (e.clientX - startX) * direction * 2
    const container = document.querySelector('.container--notepad')
    const maxW = container ? container.clientWidth - 48 : window.innerWidth - 48
    notepadWidth.value = Math.max(400, Math.min(maxW, startWidth + delta))
  }

  function onUp() {
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
    localStorage.setItem('notepadWidth', String(notepadWidth.value))
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
    cleanupResize = null
  }

  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
  cleanupResize = onUp
}

onUnmounted(() => { if (cleanupResize) cleanupResize() })

const editorName = computed(
  () =>
    note.value?.updated_by_user?.username ||
    note.value?.created_by_user?.username ||
    'tidak diketahui',
)
const formattedDate = computed(() =>
  note.value
    ? new Date(note.value.updated_at).toLocaleDateString('id-ID', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
      })
    : '',
)

async function load() {
  note.value = await notesStore.fetchNote(route.params.id)
}

async function onDelete() {
  if (!confirm('Hapus tulisan ini?')) return
  await notesStore.deleteNote(route.params.id)
  router.push('/')
}

onMounted(load)
</script>

<style scoped>
/* ── Resizable wrapper (centered, symmetrical) ── */
.notepad-wrapper {
  position: relative;
  margin: 0 auto;
  min-width: 400px;
}

/* ── Resize handles ── */
.resize-handle {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 8px;
  cursor: e-resize;
  z-index: 20;
  border-radius: 4px;
  transition: background 0.15s;
}
.resize-handle--left  { left: -4px; }
.resize-handle--right { right: -4px; }
.resize-handle:hover  { background: rgba(60, 122, 176, 0.35); }

/* ── Window chrome ── */
.win-window {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 160px);
  min-height: 480px;
  box-shadow: 2px 2px 8px rgba(0, 0, 0, 0.5), 0 0 0 1px #3c7ab0;
  border: 1px solid #8ab4d4;
}

/* ── Title bar ── */
.titlebar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(to bottom, #aed4f0 0%, #78b2de 40%, #5c9ecf 100%);
  padding: 4px 4px 4px 8px;
  border-bottom: 1px solid #3c7ab0;
  flex-shrink: 0;
  user-select: none;
}

.titlebar-left {
  display: flex;
  align-items: center;
  gap: 6px;
  overflow: hidden;
}

.titlebar-icon {
  font-size: 16px;
  line-height: 1;
}

.titlebar-title {
  font-size: 13px;
  font-weight: 400;
  color: white;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.4);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* window control buttons */
.titlebar-controls {
  display: flex;
  gap: 2px;
  flex-shrink: 0;
}

.wbtn {
  width: 26px;
  height: 20px;
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 2px;
  background: linear-gradient(to bottom, #c8e0f4, #92c0e0);
  color: #1a1a1a;
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  line-height: 1;
}

.wbtn:hover {
  background: linear-gradient(to bottom, #d8eaff, #aacfee);
}

.wbtn:disabled {
  opacity: 0.5;
  cursor: default;
}

.wbtn-close {
  background: linear-gradient(to bottom, #f0a0a0, #cc4444);
  color: white;
}

.wbtn-close:hover {
  background: linear-gradient(to bottom, #f8b8b8, #e05050);
}

/* ── Menu bar ── */
.menubar {
  display: flex;
  align-items: stretch;
  background: #f0f0f0;
  border-bottom: 1px solid #b8b8b8;
  padding: 0 4px;
  flex-shrink: 0;
}

.menu-item {
  position: relative;
  padding: 3px 8px;
  font-size: 13px;
  cursor: pointer;
  color: #1a1a1a;
  user-select: none;
}

.menu-item:hover > .dropdown {
  display: flex;
}

.menu-item:hover:not(.menu-item--disabled) {
  background: #d8e8f8;
  outline: 1px solid #a0c0e0;
}

.menu-item--disabled {
  color: #888;
  cursor: default;
}

.menu-item--disabled:hover {
  background: none;
  outline: none;
}

/* dropdown */
.dropdown {
  display: none;
  flex-direction: column;
  position: absolute;
  top: 100%;
  left: 0;
  background: #f8f8f8;
  border: 1px solid #b0b0b0;
  box-shadow: 2px 2px 6px rgba(0, 0, 0, 0.2);
  min-width: 160px;
  z-index: 10;
  padding: 2px 0;
}

.dd-item {
  display: block;
  padding: 5px 16px;
  font-size: 13px;
  color: #1a1a1a;
  white-space: nowrap;
  cursor: pointer;
}

.dd-item:hover {
  background: #d8e8f8;
  color: #000;
}

.dd-danger {
  color: #b00000;
}

.dd-sep {
  display: block;
  height: 1px;
  background: #c8c8c8;
  margin: 3px 0;
}

/* ── Editor area ── */
.editor-area {
  flex: 1;
  overflow-y: auto;
  background: #ffffff;
  padding: 8px 12px;
  /* replicate textarea scrollbar style roughly */
  scrollbar-width: thin;
  scrollbar-color: #c0c0c0 #f0f0f0;
}

/* "<!-- Tags: ... -->" comment shown as monospace hint */
.editor-tags {
  font-family: Consolas, 'Courier New', monospace;
  font-size: 13px;
  color: #888;
  margin-bottom: 4px;
}

/* override md-editor preview to look like Notepad text */
.editor-area :deep(.md-editor-preview-wrapper) {
  padding: 0 !important;
  background: transparent !important;
}

.editor-area :deep(.md-editor-preview) {
  font-family: Consolas, 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.55;
  color: #000;
  background: transparent;
}

.editor-area :deep(.md-editor-preview h1),
.editor-area :deep(.md-editor-preview h2),
.editor-area :deep(.md-editor-preview h3),
.editor-area :deep(.md-editor-preview h4) {
  font-family: Consolas, 'Courier New', monospace;
  font-size: 14px;
  font-weight: bold;
  border: none;
  padding: 0;
  margin: 8px 0 2px;
  color: #000;
}

.editor-area :deep(.md-editor-preview strong) {
  font-weight: bold;
}

.editor-area :deep(.md-editor-preview em) {
  font-style: italic;
}

.editor-area :deep(.md-editor-preview code),
.editor-area :deep(.md-editor-preview pre) {
  font-family: Consolas, 'Courier New', monospace;
  font-size: 13px;
  background: transparent;
  border: none;
}

.editor-area :deep(.md-editor-preview blockquote) {
  border-left: none;
  padding: 0 0 0 4ch;
  background: transparent;
  color: #444;
}

.editor-area :deep(.md-editor-preview a) {
  color: #0000ee;
  text-decoration: underline;
}

.editor-area :deep(.md-editor-preview ul),
.editor-area :deep(.md-editor-preview ol) {
  padding-left: 3ch;
}

.editor-area :deep(.md-editor-preview hr) {
  border: none;
  border-top: 1px solid #ccc;
}

/* ── Status bar ── */
.statusbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #f0f0f0;
  border-top: 1px solid #b8b8b8;
  padding: 2px 8px;
  font-size: 12px;
  color: #444;
  flex-shrink: 0;
  gap: 8px;
  min-height: 22px;
}

.statusbar-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  overflow: hidden;
}

.status-tag {
  background: #dce8f4;
  border: 1px solid #a8c4e0;
  border-radius: 2px;
  padding: 0 5px;
  font-size: 11px;
  color: #336;
}

.statusbar-muted {
  color: #999;
  font-style: italic;
}

.statusbar-right {
  white-space: nowrap;
  flex-shrink: 0;
  color: #555;
}

.win-loading {
  text-align: center;
  color: #fff;
  padding: 60px;
  font-size: 14px;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.4);
}

/* ── Dark mode ── */
[data-theme='dark'] .win-window {
  box-shadow: 2px 2px 12px rgba(0, 0, 0, 0.7), 0 0 0 1px #2a2a3a;
  border-color: #3a3a4a;
}

[data-theme='dark'] .titlebar {
  background: linear-gradient(to bottom, #3a3a4e 0%, #252534 40%, #1e1e2e 100%);
  border-bottom-color: #1a1a28;
}

[data-theme='dark'] .wbtn {
  background: linear-gradient(to bottom, #4a4a5e, #333344);
  border-color: rgba(255, 255, 255, 0.15);
  color: #ccc;
}
[data-theme='dark'] .wbtn:hover {
  background: linear-gradient(to bottom, #5a5a6e, #424254);
}
[data-theme='dark'] .wbtn-close {
  background: linear-gradient(to bottom, #9b3030, #7a1e1e);
}
[data-theme='dark'] .wbtn-close:hover {
  background: linear-gradient(to bottom, #b83838, #922424);
}

[data-theme='dark'] .menubar {
  background: #2b2b2b;
  border-bottom-color: #1a1a1a;
}
[data-theme='dark'] .menu-item {
  color: #d0d0d0;
}
[data-theme='dark'] .menu-item:hover:not(.menu-item--disabled) {
  background: #3a4a5e;
  outline-color: #4a6a8a;
}
[data-theme='dark'] .menu-item--disabled {
  color: #666;
}
[data-theme='dark'] .dropdown {
  background: #2d2d2d;
  border-color: #4a4a4a;
  box-shadow: 2px 2px 8px rgba(0, 0, 0, 0.5);
}
[data-theme='dark'] .dd-item {
  color: #d0d0d0;
}
[data-theme='dark'] .dd-item:hover {
  background: #3a4a5e;
}
[data-theme='dark'] .dd-sep {
  background: #444;
}

[data-theme='dark'] .editor-area {
  background: #0d0d0d;
  scrollbar-color: #444 #111;
}

/* override md-editor-v3 internal wrapper */
[data-theme='dark'] .editor-area :deep(.md-editor-preview-wrapper),
[data-theme='dark'] .editor-area :deep(.md-editor-preview) {
  background: #0d0d0d !important;
  color: #c9c9c9;
}

/* headings — biru susu */
[data-theme='dark'] .editor-area :deep(.md-editor-preview h1) {
  color: #cde4ff;
  border-bottom-color: #1e3a5a;
}
[data-theme='dark'] .editor-area :deep(.md-editor-preview h2) {
  color: #b8d4f5;
  border-bottom-color: #1e3a5a;
}
[data-theme='dark'] .editor-area :deep(.md-editor-preview h3),
[data-theme='dark'] .editor-area :deep(.md-editor-preview h4),
[data-theme='dark'] .editor-area :deep(.md-editor-preview h5),
[data-theme='dark'] .editor-area :deep(.md-editor-preview h6) {
  color: #a8c4e8;
}

/* links */
[data-theme='dark'] .editor-area :deep(.md-editor-preview a) {
  color: #79c0ff;
  text-decoration-color: rgba(121, 192, 255, 0.4);
}
[data-theme='dark'] .editor-area :deep(.md-editor-preview a:hover) {
  color: #b0d8ff;
}

/* bold & italic */
[data-theme='dark'] .editor-area :deep(.md-editor-preview strong) {
  color: #f0f0f0;
}
[data-theme='dark'] .editor-area :deep(.md-editor-preview em) {
  color: #d8c4a0;
}

/* inline code — oranye */
[data-theme='dark'] .editor-area :deep(.md-editor-preview code:not(pre code)) {
  background: #1e1a14;
  color: #f09060;
  border: 1px solid #3a2a18;
  border-radius: 3px;
  padding: 1px 5px;
}

/* code block */
[data-theme='dark'] .editor-area :deep(.md-editor-preview pre) {
  background: #0f1923 !important;
  border: 1px solid #1e3040;
  border-radius: 6px;
}
[data-theme='dark'] .editor-area :deep(.md-editor-preview pre code) {
  background: transparent !important;
  color: #adbac7;
  border: none;
  padding: 0;
}

/* blockquote */
[data-theme='dark'] .editor-area :deep(.md-editor-preview blockquote) {
  color: #8fa8c8;
  border-left-color: #3a6ea8;
  background: #0c1520;
}

/* table */
[data-theme='dark'] .editor-area :deep(.md-editor-preview table th) {
  background: #131e2e;
  color: #9ab8d8;
  border-color: #1e3a5a;
}
[data-theme='dark'] .editor-area :deep(.md-editor-preview table td) {
  border-color: #1e2a3a;
  color: #b8b8b8;
}
[data-theme='dark'] .editor-area :deep(.md-editor-preview table tr:nth-child(even)) {
  background: #0c1118;
}

/* HR */
[data-theme='dark'] .editor-area :deep(.md-editor-preview hr) {
  border-color: #1e3040;
}

/* list markers */
[data-theme='dark'] .editor-area :deep(.md-editor-preview ul li::marker),
[data-theme='dark'] .editor-area :deep(.md-editor-preview ol li::marker) {
  color: #5a8ab8;
}

[data-theme='dark'] .statusbar {
  background: #2b2b2b;
  border-top-color: #1a1a1a;
  color: #888;
}
[data-theme='dark'] .status-tag {
  background: #2a3a50;
  border-color: #3a5070;
  color: #8ab8e0;
}
[data-theme='dark'] .statusbar-right {
  color: #777;
}
</style>
