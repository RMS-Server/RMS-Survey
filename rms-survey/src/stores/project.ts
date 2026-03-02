import { defineStore } from 'pinia'
import { ref } from 'vue'
import { projectApi } from '@/api/project'
import type { ProjectView, ProjectQuery, ProjectRequest } from '@/types/project'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<ProjectView[]>([])
  const currentProject = ref<ProjectView | null>(null)
  const total = ref(0)
  const loading = ref(false)

  async function fetchProjects(params: ProjectQuery) {
    loading.value = true
    try {
      const result = await projectApi.list(params)
      projects.value = result.list
      total.value = result.total
      return result
    } finally {
      loading.value = false
    }
  }

  async function fetchProject(id: string) {
    loading.value = true
    try {
      const project = await projectApi.get(id)
      currentProject.value = project
      return project
    } finally {
      loading.value = false
    }
  }

  async function createProject(data: ProjectRequest) {
    const project = await projectApi.create(data)
    return project
  }

  async function updateProject(data: ProjectRequest) {
    const project = await projectApi.update(data)
    if (currentProject.value?.id === data.id) {
      currentProject.value = project
    }
    return project
  }

  async function deleteProject(id: string) {
    await projectApi.delete(id)
    projects.value = projects.value.filter(p => p.id !== id)
  }

  async function restoreProject(id: string) {
    await projectApi.restore(id)
  }

  function clearCurrentProject() {
    currentProject.value = null
  }

  return {
    projects,
    currentProject,
    total,
    loading,
    fetchProjects,
    fetchProject,
    createProject,
    updateProject,
    deleteProject,
    restoreProject,
    clearCurrentProject
  }
})
