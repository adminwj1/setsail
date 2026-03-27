<template>
  <el-container class="layout-container">
    <el-aside width="200px" class="aside">
      <div class="logo">
        <el-icon><Box /></el-icon>
        <span>ProjectHub</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="aside-menu"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon>
          <span>首页</span>
        </el-menu-item>

        <el-sub-menu index="system" v-if="hasMenu('/system')">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/system/user" v-if="hasMenu('/system/user')">用户管理</el-menu-item>
          <el-menu-item index="/system/role" v-if="hasMenu('/system/role')">角色管理</el-menu-item>
          <el-menu-item index="/system/menu" v-if="hasMenu('/system/menu')">菜单管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="product" v-if="hasMenu('/product')">
          <template #title>
            <el-icon><Goods /></el-icon>
            <span>产品管理</span>
          </template>
          <el-menu-item index="/product/list">产品列表</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="requirement" v-if="hasMenu('/requirement')">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>需求管理</span>
          </template>
          <el-menu-item index="/requirement/list">需求列表</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="route.meta.title">
              {{ route.meta.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-dropdown">
              <el-icon><User /></el-icon>
              {{ userStore.userInfo?.nickname || userStore.userInfo?.username }}
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const hasMenu = (path) => {
  if (!userStore.menus || userStore.menus.length === 0) return true
  const flatMenus = flatMenuTree(userStore.menus)
  return flatMenus.some(m => m.path === path)
}

const flatMenuTree = (menus) => {
  const result = []
  const traverse = (items) => {
    items.forEach(item => {
      if (item.path) {
        result.push(item)
      }
      if (item.children && item.children.length) {
        traverse(item.children)
      }
    })
  }
  traverse(menus)
  return result
}

const handleCommand = async (command) => {
  if (command === 'logout') {
    await userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped lang="scss">
.layout-container {
  height: 100vh;
}

.aside {
  background-color: #304156;

  .logo {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 60px;
    background-color: #2b3a4b;
    color: #fff;
    font-size: 18px;
    font-weight: bold;

    .el-icon {
      margin-right: 8px;
      font-size: 24px;
    }
  }

  .aside-menu {
    border-right: none;
  }
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);

  .header-left {
    display: flex;
    align-items: center;
  }

  .header-right {
    display: flex;
    align-items: center;

    .user-dropdown {
      display: flex;
      align-items: center;
      cursor: pointer;
      padding: 0 12px;

      &:hover {
        background-color: #f5f7fa;
      }
    }
  }
}

.main {
  background-color: #f5f7fa;
  padding: 16px;
}
</style>
