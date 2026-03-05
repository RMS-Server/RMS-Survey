<template>
  <a-layout class="layout">
    <!-- Desktop Sider - visible on lg screens -->
    <a-layout-sider
      v-model:collapsed="collapsed"
      collapsible
      theme="light"
      class="glass-sidebar desktop-sider"
      breakpoint="lg"
      :collapsed-width="80"
    >
      <div class="logo">
        <span v-if="!collapsed">RMS Survey</span>
        <span v-else>RS</span>
      </div>
      <a-menu
        v-model:selectedKeys="selectedKeys"
        theme="dark"
        mode="inline"
        @click="handleMenuClick"
      >
        <a-menu-item key="project">
          <template #icon>
            <FileTextOutlined />
          </template>
          <span>问卷管理</span>
        </a-menu-item>
        <a-menu-item key="answer">
          <template #icon>
            <SolutionOutlined />
          </template>
          <span>答题管理</span>
        </a-menu-item>
        <a-menu-item key="template">
          <template #icon>
            <CopyOutlined />
          </template>
          <span>模板管理</span>
        </a-menu-item>
        <a-menu-item key="trash">
          <template #icon>
            <DeleteOutlined />
          </template>
          <span>回收站</span>
        </a-menu-item>
        <a-sub-menu key="system">
          <template #icon>
            <SettingOutlined />
          </template>
          <template #title>系统管理</template>
          <a-menu-item key="system/user">用户管理</a-menu-item>
          <a-menu-item key="system/role">角色管理</a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>

    <!-- Mobile Drawer Menu -->
    <a-drawer
      v-model:open="drawerVisible"
      placement="left"
      :closable="true"
      :width="280"
      class="mobile-drawer"
      @close="drawerVisible = false"
    >
      <template #title>
        <span class="drawer-logo">RMS Survey</span>
      </template>
      <a-menu
        v-model:selectedKeys="selectedKeys"
        mode="inline"
        @click="handleMenuClick"
      >
        <a-menu-item key="project">
          <template #icon>
            <FileTextOutlined />
          </template>
          <span>问卷管理</span>
        </a-menu-item>
        <a-menu-item key="answer">
          <template #icon>
            <SolutionOutlined />
          </template>
          <span>答题管理</span>
        </a-menu-item>
        <a-menu-item key="template">
          <template #icon>
            <CopyOutlined />
          </template>
          <span>模板管理</span>
        </a-menu-item>
        <a-menu-item key="trash">
          <template #icon>
            <DeleteOutlined />
          </template>
          <span>回收站</span>
        </a-menu-item>
        <a-sub-menu key="system">
          <template #icon>
            <SettingOutlined />
          </template>
          <template #title>系统管理</template>
          <a-menu-item key="system/user">用户管理</a-menu-item>
          <a-menu-item key="system/role">角色管理</a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-drawer>

    <a-layout>
      <a-layout-header class="header glass-header">
        <!-- Mobile menu button -->
        <button class="mobile-menu-btn" @click="drawerVisible = true">
          <MenuOutlined />
        </button>
        <div class="header-right">
          <a-dropdown>
            <span class="user-info">
              <UserOutlined />
              <span class="username">{{ authStore.user?.name || authStore.user?.username }}</span>
            </span>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="handleLogout">
                  <LogoutOutlined />
                  <span>退出登录</span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>
      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  FileTextOutlined,
  SolutionOutlined,
  SettingOutlined,
  UserOutlined,
  LogoutOutlined,
  CopyOutlined,
  DeleteOutlined,
  MenuOutlined
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const collapsed = ref(false)
const drawerVisible = ref(false)
const selectedKeys = ref<string[]>(['project'])

const currentRoute = computed(() => route.path)

watch(currentRoute, (path) => {
  if (path.startsWith('/project')) {
    selectedKeys.value = ['project']
  } else if (path.startsWith('/answer')) {
    selectedKeys.value = ['answer']
  } else if (path.startsWith('/template')) {
    selectedKeys.value = ['template']
  } else if (path.startsWith('/trash')) {
    selectedKeys.value = ['trash']
  } else if (path.startsWith('/system/user')) {
    selectedKeys.value = ['system/user']
  } else if (path.startsWith('/system/role')) {
    selectedKeys.value = ['system/role']
  }
}, { immediate: true })

function handleMenuClick({ key }: { key: string }) {
  router.push('/' + key)
  drawerVisible.value = false // Close mobile drawer after navigation
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout {
  min-height: 100vh;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-primary);
  font-size: 18px;
  font-weight: bold;
  border-bottom: 1px solid var(--border-glass);
  box-shadow: var(--shadow-inset);
}

.header {
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-glass);
}

.header-right {
  display: flex;
  align-items: center;
  margin-left: auto;
}

.user-info {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text-main);
}

.username {
  font-size: 14px;
}

.content {
  margin: 0;
  overflow: auto;
  background: transparent;
}

/* Mobile menu button */
.mobile-menu-btn {
  display: none;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: var(--color-text-main);
  font-size: 20px;
  border-radius: var(--radius-sm);
  transition: all 0.3s ease;
}

.mobile-menu-btn:hover {
  background: rgba(252, 121, 97, 0.1);
  color: var(--color-primary);
}

@media (max-width: 991px) {
  .mobile-menu-btn {
    display: flex;
  }
}

/* Drawer logo */
.drawer-logo {
  color: var(--color-primary);
  font-size: 18px;
  font-weight: bold;
}

/* Mobile drawer styling */
:deep(.mobile-drawer .ant-drawer-content) {
  background: var(--surface-glass-strong);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
}

:deep(.mobile-drawer .ant-drawer-header) {
  background: transparent;
  border-bottom: 1px solid var(--border-glass);
}

:deep(.mobile-drawer .ant-menu) {
  background: transparent;
  border: none;
}

:deep(.mobile-drawer .ant-menu-item) {
  color: var(--color-text-main);
  margin: 4px 0;
  border-radius: var(--radius-sm);
}

:deep(.mobile-drawer .ant-menu-item:hover) {
  background: rgba(252, 121, 97, 0.1);
  color: var(--color-primary);
}

:deep(.mobile-drawer .ant-menu-item-selected) {
  background: rgba(252, 121, 97, 0.15);
  color: var(--color-primary);
}

:deep(.mobile-drawer .ant-menu-submenu-title) {
  color: var(--color-text-main);
}

:deep(.mobile-drawer .ant-menu-submenu-title:hover) {
  color: var(--color-primary);
}
</style>
