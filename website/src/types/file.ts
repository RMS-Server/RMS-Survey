export interface FileView {
  id: string
  originalName: string
  fileName: string
  filePath: string
  storageType: number | null
}

export interface FileQuery {
  projectId?: string
}

export interface FileUploadResult {
  id: string
  originalName: string
  fileName: string
  filePath: string
  storageType: number | null
}
