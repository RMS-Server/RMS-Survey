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
