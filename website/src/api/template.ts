import request from './index'
import type { TemplateQuery, TemplateView, TemplateRequest, CategoryView, TagView, CategoryQuery, TagQuery } from '@/types/template'
import type { PageResponse } from '@/types/api'

export const templateApi = {
  // List templates with pagination
  list(params: TemplateQuery): Promise<PageResponse<TemplateView>> {
    return request.get('/template/list', { params })
  },

  // Get single template by ID
  get(id: string): Promise<TemplateView> {
    return request.get('/template', { params: { id } })
  },

  // Create new template
  create(data: TemplateRequest): Promise<string> {
    return request.post('/template/create', data)
  },

  // Update existing template
  update(data: TemplateRequest): Promise<void> {
    return request.post('/template/update', data)
  },

  // Delete template
  delete(id: string): Promise<void> {
    return request.post('/template/delete', { id })
  },

  // List template categories
  listCategories(params: CategoryQuery): Promise<CategoryView[]> {
    return request.get('/template/category/list', { params })
  },

  // List template tags
  listTags(params: TagQuery): Promise<TagView[]> {
    return request.get('/template/tag/list', { params })
  }
}
