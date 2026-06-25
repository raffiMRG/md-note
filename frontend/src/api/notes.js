import client from './client'

export function listNotes(params) {
  return client.get('/notes', { params })
}

export function searchNotes(q, params = {}) {
  return client.get('/notes/search', { params: { ...params, q } })
}

export function getNote(id) {
  return client.get(`/notes/${id}`)
}

export function createNote(payload) {
  return client.post('/notes', payload)
}

export function updateNote(id, payload) {
  return client.put(`/notes/${id}`, payload)
}

export function deleteNote(id) {
  return client.delete(`/notes/${id}`)
}
