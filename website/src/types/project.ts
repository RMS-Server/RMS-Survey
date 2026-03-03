export interface ProjectQuery {
  pageIndex?: number
  pageSize?: number
  name?: string
  status?: number
  parentId?: string
  deleted?: boolean
}

export interface ProjectRequest {
  id?: string
  parentId?: string
  name: string
  survey?: object
  setting?: object
  status?: number
  mode?: string
  priority?: number
  ids?: string[]
}

export interface ProjectView {
  id: string
  parentId: string
  name: string
  survey: object
  setting: object
  status: number
  mode: string
  priority: number
  createBy: string
  createAt: string
  updateAt: string
}

// Project partner view
export interface ProjectPartnerView {
  id: string
  projectId: string
  userId: string
  userName: string
  type: number  // 1=system user, 2=imported user
  status: number  // 0=not answered, 1=answered, 2=passed, 3=failed
  createAt: string
}

// Project partner query
export interface ProjectPartnerQuery {
  projectId: string
  pageIndex?: number
  pageSize?: number
  userName?: string
}

// Project partner request
export interface ProjectPartnerRequest {
  id?: string
  projectId: string
  userId?: string
  userName?: string
  type?: number
  ids?: string[]
}

// Select user view (for editor picker)
export interface SelectUserView {
  id: string
  name: string
  username: string
  phone?: string
  email?: string
}

// Select role view (for editor picker)
export interface SelectRoleView {
  id: string
  name: string
  code: string
}

// Select template view (for editor picker)
export interface SelectTemplateView {
  id: string
  name: string
  mode: string
  category: string
}

// Project with deleted flag for trash
export interface ProjectTrashQuery {
  pageIndex?: number
  pageSize?: number
  deleted?: boolean
}
