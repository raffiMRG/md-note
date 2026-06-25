import { defineStore } from 'pinia'
import {
  listNotes as listNotesApi,
  searchNotes as searchNotesApi,
  getNote as getNoteApi,
  createNote as createNoteApi,
  updateNote as updateNoteApi,
  deleteNote as deleteNoteApi,
} from '../api/notes'
import {
  listTags as listTagsApi,
  createTag as createTagApi,
  updateTag as updateTagApi,
  deleteTag as deleteTagApi,
} from '../api/tags'

export const useNotesStore = defineStore('notes', {
  state: () => ({
    notes: [],
    total: 0,
    currentNote: null,
    tags: [],
  }),
  actions: {
    async fetchNotes(params = {}) {
      const { data } = await listNotesApi(params)
      this.notes = data.notes
      this.total = data.total
    },
    async searchNotes(q, params = {}) {
      const { data } = await searchNotesApi(q, params)
      this.notes = data.notes
      this.total = data.total
    },
    async fetchNote(id) {
      const { data } = await getNoteApi(id)
      this.currentNote = data.note
      return data.note
    },
    async createNote(payload) {
      const { data } = await createNoteApi(payload)
      return data.note
    },
    async updateNote(id, payload) {
      const { data } = await updateNoteApi(id, payload)
      return data.note
    },
    async deleteNote(id) {
      await deleteNoteApi(id)
    },
    async fetchTags() {
      const { data } = await listTagsApi()
      this.tags = data.tags
      return data.tags
    },
    async createTag(name) {
      const { data } = await createTagApi(name)
      this.tags.push(data.tag)
      return data.tag
    },
    async updateTag(id, name) {
      const { data } = await updateTagApi(id, name)
      const idx = this.tags.findIndex((tag) => tag.id === id)
      if (idx !== -1) this.tags[idx] = data.tag
      return data.tag
    },
    async deleteTag(id) {
      await deleteTagApi(id)
      this.tags = this.tags.filter((tag) => tag.id !== id)
    },
  },
})
