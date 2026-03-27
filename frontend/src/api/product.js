import request from './request'

export const getProductList = (params) => {
  return request.get('/products', { params })
}

export const getProduct = (id) => {
  return request.get(`/products/${id}`)
}

export const createProduct = (data) => {
  return request.post('/products', data)
}

export const updateProduct = (id, data) => {
  return request.put(`/products/${id}`, data)
}

export const deleteProduct = (id) => {
  return request.delete(`/products/${id}`)
}
