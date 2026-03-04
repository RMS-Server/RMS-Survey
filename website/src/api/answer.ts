import request from './index'
import type { AnswerQuery, AnswerView, DownloadQuery, AnswerUploadResult, AnswerCreateRequest } from '@/types/answer'
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
  },

  // Get deleted answers (recycle bin)
  getTrash(params: AnswerQuery): Promise<PageResponse<AnswerView>> {
    return request.get('/answer/trash', { params })
  },

  // Manually create answer
  create(data: AnswerCreateRequest): Promise<void> {
    return request.post('/answer/create', data)
  },

  // Manually update answer
  update(data: AnswerCreateRequest): Promise<void> {
    return request.post('/answer/update', data)
  },

  // Permanently delete answer
  destroy(id: string): Promise<void> {
    return request.post('/answer/destroy', { id })
  },

  // Batch destroy answers
  batchDestroy(ids: string[]): Promise<void> {
    return request.post('/answer/destroy', { ids })
  },

  // Upload answers from Excel (extended version)
  uploadAnswers(projectId: string, file: File, autoSchema?: boolean, parentId?: string): Promise<AnswerUploadResult> {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('projectId', projectId)
    if (autoSchema !== undefined) {
      formData.append('autoSchema', String(autoSchema))
    }
    if (parentId) {
      formData.append('parentId', parentId)
    }
    return request.post('/answer/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // Mark answer as read
  markRead(id: string): Promise<void> {
    return request.post('/answer/read', { id })
  },

  // Mark answer as unread
  markUnread(id: string): Promise<void> {
    return request.post('/answer/unread', { id })
  }
}
