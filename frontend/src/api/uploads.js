import client from './client'

export const uploadImage = (file) => {
  const form = new FormData()
  form.append('file', file)
  return client.post('/uploads', form)
}
