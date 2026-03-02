import request from './index'
import type { ProjectQuery, ProjectRequest, ProjectView } from '@/types/project'
import type { PageResponse } from '@/types/api'

export const projectApi = {
  // List projects with pagination
  list(params: ProjectQuery): Promise<PageResponse<ProjectView>> {
    return request.get('/project/list', { params })
  },

  // Get single project by ID
  get(id: string): Promise<ProjectView> {
    return request.get('/project', { params: { id } })
  },

  // Create new project
  create(data: ProjectRequest): Promise<ProjectView> {
    return request.post('/project/create', data)
  },

  // Update existing project
  update(data: ProjectRequest): Promise<ProjectView> {
    return request.post('/project/update', data)
  },

  // Delete project (soft delete)
  delete(id: string): Promise<void> {
    return request.post('/project/delete', { id })
  },

  // Batch delete projects
  batchDelete(ids: string[]): Promise<void> {
    return request.post('/project/delete', { ids })
  },

  // Restore deleted project
  restore(id: string): Promise<void> {
    return request.post('/project/restore', { id })
  },

  // Get project by code (for public survey)
  getByCode(code: string): Promise<ProjectView> {
    return request.get('/project/getByCode', { params: { code } })
  }
}
