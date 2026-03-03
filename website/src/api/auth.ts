import request from './index'
import type { LoginRequest, LoginResponse, RegisterRequest, UserView, UserQueryRequest, CreateUserRequest, UpdateUserRequest, UserOverview, RegisterRoleView, UserTaskView, UserTaskQuery } from '@/types/user'
import type { PageResponse } from '@/types/api'

// Public APIs (no auth required)
export const authApi = {
  // Get RSA public key for password encryption
  getRsaPublicKey(): Promise<string> {
    return request.get('/public/rsaPublicKey')
  },

  // Login with username and encrypted password
  login(data: LoginRequest): Promise<LoginResponse> {
    return request.post('/public/login', data)
  },

  // Logout
  logout(): Promise<void> {
    return request.post('/public/logout')
  },

  // Register new user
  register(data: RegisterRequest): Promise<void> {
    return request.post('/public/register', data)
  },

  // Get available roles for registration
  getRegisterRoles(): Promise<RegisterRoleView[]> {
    return request.get('/public/listRegisterRole')
  }
}

// User APIs (auth required)
export const userApi = {
  // Get current user info
  getCurrentUser(): Promise<UserView> {
    return request.get('/currentUser')
  },

  // Get user overview stats
  getUserOverview(): Promise<UserOverview> {
    return request.get('/userOverview')
  },

  // List users with pagination
  listUsers(params: UserQueryRequest): Promise<PageResponse<UserView>> {
    return request.get('/user/list', { params })
  },

  // Get user by ID
  getUser(id: string): Promise<UserView> {
    return request.get(`/user/${id}`)
  },

  // Create user (admin)
  createUser(data: CreateUserRequest): Promise<void> {
    return request.post('/user', data)
  },

  // Update user
  updateUser(data: UpdateUserRequest): Promise<void> {
    return request.put('/user', data)
  },

  // Delete user
  deleteUser(id: string): Promise<void> {
    return request.delete(`/user/${id}`)
  },

  // Bind roles to user
  bindRole(userId: string, roleIds: string[]): Promise<void> {
    return request.post('/user/bindRole', { userId, roleIds })
  },

  // Update password
  updatePassword(oldPassword: string, newPassword: string): Promise<void> {
    return request.put('/user/updatePassword', { oldPassword, newPassword })
  },

  // Import users from Excel
  importUser(file: File): Promise<void> {
    const formData = new FormData()
    formData.append('file', file)
    return request.post('/importUser', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // Get current user's tasks
  listUserTasks(params: UserTaskQuery): Promise<PageResponse<UserTaskView>> {
    return request.get('/listUserTask', { params })
  },

  // Get current user's history tasks
  listHistoryTasks(params: UserTaskQuery): Promise<PageResponse<UserTaskView>> {
    return request.get('/listHistoryTask', { params })
  }
}
