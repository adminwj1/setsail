import request from './request'

export const login = (data) => {
  return request.post('/login', data)
}

export const logout = () => {
  return request.post('/auth/logout')
}

export const getUserInfo = () => {
  return request.get('/auth/userinfo')
}
