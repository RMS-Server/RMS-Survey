<template>
  <div class="page-shell">
    <div class="page-surface">
      <div class="page-surface__inner">
        <div class="page-content">
          <div class="loading-container">
            <div class="loading-spinner"></div>
            <p class="loading-text">正在登录...</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

onMounted(async () => {
  const code = route.query.code as string
  const state = route.query.state as string

  if (!code || !state) {
    message.error('OAuth回调参数缺失')
    router.push('/login')
    return
  }

  try {
    await authStore.handleOAuthCallback(code, state)
    message.success('登录成功')

    // Redirect to originally requested page or project list
    const redirect = route.query.redirect as string
    router.push(redirect || '/project')
  } catch (error) {
    console.error('OAuth callback error:', error)
    message.error('登录失败，请重试')
    router.push('/login')
  }
})
</script>

<style scoped>
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-lg);
}

.loading-text {
  color: var(--color-text-muted);
  font-size: 1rem;
}
</style>
