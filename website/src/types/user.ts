export interface UserView {
  id: string
  name: string
  username: string
  phone: string
  email: string
  avatar: string
  gender: string
  deptId: string
  deptName: string
  status: number
  roles: RoleView[]
}

export interface RoleView {
  id: string
  name: string
  code: string
}

export interface LoginRequest {
  username: string
  password: string
  captchaVerification?: string
}

export interface LoginResponse {
  token: string
  user: UserView
}

export interface RegisterRequest {
  username: string
  password: string
  name?: string
  role?: string
}

export interface CreateUserRequest {
  username: string
  password: string
  name: string
  phone?: string
  email?: string
  gender?: string
  deptId?: string
  status: number
  roles?: string[]
}

export interface UpdateUserRequest {
  id: string
  name?: string
  phone?: string
  email?: string
  gender?: string
  avatar?: string
  deptId?: string
  status?: number
  username?: string
  password?: string
  oldPassword?: string
  roles?: string[]
}

export interface UserQueryRequest {
  pageIndex?: number
  pageSize?: number
  name?: string
  roleId?: string
  deptId?: string
}

export interface UserOverview {
  totalProjects: number
  totalAnswers: number
  todayAnswers: number
}

export interface RegisterRoleView {
  id: string
  name: string
  code: string
}

// User task view
export interface UserTaskView {
  id: string
  projectId: string
  projectName: string
  status: number
  createAt: string
  updateAt: string
}

// User task query
export interface UserTaskQuery {
  pageIndex?: number
  pageSize?: number
  status?: number
}

// System info update request
export interface SystemInfoRequest {
  name?: string
  description?: string
  avatar?: string
}

// AI settings
export interface AISetting {
  enabled: boolean
  provider?: string
  model?: string
  apiKey?: string
  baseUrl?: string
  [key: string]: any
}
