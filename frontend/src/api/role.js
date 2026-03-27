import request from './request'

export const getRoleList = () => {
  return request.get('/roles')
}

export const getRole = (id) => {
  return request.get(`/roles/${id}`)
}

export const createRole = (data) => {
  return request.post('/roles', data)
}

export const updateRole = (id, data) => {
  return request.put(`/roles/${id}`, data)
}

export const deleteRole = (id) => {
  return request.delete(`/roles/${id}`)
}

export const getRoleMenus = (id) => {
  return request.get(`/roles/${id}/menus`)
}

export const setRoleMenus = (id, data) => {
  return request.post(`/roles/${id}/menus`, data)
}
