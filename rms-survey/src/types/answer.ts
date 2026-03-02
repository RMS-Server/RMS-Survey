export interface AnswerQuery {
  pageIndex?: number
  pageSize?: number
  projectId?: string
  id?: string
  tempSave?: number
  deleted?: boolean
}

export interface AnswerRequest {
  id?: string
  projectId: string
  answer: object
  metaInfo?: object
  tempSave?: number
  ids?: string[]
}

export interface AnswerView {
  id: string
  projectId: string
  answer: object
  metaInfo: object
  tempSave: number | null
  examScore: number | null
  createBy: string
  createAt: string
  updateAt: string
}

export interface DownloadQuery {
  projectId: string
  locale?: string
}
