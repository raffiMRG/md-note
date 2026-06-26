import client from './client'

export const exportBackup = () =>
  client.get('/backup', { responseType: 'blob' })

export const importBackup = (data) =>
  client.post('/restore', data)
