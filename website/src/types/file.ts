export interface FileQuery {
  projectId?: string
}

export interface FileView {
  id: string
  name: string
  originalName: string
  mimeType: string
  size: number
  projectId: string
  createAt: string
}

export interface FileUploadResult {
  id: string
  name: string
  url: string
}
