import client from './client'

export function listTags() {
  return client.get('/tags')
}

export function createTag(name) {
  return client.post('/tags', { name })
}

export function updateTag(id, name) {
  return client.put(`/tags/${id}`, { name })
}

export function deleteTag(id) {
  return client.delete(`/tags/${id}`)
}
