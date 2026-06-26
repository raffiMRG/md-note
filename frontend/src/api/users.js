import client from './client'

export const listUsers = () => client.get('/admin/users')
export const updateUserRole = (id, role) => client.put(`/admin/users/${id}/role`, { role })
export const deleteUser = (id) => client.delete(`/admin/users/${id}`)
