import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import viewsMap from './viewsMap'

// 静态路由（无需权限）
const staticRoutes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    name: 'MainLayout',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '首页', icon: 'Odometer' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: staticRoutes
})

// 存储已添加的动态路由路径
const addedRoutes = new Set()

// 动态添加路由
export const addDynamicRoutes = (menus) => {
  if (!menus || !menus.length) return

  const processMenus = (menuList, parentPath = '') => {
    menuList.forEach(menu => {
      // 只处理类型为菜单的节点（type=2）
      if (menu.type === 2 && menu.path && viewsMap[menu.path]) {
        // 确保不在父布局下重复添加
        const fullPath = menu.path.startsWith('/') ? menu.path : `${parentPath}/${menu.path}`

        if (!addedRoutes.has(fullPath)) {
          const routePath = fullPath.replace(/^\//, '')
          // 添加到 MainLayout 的 children 下
          router.addRoute('MainLayout', {
            path: routePath,
            name: routePath.replace(/\//g, '_'),
            component: viewsMap[fullPath],
            meta: {
              title: menu.name,
              icon: menu.icon
            }
          })
          addedRoutes.add(fullPath)
        }
      }

      // 递归处理子菜单
      if (menu.children && menu.children.length) {
        processMenus(menu.children, menu.path)
      }
    })
  }

  processMenus(menus)
}

// 重置动态路由（退出登录时调用）
export const resetDynamicRoutes = () => {
  addedRoutes.forEach(path => {
    const routeName = path.replace(/^\//, '').replace(/\//g, '_')
    if (router.hasRoute(routeName)) {
      router.removeRoute(routeName)
    }
  })
  addedRoutes.clear()
}

// 路由守卫
router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()

  // 应用初始化时恢复动态路由
  if (userStore.token && userStore.menus && userStore.menus.length > 0) {
    addDynamicRoutes(userStore.menus)
  }

  if (to.path !== '/login' && !userStore.token) {
    next('/login')
  } else if (to.path === '/login' && userStore.token) {
    next('/')
  } else {
    next()
  }
})

export default router
