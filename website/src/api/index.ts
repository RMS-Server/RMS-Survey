import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { useAuthStore } from '@/stores/auth'
import { message } from 'ant-design-vue'
import router from '@/router'

const instance: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const authStore = useAuthStore()
    const token = authStore.token
    if (token) {
      config.headers.Authorization = token
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
instance.interceptors.response.use(
  (response: AxiosResponse) => {
    const { data } = response
    // Check business code
    if (data.code !== 0) {
      message.error(data.message || 'Request failed')
      return Promise.reject(new Error(data.message || 'Request failed'))
    }
    return data.data
  },
  (error) => {
    if (error.response) {
      const { status } = error.response
      if (status === 401) {
        const authStore = useAuthStore()
        authStore.logout()
        router.push('/login')
        message.error('Session expired, please login again')
      } else if (status === 403) {
        message.error('No permission')
      } else if (status === 404) {
        message.error('Resource not found')
      } else if (status >= 500) {
        message.error('Server error')
      } else {
        message.error(error.response.data?.message || 'Request failed')
      }
    } else {
      message.error('Network error')
    }
    return Promise.reject(error)
  }
)

export default instance
