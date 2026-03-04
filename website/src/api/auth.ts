import request from './index'
import type { UserView, UserQueryRequest, CreateUserRequest, UpdateUserRequest, UserOverview, UserTaskView, UserTaskQuery } from '@/types/user'
import type { PageResponse } from '@/types/api'

// OAuth APIs (public)
export const authApi = {
  // Get OAuth configuration (auth URL, client ID, etc.)
  getOAuthAuthorize(): Promise<{ authUrl: string; clientId: string; redirectUri: string; scopes: string }> {
    return request.get('/oauth/authorize')
  },

  // Handle OAuth callback with code and PKCE verifier (POST with JSON body)
  oauthCallback(data: { code: string; state: string; codeVerifier: string }): Promise<{ token: string; user: UserView }> {
    return request.post('/oauth/callback', data)
  },

  // Refresh local JWT using stored refresh token
  oauthRefresh(): Promise<{ token: string; user: UserView }> {
    return request.post('/oauth/refresh')
  },

  // Logout and clear OAuth session
  logout(): Promise<void> {
    return request.post('/oauth/logout')
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
