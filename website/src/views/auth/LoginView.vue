<template>
  <div class="page-shell">
    <div class="page-surface">
      <div class="page-surface__inner">
        <div class="page-content">
          <h1 class="glass-title">RMS Survey</h1>
          <p class="login-hint">通过 SSO 账号登录问卷系统</p>
          <a-button
            type="primary"
            size="large"
            block
            class="login-btn glow-effect"
            :loading="loading"
            @click="handleOAuthLogin"
          >
            通过 SSO 登录
          </a-button>
        </div>
      </div>
    </div>
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
.login-hint {
  text-align: center;
  color: var(--color-text-muted);
  margin-bottom: var(--spacing-xl);
  font-size: 1rem;
}

.login-btn {
  height: 48px;
  font-size: 1.1rem;
  font-weight: 600;
  border-radius: var(--radius-md);
}
</style>
