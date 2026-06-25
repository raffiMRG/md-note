import client from './client'

export function loginRequest(payload) {
  return client.post('/auth/login', payload)
}

export function registerRequest(payload) {
  return client.post('/auth/register', payload)
}

export function meRequest() {
  return client.get('/me')
}
