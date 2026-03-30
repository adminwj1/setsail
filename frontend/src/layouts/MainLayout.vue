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

        <el-sub-menu
          v-for="menu in menus"
          :key="menu.id"
          :index="String(menu.id)"
        >
          <template #title>
            <el-icon><component :is="getIcon(menu.icon)" /></el-icon>
            <span>{{ menu.name }}</span>
          </template>
          <el-menu-item
            v-for="child in (menu.children || [])"
            :key="child.id"
            :index="child.path"
          >
            <el-icon><component :is="getIcon(child.icon)" /></el-icon>
            <span>{{ child.name }}</span>
          </el-menu-item>
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
import { computed, markRaw } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  Box,
  Odometer,
  Setting,
  User,
  ArrowDown,
  Goods,
  Document,
  Menu as MenuIcon
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const menus = computed(() => userStore.menus || [])
const activeMenu = computed(() => route.path)

const iconMap = {
  'setting': markRaw(Setting),
  'user': markRaw(User),
  'goods': markRaw(Goods),
  'document': markRaw(Document),
  'menu': markRaw(MenuIcon),
  'Odometer': markRaw(Odometer)
}

const getIcon = (iconName) => {
  return iconMap[iconName] || markRaw(Setting)
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
