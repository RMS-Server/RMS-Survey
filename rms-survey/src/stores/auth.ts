import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, userApi } from '@/api/auth'
import { encryptPassword } from '@/utils/rsa'
import { getToken, setToken, removeToken } from '@/utils/storage'
import type { UserView, RegisterRequest } from '@/types/user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken())
  const user = ref<UserView | null>(null)
  const rsaPublicKey = ref<string | null>(null)

  const isLoggedIn = computed(() => !!token.value && !!user.value)

  async function fetchRsaPublicKey() {
    const key = await authApi.getRsaPublicKey()
    rsaPublicKey.value = key
    return key
  }

  async function login(username: string, password: string) {
    // Get RSA public key if not cached
    if (!rsaPublicKey.value) {
      await fetchRsaPublicKey()
    }

    // Encrypt password
    const encryptedPassword = encryptPassword(password, rsaPublicKey.value!)

    const response = await authApi.login({
      username,
      password: encryptedPassword
    })

    token.value = response.token
    user.value = response.user
    setToken(response.token)

    return response
  }

  async function register(data: RegisterRequest) {
    // Get RSA public key if not cached
    if (!rsaPublicKey.value) {
      await fetchRsaPublicKey()
    }

    // Encrypt password
    const encryptedPassword = encryptPassword(data.password, rsaPublicKey.value!)

    await authApi.register({
      ...data,
      password: encryptedPassword
    })
  }

  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // Ignore logout API errors
    } finally {
      token.value = null
      user.value = null
      removeToken()
    }
  }

  async function fetchCurrentUser() {
    if (!token.value) return null
    try {
      const userData = await userApi.getCurrentUser()
      user.value = userData
      return userData
    } catch {
      logout()
      return null
    }
  }

  return {
    token,
    user,
    rsaPublicKey,
    isLoggedIn,
    fetchRsaPublicKey,
    login,
    register,
    logout,
    fetchCurrentUser
  }
})
