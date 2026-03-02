<template>
  <div class="register-container">
    <a-card class="register-card" title="注册">
      <a-form
        :model="formState"
        :rules="rules"
        @finish="handleSubmit"
      >
        <a-form-item name="username">
          <a-input
            v-model:value="formState.username"
            placeholder="用户名"
            size="large"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item name="name">
          <a-input
            v-model:value="formState.name"
            placeholder="显示名称"
            size="large"
          >
            <template #prefix>
              <IdcardOutlined />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item name="password">
          <a-input-password
            v-model:value="formState.password"
            placeholder="密码"
            size="large"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>
        <a-form-item name="confirmPassword">
          <a-input-password
            v-model:value="formState.confirmPassword"
            placeholder="确认密码"
            size="large"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>
        <a-form-item name="role">
          <a-select
            v-model:value="formState.role"
            placeholder="选择角色"
            size="large"
            :loading="rolesLoading"
          >
            <a-select-option
              v-for="role in roles"
              :key="role.id"
              :value="role.code"
            >
              {{ role.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            block
            :loading="loading"
          >
            注册
          </a-button>
        </a-form-item>
        <div class="register-footer">
          <span>已有账号？</span>
          <router-link to="/login">登录</router-link>
        </div>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'
import { UserOutlined, LockOutlined, IdcardOutlined } from '@ant-design/icons-vue'
import type { RegisterRoleView } from '@/types/user'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const rolesLoading = ref(false)
const roles = ref<RegisterRoleView[]>([])

const formState = reactive({
  username: '',
  name: '',
  password: '',
  confirmPassword: '',
  role: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名' },
    { min: 3, message: '用户名至少3个字符' }
  ],
  password: [
    { required: true, message: '请输入密码' },
    { min: 6, message: '密码至少6个字符' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码' },
    {
      validator: (_rule: unknown, value: string) => {
        if (value !== formState.password) {
          return Promise.reject('两次密码不一致')
        }
        return Promise.resolve()
      }
    }
  ]
}

onMounted(async () => {
  rolesLoading.value = true
  try {
    roles.value = await authApi.getRegisterRoles()
  } catch {
    // Ignore error
  } finally {
    rolesLoading.value = false
  }
})

async function handleSubmit() {
  loading.value = true
  try {
    await authStore.register({
      username: formState.username,
      name: formState.name || formState.username,
      password: formState.password,
      role: formState.role
    })
    message.success('注册成功，请登录')
    router.push('/login')
  } catch (error) {
    message.error('注册失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-card {
  width: 400px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.register-footer {
  text-align: center;
  margin-top: 16px;
}
</style>
