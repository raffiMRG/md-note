import client from './client'

export const listCORSOrigins = () => client.get('/cors-origins')
export const addCORSOrigin = (origin) => client.post('/cors-origins', { origin })
export const deleteCORSOrigin = (id) => client.delete(`/cors-origins/${id}`)
