import request from './request'

export const getRequirementList = (params) => {
  return request.get('/requirements', { params })
}

export const getRequirement = (id) => {
  return request.get(`/requirements/${id}`)
}

export const createRequirement = (data) => {
  return request.post('/requirements', data)
}

export const updateRequirement = (id, data) => {
  return request.put(`/requirements/${id}`, data)
}

export const deleteRequirement = (id) => {
  return request.delete(`/requirements/${id}`)
}

export const getRequirementsByProduct = (productId, params) => {
  return request.get(`/products/${productId}/requirements`, { params })
}
