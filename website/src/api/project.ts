import request from './index'
import type { ProjectQuery, ProjectRequest, ProjectMetaUpdate, ProjectSurveyUpdate, ProjectView, ProjectPartnerQuery, ProjectPartnerRequest, ProjectPartnerView, SelectUserView, SelectRoleView, SelectTemplateView } from '@/types/project'
import type { TemplateQuery } from '@/types/template'
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

  // Update project meta (name, status, mode, priority, parentId).
  // CANNOT modify survey or setting — use updateSurvey / surveyApi.updateSetting.
  update(data: ProjectMetaUpdate): Promise<void> {
    return request.post('/project/update', data)
  },

  // Replace the survey schema JSON. Touches nothing else.
  updateSurvey(data: ProjectSurveyUpdate): Promise<void> {
    return request.post('/project/updateSurvey', data)
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
  },

  // Get deleted projects (recycle bin)
  getTrash(): Promise<ProjectView[]> {
    return request.get('/project/trash')
  },

  // Permanently delete project
  destroy(id: string): Promise<void> {
    return request.post('/project/destroy', { id })
  },

  // Batch destroy projects
  batchDestroy(ids: string[]): Promise<void> {
    return request.post('/project/destroy', { ids })
  },

  // List project partners
  listPartners(params: ProjectPartnerQuery): Promise<PageResponse<ProjectPartnerView>> {
    return request.get('/project/partner/list', { params })
  },

  // Add project partner
  addPartner(data: ProjectPartnerRequest): Promise<void> {
    return request.post('/project/partner/create', data)
  },

  // Remove project partner
  removePartner(data: ProjectPartnerRequest): Promise<void> {
    return request.post('/project/partner/delete', data)
  },

  // Export partners to Excel
  downloadPartners(projectId: string): Promise<Blob> {
    return request.get('/project/partner/download', {
      params: { projectId },
      responseType: 'blob'
    })
  },

  // Import partners from Excel
  importPartners(projectId: string, file: File): Promise<void> {
    const formData = new FormData()
    formData.append('file', file)
    return request.post('/project/partner/import', formData, {
      params: { projectId },
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // Select user for editor picker
  selectUser(name: string): Promise<SelectUserView[]> {
    return request.post('/project/selectUser', { name })
  },

  // Select role for editor picker
  selectRole(name: string): Promise<SelectRoleView[]> {
    return request.post('/project/selectRole', { name })
  },

  // Select template for editor picker
  selectTemplate(params: TemplateQuery): Promise<PageResponse<SelectTemplateView>> {
    return request.post('/project/selectTemplate', params)
  }
}
