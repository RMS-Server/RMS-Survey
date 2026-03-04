<template>
  <div class="oauth-callback-container">
    <a-spin size="large" tip="正在登录..." />
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
.oauth-callback-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
</style>
