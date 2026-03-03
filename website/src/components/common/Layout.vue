<template>
  <a-layout class="layout">
    <a-layout-sider v-model:collapsed="collapsed" collapsible>
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
    <a-layout>
      <a-layout-header class="header">
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
  DeleteOutlined
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const collapsed = ref(false)
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
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  background: rgba(255, 255, 255, 0.1);
}

.header {
  background: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}

.username {
  font-size: 14px;
}

.content {
  margin: 0;
  overflow: auto;
  background: #f0f2f5;
}
</style>
