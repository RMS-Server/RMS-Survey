import request from './index'
import type { FileQuery, FileView, FileUploadResult } from '@/types/file'

export const fileApi = {
  // Get file by ID
  get(id: string): Promise<FileView> {
    return request.get('/file', { params: { id } })
  },

  // List files (with projectId filter)
  list(params: FileQuery): Promise<FileView[]> {
    return request.get('/file/list', { params })
  },

  // Upload file
  upload(file: File, projectId?: string): Promise<FileUploadResult> {
    const formData = new FormData()
    formData.append('file', file)
    if (projectId) {
      formData.append('projectId', projectId)
    }
    return request.post('/file/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // Delete file
  delete(id: string): Promise<void> {
    return request.post('/file/delete', { id })
  },

  // Download import template
  downloadTemplate(name?: string): Promise<Blob> {
    return request.get('/file/downloadTemplate', {
      params: { name },
      responseType: 'blob'
    })
  }
}
