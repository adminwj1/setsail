// 菜单路径到视图组件的映射
const viewsMap = {
  '/dashboard': () => import('@/views/Dashboard.vue'),
  '/system/user': () => import('@/views/system/User.vue'),
  '/system/role': () => import('@/views/system/Role.vue'),
  '/system/menu': () => import('@/views/system/Menu.vue'),
  '/product/list': () => import('@/views/product/List.vue'),
  '/requirement/list': () => import('@/views/requirement/List.vue')
}

export default viewsMap
