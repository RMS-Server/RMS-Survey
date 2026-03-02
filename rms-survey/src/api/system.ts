import request from './index'
import type { RoleView } from '@/types/user'
import type { PageResponse } from '@/types/api'

// Role APIs
export const roleApi = {
  list(params: { pageIndex?: number; pageSize?: number; name?: string }): Promise<PageResponse<RoleView>> {
    return request.get('/system/role/list', { params })
  },

  get(id: string): Promise<RoleView> {
    return request.get('/system/role/get', { params: { id } })
  },

  create(data: { name: string; code: string; remark?: string; authority?: string }): Promise<void> {
    return request.post('/system/role/create', data)
  },

  update(data: { id: string; name?: string; code?: string; remark?: string; authority?: string; status?: boolean }): Promise<void> {
    return request.post('/system/role/update', data)
  },

  delete(id: string): Promise<void> {
    return request.post('/system/role/delete', { id })
  }
}

// System info APIs
export const systemApi = {
  getSysInfo(): Promise<{ name: string; description: string; avatar: string }> {
    return request.get('/system')
  }
}
