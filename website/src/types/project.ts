export interface ProjectQuery {
  pageIndex?: number
  pageSize?: number
  name?: string
  status?: number
  parentId?: string
  deleted?: boolean
}

// ProjectRequest is for create + id-only ops (delete/restore/destroy).
// Update paths use the sparse types below — mixing full and partial on
// one endpoint is how `setting: {}` wiped IP rules in the past.
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

// Sparse meta update: every field optional, absent fields are not touched.
// Server-side endpoint cannot accept survey/setting at all.
export interface ProjectMetaUpdate {
  id: string
  parentId?: string
  name?: string
  status?: number
  mode?: string
  priority?: number
}

// Survey-only update: replaces the survey schema and nothing else.
export interface ProjectSurveyUpdate {
  id: string
  survey: object
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
  isOwner?: boolean  // true if current user is the project owner
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
