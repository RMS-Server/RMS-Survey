// Template query parameters
export interface TemplateQuery {
  pageIndex?: number
  pageSize?: number
  id?: string
  repoId?: string
  name?: string
  category?: string
  mode?: string
  questionType?: string
}

// Template view response
export interface TemplateView {
  id: string
  repoId: string
  serialNo: string
  name: string
  questionType: string
  template: object
  mode: string
  category: string
  tag: string
  priority: number | null
  previewUrl: string
  shared: boolean | null
}

// Template create/update request
export interface TemplateRequest {
  id?: string
  repoId?: string
  name: string
  questionType?: string
  template?: object
  mode?: string
  category?: string
  tag?: string
  priority?: number
  previewUrl?: string
  shared?: boolean
}

// Category query
export interface CategoryQuery {
  mode?: string
}

// Category view
export interface CategoryView {
  id: string
  name: string
}

// Tag query
export interface TagQuery {
  mode?: string
  category?: string
}

// Tag view
export interface TagView {
  id: string
  name: string
}
