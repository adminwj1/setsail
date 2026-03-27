import request from './request'

export const getMenus = () => {
  return request.get('/menus/router')
}

export const getMenuTree = () => {
  return request.get('/menus/tree')
}

export const getMenuList = () => {
  return request.get('/menus')
}

export const getMenu = (id) => {
  return request.get(`/menus/${id}`)
}

export const createMenu = (data) => {
  return request.post('/menus', data)
}

export const updateMenu = (id, data) => {
  return request.put(`/menus/${id}`, data)
}

export const deleteMenu = (id) => {
  return request.delete(`/menus/${id}`)
}
