import request from './index'
import type { AnswerQuery, AnswerView, DownloadQuery } from '@/types/answer'
import type { PageResponse } from '@/types/api'

export const answerApi = {
  // List answers with pagination
  list(params: AnswerQuery): Promise<PageResponse<AnswerView>> {
    return request.get('/answer/list', { params })
  },

  // Get single answer by ID
  get(id: string): Promise<AnswerView> {
    return request.get('/answer', { params: { id } })
  },

  // Delete answer
  delete(id: string): Promise<void> {
    return request.post('/answer/delete', { id })
  },

  // Batch delete answers
  batchDelete(ids: string[]): Promise<void> {
    return request.post('/answer/delete', { ids })
  },

  // Restore deleted answer
  restore(id: string): Promise<void> {
    return request.post('/answer/restore', { id })
  },

  // Download answers as Excel
  download(params: DownloadQuery): Promise<Blob> {
    return request.get('/answer/download', {
      params,
      responseType: 'blob'
    })
  },

  // Import answers from Excel
  import(projectId: string, file: File): Promise<{ projectId: string; schema: object }> {
    const formData = new FormData()
    formData.append('projectId', projectId)
    formData.append('file', file)
    return request.post('/answer/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}
