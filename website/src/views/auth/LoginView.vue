<template>
  <div class="login-container">
    <a-card class="login-card" title="登录">
      <div class="login-content">
        <p class="login-hint">通过 SSO 账号登录问卷系统</p>
        <a-button
          type="primary"
          size="large"
          block
          :loading="loading"
          @click="handleOAuthLogin"
        >
          通过 SSO 登录
        </a-button>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const loading = ref(false)

async function handleOAuthLogin() {
  loading.value = true
  try {
    await authStore.initiateOAuthLogin()
  } catch (error) {
    message.error('发起登录失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.login-content {
  padding: 24px 0;
}

.login-hint {
  text-align: center;
  color: #666;
  margin-bottom: 24px;
}
</style>
