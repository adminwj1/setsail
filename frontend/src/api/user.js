import request from './request'

export const getUserList = (params) => {
  return request.get('/users', { params })
}

export const getUser = (id) => {
  return request.get(`/users/${id}`)
}

export const createUser = (data) => {
  return request.post('/users', data)
}

export const updateUser = (id, data) => {
  return request.put(`/users/${id}`, data)
}

export const deleteUser = (id) => {
  return request.delete(`/users/${id}`)
}
